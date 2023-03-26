package agent

import (
    "github.com/stretchr/testify/assert"
    "strings"
    "testing"
)

func TestLangFormConstructor(t *testing.T) {
    agent := newEmptyWorldAgent()
    soLn := agent.language.newLangNode(toReflect[*simpleObject]())
    soLf := soLn.newLangForm()
    assert.Equal(t, soLn, soLf.node)
    assert.Zero(t, soLf.used)
}

func TestLangFormMatch(t *testing.T) {
    agent := newEmptyWorldAgent()
    class := toReflect[*simpleObject]()
    soLn := agent.language.newLangNode(class)
    soLf := soLn.newLangForm()
    soWlp := soLf.newWordLangPart(class, "it")
    soClp := soLf.newConceptLangPart(class, 123)
    assert.False(t, soLf.match(soWlp))
    assert.False(t, soLf.match(agent.language.newLangNode(toReflect[*selfObject]()).newLangForm()))

    newSoLf := soLn.newLangForm()
    assert.True(t, soLf.match(newSoLf))
    soLf.parts = append(soLf.parts, soWlp)
    newSoLf.parts = append(newSoLf.parts, soWlp)
    assert.True(t, soLf.match(newSoLf))
    newSoLf.parts[0] = soClp
    assert.False(t, soLf.match(newSoLf))
}

func TestLangFormInstantiate(t *testing.T) {
    agent := newEmptyWorldAgent()
    class := toReflect[*simpleObject]()
    soLn := agent.language.newLangNode(class)
    soLf := soLn.newLangForm()
    soWlp := soLf.newWordLangPart(class, "it")
    soLf.parts = append(soLf.parts, soWlp)
    sn := soLf.instantiate(nil, nil)
    assert.Equal(t, soLf, sn.parent())
    assert.Equal(t, soWlp.w, strings.Join(sn.str(), " "))
}

func TestLangFormFitWord(t *testing.T) {
    agent := newEmptyWorldAgent()
    class := toReflect[*simpleObject]()
    soLn := agent.language.newLangNode(class)
    soLf := soLn.newLangForm()
    soWlp := soLf.newWordLangPart(class, "it")
    soLf.parts = append(soLf.parts, soWlp)
    ctx := agent.language.newSntcCtx(nil, nil)
    ctx.sentence = strings.Split("it is an apple", " ")
    fits := soLf.fit(0, ctx)
    assert.Len(t, fits, 1)
    assert.Equal(t, fits[0].lang, soLf)
    assert.Nil(t, fits[0].c)
    assert.Equal(t, fits[0].children[0].parent, fits[0])
}

func TestLangFormFitConcept(t *testing.T) {
    agent := newEmptyWorldAgent()
    class := toReflect[*simpleObject]()
    soLn := agent.language.newLangNode(class)
    soLf := soLn.newLangForm()
    tenseId := 123
    soClp := soLf.newConceptLangPart(class, tenseId)
    soLf.parts = append(soLf.parts, soClp)
    ctx := agent.language.newSntcCtx(nil, nil)
    ctx.sentence = strings.Split("bob is a person", " ")
    soSbj := agent.newSimpleObject(456, nil)
    agent.language.registerWordConcept(ctx.sentence[0], soSbj, tenseId)
    fits := soLf.fit(0, ctx)
    assert.Len(t, fits, 1)
    assert.Equal(t, fits[0].c, soSbj)
}
