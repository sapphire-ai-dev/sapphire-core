package agent

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
	"reflect"
	"strings"
)

type agentLanguage struct {
	agent            *Agent
	langNodes        map[reflect.Type]*langNode
	condMemory       *langCondMemory
	condGenerator    *langCondGenerator
	interpreters     map[reflect.Type]func(concepts map[int]concept, args ...any) concept
	wordPartDict     map[string]map[langPart]bool
	wordConceptDict  map[string]map[concept]bool
	trainParser      *trainSntcParser
	conceptParsers   map[string]func(d *trainSntcData, data map[string]any, args map[int]any) concept
	fieldNamePartIds map[reflect.Type]map[string]int
	assembleRecord   *assembleConceptRecord
	backupRecord     *backupConceptRecord
}

func (l *agentLanguage) findLangNode(class reflect.Type) *langNode {
	if _, seen := l.langNodes[class]; !seen {
		l.langNodes[class] = l.newLangNode(class)
	}

	return l.langNodes[class]
}

func (l *agentLanguage) toSentence(root concept, ctx *sntcCtx) sntcPart {
	return l.toSentenceRec(root, nil, ctx, map[langCond]bool{})
}

func (l *agentLanguage) registerWordPart(word string, part langPart) {
	if _, seen := l.wordPartDict[word]; !seen {
		l.wordPartDict[word] = map[langPart]bool{}
	}
	l.wordPartDict[word][part] = true
}

func (l *agentLanguage) registerWordConcept(word string, c concept, tenseId int) {
	if _, seen := l.wordConceptDict[word]; !seen {
		l.wordConceptDict[word] = map[concept]bool{}
	}
	l.wordConceptDict[word][c] = true
	c.setExplicitName(tenseId, word)
}

func (l *agentLanguage) toSentenceRec(root, parent concept, ctx *sntcCtx, extraConds map[langCond]bool) sntcPart {
	return l.langNodes[reflect.TypeOf(root)].instantiate(root, parent, ctx, extraConds)
}

func (l *agentLanguage) genConds(root concept, ctx *sntcCtx) map[langCond]bool {
	return l.condGenerator.generate(root, ctx)
}

func (l *agentLanguage) initInterpreters() {
	l.interpreters[toReflect[*simpleObject]()] = l.agent.interpretSimpleObject
}

func (l *agentLanguage) listen(msg *world.LangMessage) sntcPart {
	var src, dst object
	if msg.Src != nil {
		if *msg.Src == l.agent.self.worldId {
			src = l.agent.self
		} else {
			src = l.agent.newSimpleObject(map[int]any{partIdObjectWorldId: *msg.Src})
		}
	}
	if msg.Dst != nil {
		if *msg.Dst == l.agent.self.worldId {
			dst = l.agent.self
		} else {
			dst = l.agent.newSimpleObject(map[int]any{partIdObjectWorldId: *msg.Dst})
		}
	}

	ctx := l.newSntcCtx(src, dst)
	return l.fit(strings.Split(msg.Body, " "), ctx)
}

func (l *agentLanguage) fit(sentence []string, ctx *sntcCtx) sntcPart {
	ctx.sentence = sentence
	for pos, word := range sentence {
		for lPart := range l.wordPartDict[word] {
			lPart.fit(pos, ctx)
		}
	}

	var bestMatch *sntcFit
	for _, matches := range ctx.matches[0] {
		for match := range matches {
			if match.end != len(sentence) {
				continue
			}

			if bestMatch == nil || bestMatch.mismatchCount > match.mismatchCount {
				bestMatch = match
			}
		}
	}

	if bestMatch == nil || bestMatch.mismatchCount > 1 || bestMatch.sntc == nil {
		return nil
	}

	interpretedMatch := l.interpretMatch(bestMatch, ctx, nil)
	return interpretedMatch
}

func (l *agentLanguage) interpretMatch(match *sntcFit, ctx *sntcCtx, parent langPart) sntcPart {
	if match.sntc == nil {
		word := ctx.sentence[match.start]
		parentForm, parentIsVf := parent.(*langForm)
		currPart, currIsCvp := match.lang.(*conceptLangPart)
		if !parentIsVf || !currIsCvp {
			return nil
		}

		newC := l.interpreters[parentForm.node.class](map[int]concept{})
		newC.setExplicitName(currPart.tenseId, word)
		l.registerWordPart(word, match.lang)
		match.c = newC
		return &wordSntcPart{
			s: word,
			p: match.lang,
		}
	}

	matchNode, ok := match.sntc.(*sntcNode)
	if !ok {
		return match.sntc
	}

	currForm, isForm := match.lang.(*langForm)
	if !isForm {
		return nil
	}

	for i := range matchNode.parts {
		matchNode.parts[i] = l.interpretMatch(match.children[i], ctx, match.lang)
		if matchNode.parts[i] == nil {
			return nil
		}
	}

	if match.c == nil || match.c.isImaginary() { // try to obtain concept from child
		currForm.interpret(match, ctx)
	}

	return matchNode
}

func (a *Agent) newAgentLanguage() {
	result := &agentLanguage{
		agent:            a,
		langNodes:        map[reflect.Type]*langNode{},
		wordConceptDict:  map[string]map[concept]bool{},
		wordPartDict:     map[string]map[langPart]bool{},
		interpreters:     map[reflect.Type]func(concepts map[int]concept, args ...any) concept{},
		conceptParsers:   map[string]func(d *trainSntcData, data map[string]any, args map[int]any) concept{},
		fieldNamePartIds: map[reflect.Type]map[string]int{},
	}

	result.newCondGenerator()
	result.newCondMemory()
	result.newTrainSntcParser()
	result.initInterpreters()
	result.initConceptParsers()
	result.initFieldPartIds()
	result.newAssembleConceptRecord()
	result.newBackupConceptRecord()

	a.language = result
}

type langCondGenerator struct {
	language *agentLanguage
	funcs    []func(root concept, ctx *sntcCtx) langCond
}

func (g *langCondGenerator) generate(root concept, ctx *sntcCtx) map[langCond]bool {
	result := map[langCond]bool{}
	for _, f := range g.funcs {
		c := f(root, ctx)
		if c != nil {
			c = g.language.condMemory.findCond(c)
			result[c] = true
		}
	}

	return result
}

func (l *agentLanguage) newCondGenerator() {
	if l.condGenerator != nil {
		return
	}

	l.condGenerator = &langCondGenerator{
		language: l,
		funcs:    []func(root concept, ctx *sntcCtx) langCond{},
	}

	l.condGenerator.initFuncs()
}

func (g *langCondGenerator) initFuncs() {
	g.funcs = append(g.funcs, g.generatorObjectSpeaker)
	g.funcs = append(g.funcs, g.generatorObjectListener)
	g.funcs = append(g.funcs, g.generatorMentioned)
}

func (g *langCondGenerator) generatorObjectSpeaker(root concept, ctx *sntcCtx) langCond {
	if _, ok := root.(object); ok && ctx.src != nil {
		return &participantLangCondition{participantTypeId: participantTypeIdSpeaker}
	}
	return nil
}

func (g *langCondGenerator) generatorObjectListener(root concept, ctx *sntcCtx) langCond {
	if _, ok := root.(object); ok && ctx.dst != nil {
		return &participantLangCondition{participantTypeId: participantTypeIdListener}
	}
	return nil
}

func (g *langCondGenerator) generatorMentioned(_ concept, _ *sntcCtx) langCond {
	return &mentionedLangCondition{}
}

type langCondMemory struct {
	language *agentLanguage
	classes  map[reflect.Type]*langCondClassMemory
}

func (m *langCondMemory) findCond(c langCond) langCond {
	class := reflect.TypeOf(c)
	m.newClassMemory(class)
	return m.classes[class].findCond(c)
}

func (l *agentLanguage) newCondMemory() {
	if l.condMemory != nil {
		return
	}

	l.condMemory = &langCondMemory{
		language: l,
		classes:  map[reflect.Type]*langCondClassMemory{},
	}
}

type langCondClassMemory struct {
	parent *langCondMemory
	conds  []langCond
}

func (m *langCondClassMemory) findCond(c langCond) langCond {
	for _, oldCond := range m.conds {
		if c.match(oldCond) {
			return oldCond
		}
	}

	m.conds = append(m.conds, c)
	return c
}

func (m *langCondMemory) newClassMemory(class reflect.Type) {
	if _, seen := m.classes[class]; seen {
		return
	}

	m.classes[class] = &langCondClassMemory{
		parent: m,
		conds:  []langCond{},
	}
}
