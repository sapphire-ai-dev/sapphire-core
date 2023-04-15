package agent

import (
	"reflect"
)

// the struct to represent a permanent grammatical structure for a specific concept class
// contains either a single wordLangPart to fit pronouns, a single conceptLangPart to fit named
// concepts, or an arbitrary number of recursiveLangPart to fit phrases
type langForm struct {
	node  *langNode
	parts []langPart
	used  int
}

func (f *langForm) match(other langPart) bool {
	o, ok := other.(*langForm)
	if !ok {
		return false
	}

	if f.node != o.node || len(f.parts) != len(o.parts) {
		return false
	}

	for i := range f.parts {
		if f.parts[i] == o.parts[i] || f.parts[i].match(o.parts[i]) {
			continue
		}

		return false
	}

	return true
}

func (f *langForm) instantiate(root concept, ctx *sntcCtx) sntcPart {
	result := f.node.newSntcNode(f)
	for _, part := range f.parts {
		result.parts = append(result.parts, part.instantiate(root, ctx))
	}

	return result
}

func (f *langForm) fit(start int, ctx *sntcCtx) []*sntcFit {
	if start >= len(ctx.sentence) || ctx.fitStatus(f, start) {
		return ctx.getMatch(start, f)
	}

	result := map[*sntcFit]bool{}
	f.fitRecursive(ctx.sentence, start, start, ctx, []*sntcFit{}, nil, result)
	var resultSlice []*sntcFit
	for match := range result {
		ctx.addMatch(start, f, match)
		resultSlice = append(resultSlice, match)
	}

	ctx.setFitStatus(f, start, true)
	return resultSlice
}

const condTruthThreshold = 0.9

func (f *langForm) interpret(fit *sntcFit, ctx *sntcCtx) bool {
	ctx.interpretedConds[f] = f.assumeCondTruth()
	for i, p := range f.parts {
		if !p.interpret(fit.children[i], ctx) {
			return false
		}
	}

	f.interpretFormConcept(fit)
	if fit.c == nil {
		return false
	}

	for cond, truth := range ctx.interpretedConds[f] {
		f.safeInterpret(fit, cond, truth, ctx)
	}

	return true
}

func (f *langForm) assumeCondTruth() map[langCond]*bool {
	log := f.node.formLog()[f]
	conds, seen, seenTrue, seenFalse := map[langCond]*bool{}, map[langCond]int{}, map[langCond]int{}, map[langCond]int{}
	for _, entry := range log {
		for cond, truth := range entry.condTruth {
			if _, ok := seen[cond]; !ok {
				seen[cond] = 0
				seenTrue[cond] = 0
				seenFalse[cond] = 0
			}

			seen[cond]++
			if truth != nil {
				if *truth {
					seenTrue[cond]++
				} else {
					seenFalse[cond]++
				}
			}
		}
	}

	for cond := range seen {
		if seen[cond] == 0 {
			continue
		}

		tRatio := float64(seenTrue[cond]) / float64(seen[cond])
		fRatio := float64(seenFalse[cond]) / float64(seen[cond])
		if tRatio > condTruthThreshold {
			conds[cond] = ternary(true)
		} else if fRatio > condTruthThreshold {
			conds[cond] = ternary(false)
		} else {
			conds[cond] = nil
		}
	}

	return conds
}

func (f *langForm) interpretFormConcept(fit *sntcFit) {
	if fit.c != nil && fit.c.isImaginary() == false {
		return
	}

	if len(f.parts) == 1 { // form contains exactly 1 wordLangPart or conceptLangPart
		fit.c = fit.children[0].c
		return
	}

	// form contains multiple recursiveLangParts
	concepts := map[int]concept{}
	for i, p := range f.parts {
		rlp, ok := p.(*recursiveLangPart)
		if !ok {
			panic("illegal language form error: this should have been an recursive language part")
		}

		concepts[rlp.partId] = fit.children[i].c
	}

	fit.c = f.node.agent.language.interpreters[f.node.class](concepts)
}

// progress: list of fits currently using, result: out parameter
func (f *langForm) fitRecursive(sntc []string, start, curr int, ctx *sntcCtx,
	progress []*sntcFit, c concept, result map[*sntcFit]bool) {
	//fmt.Println(sntc, f.node.class, "start:", start, "curr:", curr, "progress:", len(progress), "/", len(f.parts))
	if len(progress) == len(f.parts) {
		f.fitForm(start, curr, progress, c, result)
		return
	}

	matches := f.parts[len(progress)].fit(curr, ctx)
	for _, match := range matches {
		newConcept := c
		if newConcept == nil && match.c != nil {
			newConcept = match.c
		}

		f.fitRecursive(sntc, start, match.end, ctx, append(progress, match), newConcept, result)
	}
}

func (f *langForm) fitForm(start, curr int, progress []*sntcFit, c concept, result map[*sntcFit]bool) {
	sp := f.node.newSntcNode(f)
	mismatchCount := 0
	var matchChildren []*sntcFit

	for _, match := range progress {
		sp.parts = append(sp.parts, match.sntc)
		mismatchCount += match.mismatchCount
		matchChildren = append(matchChildren, match)
	}

	if !isNil(c) && reflect.TypeOf(c) != f.node.class {
		c = f.assembleFormConcept(progress)
	}

	newMatch := newSntcFit(start, curr, sp, c, f, mismatchCount)
	result[newMatch] = true
	newMatch.children = matchChildren
	for _, child := range matchChildren {
		child.parent = newMatch
	}
}

func (f *langForm) assembleFormConcept(progress []*sntcFit) concept {
	args := map[int]any{}
	for i, part := range f.parts {
		rlp, ok := part.(*recursiveLangPart)
		if !ok {
			//fmt.Println(f.node.class, reflect.TypeOf(part), reflect.TypeOf(progress[i].c), len(f.parts), part.(*wordLangPart).w)
			panic("part is not recursive")
		}

		args[rlp.partId] = progress[i].c
	}

	for partId, partInfo := range f.node.infos {
		// len == 1 is a workaround, todo replace with some probability threshold constant
		if _, seen := args[partId]; !seen && len(partInfo.implicitIds) == 1 {
			implicitId := 0
			for iid := range partInfo.implicitIds {
				implicitId = iid
			}
			implicitConcept := f.node.agent.language.generateImplicitConcept(implicitId)
			if !isNil(implicitConcept) {
				args[partId] = implicitConcept
			}
		}
	}

	return f.node.agent.language.assembleRecord.assemble(f.node.class, args)
}

func (f *langForm) collectMatchableLangParts(progress []*sntcFit) map[langPart]bool {
	nextPart := f.parts[len(progress)]
	matchableParts := map[langPart]bool{nextPart: true}
	if rvp, ok := nextPart.(*recursiveLangPart); ok {
		delete(matchableParts, rvp)
		for recImpl := range f.node.agent.record.partImpls(f.node.class, rvp.partId) {
			if f.node.agent.language.langNodes[recImpl] != nil {
				for _, recForm := range f.node.agent.language.langNodes[recImpl].forms {
					matchableParts[recForm] = true
				}
			}
		}
	}

	return matchableParts
}

func (f *langForm) safeInterpret(fit *sntcFit, cond langCond, truth *bool, ctx *sntcCtx) {
	var parentC concept
	if fit.parent != nil {
		parentC = fit.parent.c
	}

	satisfied := cond.satisfied(fit.c, parentC, ctx)
	if truth == nil || satisfied == nil || truth != satisfied {
		return
	}

	_, newConcepts := cond.interpret(fit.c, parentC, truth, ctx)
	for _, c := range newConcepts {
		ctx.newConcepts[c.id()] = c
		if r, ok := c.(relation); ok {
			r.interpret()
		}
	}
}

func (n *langNode) newLangForm() *langForm {
	return &langForm{
		node: n,
		used: 0,
	}
}
