package agent

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
	"github.com/stretchr/testify/assert"
	"testing"
)

func generateTime(agent *Agent, start int, end int) ([]*timePointObject, [][]*timeSegmentObject) {
	tpos := make([]*timePointObject, end)
	tsos := make([][]*timeSegmentObject, end)
	for i := start; i < end; i++ {
		tpos[i] = agent.newTimePointObject(&i, nil)
	}

	for i := start; i < end-1; i++ {
		tsos[i] = make([]*timeSegmentObject, end)
		for j := i + 1; j < end; j++ {
			tsos[i][j] = agent.newTimeSegmentObject(tpos[i], tpos[j], nil).(*timeSegmentObject)
		}
	}

	return tpos, tsos
}

func assertAspectModifierSeen(t *testing.T, modifs map[int]modifier, amt *aspectModifierType,
	target concept) *aspectModifier {
	var result *aspectModifier

	for _, modif := range modifs {
		if am, ok := modif.(*aspectModifier); ok && am._type() == amt {
			result = am
			assert.Equal(t, target, am.target())
		}
	}

	assert.NotNil(t, result)
	return result
}

func TestVersioningOverlapSingleSide(t *testing.T) {
	agent := newEmptyWorldAgent()
	so := agent.newSimpleObject(123, nil)
	amt1 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "red"), nil)
	amt2 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "blue"), nil)
	tpos, tsos := generateTime(agent, 0, 6)
	amt1.instantiate(so, conceptSourceObservation, map[int]any{conceptArgTime: tsos[1][4]})
	am1 := assertAspectModifierSeen(t, so.modifiers(nil), amt1, so)
	assert.Equal(t, am1.time().start(), tpos[1])
	assert.Equal(t, am1.time().end(), tpos[4])

	amt2.instantiate(so, conceptSourceObservation, map[int]any{conceptArgTime: tsos[0][2]})
	am1 = assertAspectModifierSeen(t, so.modifiers(nil), amt1, so)
	am2 := assertAspectModifierSeen(t, so.modifiers(nil), amt2, so)
	assert.Equal(t, am1.time().start(), tpos[2])
	assert.Equal(t, am1.time().end(), tpos[4])
	assert.Equal(t, am2.time().start(), tpos[0])
	assert.Equal(t, am2.time().end(), tpos[2])

	amt2.instantiate(so, conceptSourceObservation, map[int]any{conceptArgTime: tsos[3][5]})
	am1 = assertAspectModifierSeen(t, so.modifiers(nil), amt1, so)
	assert.Equal(t, am1.time().start(), tpos[2])
	assert.Equal(t, am1.time().end(), tpos[3])
	assert.Len(t, so.modifiers(nil), 3)
}

func TestVersioningOverlapSingleSideOriginalNils(t *testing.T) {
	agent := newEmptyWorldAgent()
	so1 := agent.newSimpleObject(123, nil)
	amt1 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "red"), nil)
	amt2 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "blue"), nil)
	tpos, tsos := generateTime(agent, 0, 6)
	amt1.instantiate(so1, conceptSourceObservation, map[int]any{
		conceptArgTime: agent.newTimeSegmentObject(nil, tpos[3], nil),
	})
	am1 := assertAspectModifierSeen(t, so1.modifiers(nil), amt1, so1)
	assert.Nil(t, am1.time().start())
	assert.Equal(t, am1.time().end(), tpos[3])

	amt2.instantiate(so1, conceptSourceObservation, map[int]any{conceptArgTime: tsos[0][4]})
	am1 = assertAspectModifierSeen(t, so1.modifiers(nil), amt1, so1)
	assert.Nil(t, am1.time().start())
	assert.Equal(t, am1.time().end(), tpos[0])

	so2 := agent.newSimpleObject(234, nil)
	amt3 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "yellow"), nil)
	amt3.instantiate(so2, conceptSourceObservation, map[int]any{
		conceptArgTime: agent.newTimeSegmentObject(tpos[2], nil, nil),
	})
	am3 := assertAspectModifierSeen(t, so2.modifiers(nil), amt3, so2)
	assert.Equal(t, am3.time().start(), tpos[2])
	assert.Nil(t, am3.time().end())

	amt2.instantiate(so2, conceptSourceObservation, map[int]any{conceptArgTime: tsos[0][4]})
	am3 = assertAspectModifierSeen(t, so2.modifiers(nil), amt3, so2)
	assert.Equal(t, am3.time().start(), tpos[4])
	assert.Nil(t, am3.time().end())
}
