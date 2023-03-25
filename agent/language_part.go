package agent

import (
	"reflect"
)

// match: prevent two of the same langParts from being created
// instantiate: given concept, create sentence part
// fit: given sentence, create sentence fit (contains sentence part and some concepts)
// interpret, given fit, fill in all concepts
type langPart interface {
	match(other langPart) bool
	instantiate(root concept, ctx *sntcCtx) sntcPart
	fit(start int, ctx *sntcCtx) []*sntcFit
	interpret(fit *sntcFit, ctx *sntcCtx) bool
}

type abstractLangPart struct {
	f     *langForm
	class reflect.Type
}

func (p *abstractLangPart) match(other *abstractLangPart) bool {
	return p.class == other.class
}

func (f *langForm) newAbstractLangPart(class reflect.Type) *abstractLangPart {
	return &abstractLangPart{
		f:     f,
		class: class,
	}
}

type wordLangPart struct {
	*abstractLangPart
	w string
}

func (p *wordLangPart) match(other langPart) bool {
	o, ok := other.(*wordLangPart)
	return ok && p.abstractLangPart.match(o.abstractLangPart) && p.w == o.w
}

func (p *wordLangPart) instantiate(_ concept, _ *sntcCtx) sntcPart {
	return &wordSntcPart{
		s: p.w,
		p: p,
	}
}

func (p *wordLangPart) fit(start int, ctx *sntcCtx) []*sntcFit {
	if start >= len(ctx.sentence) || ctx.fitStatus(p, start) {
		return ctx.getMatch(start, p)
	}

	var result []*sntcFit
	if ctx.sentence[start] == p.w {
		match := newSntcFit(start, start+1, p.instantiate(nil, nil), nil, p, 0)
		ctx.addMatch(start, p, match)
		result = append(result, match)
	}

	ctx.setFitStatus(p, start, true)
	return result
}

func (p *wordLangPart) interpret(fit *sntcFit, ctx *sntcCtx) bool {
	if fit.c != nil || p.class == nil {
		return true
	}

	candidates := ctx.filterConcepts(p.class)
	conds, seen := ctx.interpretedConds[p.f]
	if !seen {
		conds = map[langCond]*bool{}
	}

	var parentC concept
	if fit.parent != nil && fit.parent.c != nil {
		parentC = fit.parent.c
	}
	candidates = p.f.node.filterConcepts(parentC, candidates, conds, ctx)
	if len(candidates) != 1 {
		return false
	}

	for _, c := range candidates {
		fit.c = c
	}

	return true
}

func (f *langForm) newWordLangPart(class reflect.Type, w string) *wordLangPart {
	return &wordLangPart{
		abstractLangPart: f.newAbstractLangPart(class),
		w:                w,
	}
}

// this is responsible for a single word, difference is that this word depends on the concept
// for a recursively expanding lang part, see recursiveLangPart
type conceptLangPart struct {
	*abstractLangPart
	tenseId int
}

func (p *conceptLangPart) match(other langPart) bool {
	o, ok := other.(*conceptLangPart)
	return ok && p.abstractLangPart.match(o.abstractLangPart) && p.tenseId == o.tenseId
}

func (p *conceptLangPart) instantiate(root concept, _ *sntcCtx) sntcPart {
	return &wordSntcPart{
		s: root.explicitName(p.tenseId),
		p: p,
	}
}

func (p *conceptLangPart) fit(start int, ctx *sntcCtx) []*sntcFit {
	if start >= len(ctx.sentence) || ctx.fitStatus(p, start) {
		return ctx.getMatch(start, p)
	}

	var result []*sntcFit
	for c := range ctx.convCtx.agent.language.wordConceptDict[ctx.sentence[start]] {
		if reflect.TypeOf(c) == p.class && c.explicitName(p.tenseId) == ctx.sentence[start] {
			result = append(result, newSntcFit(start, start+1, p.instantiate(c, nil), c, p, 0))
		}
	}

	// make an attempt to learn new word
	if len(result) == 0 {
		result = append(result, newSntcFit(start, start+1, nil, nil, p, 1))
	}

	for _, match := range result {
		ctx.addMatch(start, p, match)
	}

	ctx.setFitStatus(p, start, true)
	return result
}

func (p *conceptLangPart) interpret(fit *sntcFit, ctx *sntcCtx) bool {
	if fit.c != nil {
		return true
	}

	if fit.start >= len(ctx.sentence) {
		return false
	}

	word := ctx.sentence[fit.start]
	c := p.f.node.agent.language.interpreters[p.class](map[int]concept{})
	p.f.node.agent.language.registerWordConcept(word, c, p.tenseId)
	p.f.node.agent.language.registerWordPart(word, fit.lang)
	fit.c = c
	fit.sntc = &wordSntcPart{
		s: word,
		p: p,
	}

	return true
}

func (f *langForm) newConceptLangPart(class reflect.Type, tenseId int) *conceptLangPart {
	return &conceptLangPart{
		abstractLangPart: f.newAbstractLangPart(class),
		tenseId:          tenseId,
	}
}

type recursiveLangPart struct {
	*abstractLangPart
	partId int
}

func (p *recursiveLangPart) match(other langPart) bool {
	o, ok := other.(*recursiveLangPart)
	return ok && p.abstractLangPart.match(o.abstractLangPart) && p.partId == o.partId
}

func (p *recursiveLangPart) instantiate(root concept, ctx *sntcCtx) sntcPart {
	child := root.part(p.partId)
	return child.abs().agent.language.toSentenceRec(child, root, ctx, map[langCond]bool{
		&recursiveLangCondition{partId: p.partId}: true,
	})
}

func (p *recursiveLangPart) fit(start int, ctx *sntcCtx) []*sntcFit {
	if start >= len(ctx.sentence) || ctx.fitStatus(p, start) {
		return ctx.getMatch(start, p)
	}

	var result []*sntcFit
	agent := p.f.node.agent
	for recImpl := range agent.record.classImpls(p.class) {
		if agent.language.langNodes[recImpl] != nil {
			for _, recForm := range agent.language.langNodes[recImpl].forms {
				result = append(result, recForm.fit(start, ctx)...)
			}
		}
	}

	for _, match := range result {
		ctx.addMatch(start, p, match)
	}
	ctx.setFitStatus(p, start, true)
	return result
}

func (p *recursiveLangPart) interpret(_ *sntcFit, _ *sntcCtx) bool {
	return false
}

func (f *langForm) newRecursiveLangPart(class reflect.Type, partId int) *recursiveLangPart {
	return &recursiveLangPart{
		abstractLangPart: f.newAbstractLangPart(class),
		partId:           partId,
	}
}
