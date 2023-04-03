package agent

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConceptComponentGeneralizationEmpty(t *testing.T) {
	agent := newEmptyWorldAgent()
	amt1 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "red"), nil)
	amt2 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "blue"), nil)
	sot1 := agent.newSimpleObjectType(conceptSourceObservation, map[int]modifierType{
		amt1.cid: amt1,
	}, nil) // red and sweet
	sot2 := agent.newSimpleObjectType(conceptSourceObservation, map[int]modifierType{
		amt2.cid: amt2,
	}, nil) // blue and sweet
	sot1.generalize(sot2)
	assert.Empty(t, sot1.generalizations())
	assert.Empty(t, sot2.generalizations())
	assert.Empty(t, sot1.specifications())
	assert.Empty(t, sot2.specifications())
}

func TestConceptComponentGeneralizationBasic(t *testing.T) {
	agent := newEmptyWorldAgent()
	amt1 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "red"), nil)
	amt2 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "blue"), nil)
	amt3 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "taste", "sweet"), nil)
	sot1 := agent.newSimpleObjectType(conceptSourceObservation, map[int]modifierType{
		amt1.cid: amt1, amt3.cid: amt3,
	}, nil) // red and sweet
	sot2 := agent.newSimpleObjectType(conceptSourceObservation, map[int]modifierType{
		amt2.cid: amt2, amt3.cid: amt3,
	}, nil) // blue and sweet

	assert.Empty(t, sot1.generalizations())
	assert.Empty(t, sot2.generalizations())
	assert.Empty(t, sot1.specifications())
	assert.Empty(t, sot2.specifications())

	sot1.generalize(sot2)
	assert.Len(t, sot1.generalizations(), 1)
	assert.Len(t, sot2.generalizations(), 1)
	assert.Empty(t, sot1.specifications())
	assert.Empty(t, sot2.specifications())
	var sotG *simpleObjectType
	for _, gen := range sot1.generalizations() {
		sotG = gen.(*simpleObjectType)
	}

	assert.Len(t, sotG.specifications(), 2)
	assert.Contains(t, sotG.specifications(), sot1.id())
	assert.Contains(t, sotG.specifications(), sot2.id())

	assert.Equal(t, sot1.lowestCommonGeneralization(sot2), sotG)
}

func TestConceptComponentGeneralizationTwoLevel(t *testing.T) {
	agent := newEmptyWorldAgent()
	amt11 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "red"), nil)
	amt12 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "blue"), nil)
	amt21 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "taste", "sour"), nil)
	amt22 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "taste", "sweet"), nil)
	amt31 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "shape", "round"), nil)
	sot111 := agent.newSimpleObjectType(conceptSourceObservation, map[int]modifierType{
		amt11.cid: amt11, amt21.cid: amt21, amt31.cid: amt31,
	}, nil) // red and sour and round
	sot211 := agent.newSimpleObjectType(conceptSourceObservation, map[int]modifierType{
		amt12.cid: amt12, amt21.cid: amt21, amt31.cid: amt31,
	}, nil) // blue and sour and round
	sot121 := agent.newSimpleObjectType(conceptSourceObservation, map[int]modifierType{
		amt11.cid: amt11, amt22.cid: amt22, amt31.cid: amt31,
	}, nil) // red and sweet and round
	sot221 := agent.newSimpleObjectType(conceptSourceObservation, map[int]modifierType{
		amt12.cid: amt12, amt22.cid: amt22, amt31.cid: amt31,
	}, nil) // blue and sweet and round
	sot111.generalize(sot211)
	sot121.generalize(sot221)

	sotX11 := sot111.lowestCommonGeneralization(sot211)
	sotX21 := sot121.lowestCommonGeneralization(sot221)
	assert.Len(t, sotX11.specifications(), 2)
	assert.Len(t, sotX21.specifications(), 2)
	assert.Contains(t, sotX11.specifications(), sot111.id())
	assert.Contains(t, sotX11.specifications(), sot211.id())
	assert.Contains(t, sotX21.specifications(), sot121.id())
	assert.Contains(t, sotX21.specifications(), sot221.id())
	assert.Empty(t, sotX11.generalizations())
	assert.Empty(t, sotX21.generalizations())

	sotX11.generalize(sotX21)
	assert.Len(t, sotX11.generalizations(), 1)
	assert.Len(t, sotX21.generalizations(), 1)
	sotXX1 := sotX11.lowestCommonGeneralization(sotX21)
	assert.Len(t, sotXX1.specifications(), 6)

	assert.Len(t, sot111._generalizations, 1)
	assert.Len(t, sot111.generalizations(), 2)
	assert.Len(t, sot221._generalizations, 1)
	assert.Len(t, sot221.generalizations(), 2)
	assert.Equal(t, sotXX1, sot111.lowestCommonGeneralization(sot221))

	sot111.generalize(sot221)
	assert.Len(t, sot111._generalizations, 1)
	assert.Len(t, sot111.generalizations(), 2)
	assert.Len(t, sot221._generalizations, 1)
	assert.Len(t, sot221.generalizations(), 2)
	assert.Equal(t, sotXX1, sot111.lowestCommonGeneralization(sot221))
}

func TestConceptComponentGeneralizationDifferentLevel(t *testing.T) {
	agent := newEmptyWorldAgent()
	amt1 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "red"), nil)
	amt2 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "blue"), nil)
	amt3 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "taste", "sweet"), nil)
	sot1 := agent.newSimpleObjectType(conceptSourceObservation, map[int]modifierType{
		amt1.cid: amt1, amt3.cid: amt3,
	}, nil) // red and sweet
	sot2 := agent.newSimpleObjectType(conceptSourceObservation, map[int]modifierType{
		amt2.cid: amt2, amt3.cid: amt3,
	}, nil) // blue and sweet
	sot1.generalize(sot2)
	sotG := sot1.lowestCommonGeneralization(sot2)
	assert.Len(t, sot1.generalizations(), 1)
	assert.Len(t, sot2.generalizations(), 1)
	assert.Len(t, sotG.specifications(), 2)

	sot1.generalize(sotG)
	assert.Len(t, sot1.generalizations(), 1)
	assert.Len(t, sot2.generalizations(), 1)
	assert.Len(t, sotG.specifications(), 2)

	assert.Equal(t, sot1, sot1.lowestCommonGeneralization(sot1))
	assert.Equal(t, sotG, sot1.lowestCommonGeneralization(sotG))
	assert.Equal(t, sotG, sotG.lowestCommonGeneralization(sot1))
}
