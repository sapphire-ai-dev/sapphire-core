package agent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func newContextObject(agent *Agent) *contextObject {
	ccat := agent.newCreateContextActionType()
	cca := agent.newCreateContextAction(ccat, agent.self)
	cot := agent.newContextObjectType(conceptSourceObservation)
	co := agent.newContextObject(cca)
	co.addType(cot)
	return co
}

func TestAuxiliaryRelationTypeConstructor(t *testing.T) {
	agent := newEmptyWorldAgent()
	art := agent.newAuxiliaryRelationType(auxiliaryTypeIdWant, nil)
	assert.Equal(t, art.auxiliaryTypeId, auxiliaryTypeIdWant)
}

func TestAuxiliaryRelationConstructor(t *testing.T) {
	agent := newEmptyWorldAgent()
	art := agent.newAuxiliaryRelationType(auxiliaryTypeIdWant, nil)

	tai := newTestActionInterface()
	tai.ReadyResult, tai.StepCount = true, 0
	aat := agent.newAtomicActionType(tai.instantiate(), nil)
	co := newContextObject(agent)
	aa := agent.newAtomicAction(aat, agent.self, map[int]any{co.cid: co})

	asct := agent.newActionStateChangeType(aat, nil)
	asc := agent.newActionStateChange(asct, aa, nil)
	ar := agent.newAuxiliaryRelation(art, agent.self, aa, map[int]any{
		conceptArgRelationAuxiliaryWantChange: asc,
	})

	assert.Equal(t, art, ar._type())
	assert.Equal(t, agent.self, ar.lObject())
	assert.Equal(t, aa, ar.rTarget())
	assert.Equal(t, asc, ar.wantChange())

	assert.True(t, tai.ReadyResult)
	assert.Zero(t, tai.StepCount)

	agent.mind.add(aat)
	agent.cycle()

	assert.True(t, tai.ReadyResult)
	assert.Equal(t, tai.StepCount, 1)
}
