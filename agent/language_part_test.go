package agent

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestWordLangPartConstructor(t *testing.T) {
	agent := newEmptyWorldAgent()
	class := toReflect[*simpleObject]()
	soLn := agent.language.newLangNode(class)
	soLf := soLn.newLangForm()
	word := "it"
	soWlp := soLf.newWordLangPart(class, word)
	assert.Equal(t, class, soWlp.class)
	assert.Equal(t, soLf, soWlp.f)
	assert.Equal(t, word, soWlp.w)
	assert.True(t, soWlp.match(soLf.newWordLangPart(class, word)))
	assert.False(t, soWlp.match(soLf.newWordLangPart(class, word+"_")))
}

func TestWordLangPartInstantiate(t *testing.T) {
	agent := newEmptyWorldAgent()
	class, word := toReflect[*simpleObject](), "it"
	soWlp := agent.language.newLangNode(class).newLangForm().newWordLangPart(class, word)
	wsp := soWlp.instantiate(nil, nil).(*wordSntcPart)
	assert.Equal(t, wsp.s, soWlp.w)
	assert.Equal(t, wsp.p, soWlp)
}

func TestWordLangPartFit(t *testing.T) {
	agent := newEmptyWorldAgent()
	class, word := toReflect[*simpleObject](), "it"
	soWlp := agent.language.newLangNode(class).newLangForm().newWordLangPart(class, word)
	soSrc, soDst := agent.newSimpleObject(1, nil), agent.newSimpleObject(2, nil)
	ctx := agent.language.newSntcCtx(soSrc, soDst)
	ctx.sentence = strings.Split("it is an apple", " ")
	assert.Empty(t, soWlp.fit(1, ctx))
	fits := soWlp.fit(0, ctx)
	assert.Len(t, fits, 1)
	assert.Equal(t, strings.Join(fits[0].sntc.str(), " "), soWlp.w)
}

func TestWordLangPartInterpret(t *testing.T) {
	agent := newEmptyWorldAgent()
	class, word := toReflect[*simpleObject](), "it"
	soWlp := agent.language.newLangNode(class).newLangForm().newWordLangPart(class, word)
	ctx := agent.language.newSntcCtx(nil, nil)
	ctx.sentence = strings.Split("it is an apple", " ")
	fit := soWlp.fit(0, ctx)[0]
	assert.Nil(t, fit.c)
	assert.False(t, soWlp.interpret(fit, ctx))

	soSbj := agent.newSimpleObject(3, nil)
	ctx.convCtx.mentioned[soSbj.id()] = soSbj
	assert.True(t, soWlp.interpret(fit, ctx))
	assert.Equal(t, soSbj, fit.c)
	assert.True(t, soWlp.interpret(fit, ctx))
}

func TestConceptLangPartConstructor(t *testing.T) {
	agent := newEmptyWorldAgent()
	class := toReflect[*simpleObject]()
	soLn := agent.language.newLangNode(class)
	soLf := soLn.newLangForm()
	tenseId := 123
	soClp := soLf.newConceptLangPart(class, tenseId)
	assert.Equal(t, class, soClp.class)
	assert.Equal(t, soLf, soClp.f)
	assert.Equal(t, tenseId, soClp.tenseId)

	assert.True(t, soClp.match(soLf.newConceptLangPart(class, tenseId)))
	assert.False(t, soClp.match(soLf.newConceptLangPart(class, tenseId+1)))
}

func TestConceptLangPartInstantiate(t *testing.T) {
	agent := newEmptyWorldAgent()
	class, tenseId := toReflect[*simpleObject](), 123
	soClp := agent.language.newLangNode(class).newLangForm().newConceptLangPart(class, tenseId)
	soSbj := agent.newSimpleObject(1, nil)
	csp := soClp.instantiate(soSbj, nil).(*wordSntcPart)
	assert.Equal(t, csp.s, "")
	assert.Equal(t, csp.p, soClp)

	word := "apple"
	soSbj.setExplicitName(tenseId, word)
	csp = soClp.instantiate(soSbj, nil).(*wordSntcPart)
	assert.Equal(t, csp.s, word)
	assert.Equal(t, csp.p, soClp)
}

func TestConceptLangPartFit(t *testing.T) {
	agent := newEmptyWorldAgent()
	class, tenseId := toReflect[*simpleObject](), 123
	soClp := agent.language.newLangNode(class).newLangForm().newConceptLangPart(class, tenseId)
	soSrc, soDst := agent.newSimpleObject(1, nil), agent.newSimpleObject(2, nil)
	ctx := agent.language.newSntcCtx(soSrc, soDst)
	ctx.sentence = strings.Split("bob is a person", " ")

	// the word is unknown, but the fit will pass in attempt to learn a new word
	fits := soClp.fit(0, ctx)
	assert.Len(t, fits, 1)
	assert.Nil(t, fits[0].c)
	assert.Nil(t, fits[0].sntc)
	assert.Equal(t, soClp, fits[0].lang)

	ctx = agent.language.newSntcCtx(soSrc, soDst)
	ctx.sentence = strings.Split("bob is a person", " ")
	// the word is known, the fit will match it
	soSbj := agent.newSimpleObject(456, nil)
	agent.language.registerWordConcept(ctx.sentence[0], soSbj, tenseId)
	fits = soClp.fit(0, ctx)
	assert.Len(t, fits, 1)
	assert.Equal(t, fits[0].c, soSbj)
	assert.Equal(t, ctx.sentence[0], strings.Join(fits[0].sntc.str(), " "))
}

func TestConceptLangPartInterpret(t *testing.T) {
	agent := newEmptyWorldAgent()
	class, tenseId := toReflect[*simpleObject](), 123
	soClp := agent.language.newLangNode(class).newLangForm().newConceptLangPart(class, tenseId)
	ctx := agent.language.newSntcCtx(nil, nil)
	ctx.sentence = strings.Split("bob is a person", " ")
	fit := soClp.fit(0, ctx)[0]
	assert.Nil(t, fit.c)
	assert.Nil(t, fit.sntc)
	assert.True(t, soClp.interpret(fit, ctx))
	assert.Equal(t, ctx.sentence[0], fit.c.explicitName(tenseId))
	assert.Equal(t, ctx.sentence[0], strings.Join(fit.sntc.str(), " "))

	assert.True(t, soClp.interpret(fit, ctx))

	fit.start, fit.c = len(ctx.sentence), nil
	assert.False(t, soClp.interpret(fit, ctx))
}

func TestRecursiveLangPartConstructor(t *testing.T) {
	agent := newEmptyWorldAgent()
	class := toReflect[*simpleObject]()
	soLn := agent.language.newLangNode(class)
	soLf := soLn.newLangForm()
	partId := partIdObjectT
	soRlp := soLf.newRecursiveLangPart(class, partId)
	assert.Equal(t, class, soRlp.class)
	assert.Equal(t, soLf, soRlp.f)
	assert.Equal(t, partId, soRlp.partId)

	assert.True(t, soRlp.match(soLf.newRecursiveLangPart(class, partId)))
	assert.False(t, soRlp.match(soLf.newRecursiveLangPart(class, partId+1)))

	assert.False(t, soRlp.interpret(nil, nil))
}

func TestRecursiveLangPartInstantiate(t *testing.T) {
	agent := newEmptyWorldAgent()
	am, amt, partId, tenseId := toReflect[*aspectModifier](), toReflect[*aspectModifierType](), partIdModifierT, 0
	amRlp := agent.language.newLangNode(am).newLangForm().newRecursiveLangPart(am, partId)
	amtLn := agent.language.newLangNode(amt)
	agent.language.langNodes[amt] = amtLn
	amtLf := amtLn.newLangForm()
	amtLn.selectForm(amtLf)
	amtClp := amtLf.newConceptLangPart(amt, tenseId)
	amtLf.parts = []langPart{amtClp}
	word := "red"
	tc := agent.newTestConcept(1, nil)
	redType := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", word), nil)
	redInst := agent.newAspectModifier(redType, tc, conceptSourceObservation, nil, map[string]any{})
	agent.language.registerWordConcept(word, redType, tenseId)

	sn := amRlp.instantiate(redInst, nil)
	assert.Equal(t, sn.parent(), amtLf)
	assert.Equal(t, word, strings.Join(sn.str(), " "))
}
