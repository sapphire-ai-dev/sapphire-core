package agent

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConceptCtxMatch(t *testing.T) {
	agent := newEmptyWorldAgent()
	tc := agent.newTestConcept(1, nil)

	ccat := agent.newCreateContextActionType()
	cca := agent.newCreateContextAction(ccat, agent.self, 0)
	cot := agent.newContextObjectType(conceptSourceObservation)
	co := agent.newContextObject(cca)
	co.addType(cot)

	assert.Equal(t, tc, agent.newTestConcept(1, nil))
	assert.NotEqual(t, tc, agent.newTestConcept(1, map[int]any{partIdConceptContext: co}))
}

func TestDecoratorGetRelations(t *testing.T) {
	agent := newEmptyWorldAgent()
	so := agent.newSimpleObject(map[int]any{partIdObjectWorldId: 123})
	amt1 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "red"), nil)
	amt2 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "blue"), nil)
	t0, t1, t2, t3, t4 := 0, 1, 2, 3, 4
	tpo0 := agent.newTimePointObject(&t0, nil)
	tpo1 := agent.newTimePointObject(&t1, nil)
	tpo2 := agent.newTimePointObject(&t2, nil)
	tpo3 := agent.newTimePointObject(&t3, nil)
	tpo4 := agent.newTimePointObject(&t4, nil)
	tso01 := agent.newTimeSegmentObject(tpo0, tpo1, nil)
	tso02 := agent.newTimeSegmentObject(tpo0, tpo2, nil)
	tso04 := agent.newTimeSegmentObject(tpo0, tpo4, nil)
	tso12 := agent.newTimeSegmentObject(tpo1, tpo2, nil)
	tso13 := agent.newTimeSegmentObject(tpo1, tpo3, nil)
	tso23 := agent.newTimeSegmentObject(tpo2, tpo3, nil)
	tso24 := agent.newTimeSegmentObject(tpo2, tpo4, nil)
	tso34 := agent.newTimeSegmentObject(tpo3, tpo4, nil)
	am1 := amt1.instantiate(so, conceptSourceObservation, map[int]any{partIdConceptTime: tso12})
	am2 := amt2.instantiate(so, conceptSourceObservation, map[int]any{partIdConceptTime: tso23})
	assert.NotEqual(t, am1, am2)
	assert.Len(t, so.modifiers(nil), 2)
	assert.Empty(t, so.modifiers(map[int]any{partIdConceptTime: tpo0}))
	assert.Len(t, so.modifiers(map[int]any{partIdConceptTime: tpo1}), 1)
	assert.Len(t, so.modifiers(map[int]any{partIdConceptTime: tpo2}), 2)
	assert.Len(t, so.modifiers(map[int]any{partIdConceptTime: tpo3}), 1)
	assert.Empty(t, so.modifiers(map[int]any{partIdConceptTime: tpo4}))
	assert.Len(t, so.modifiers(map[int]any{partIdConceptTime: tso01}), 1)
	assert.Len(t, so.modifiers(map[int]any{partIdConceptTime: tso02}), 2)
	assert.Len(t, so.modifiers(map[int]any{partIdConceptTime: tso04}), 2)
	assert.Len(t, so.modifiers(map[int]any{partIdConceptTime: tso13}), 2)
	assert.Len(t, so.modifiers(map[int]any{partIdConceptTime: tso24}), 2)
	assert.Len(t, so.modifiers(map[int]any{partIdConceptTime: tso34}), 1)
}
