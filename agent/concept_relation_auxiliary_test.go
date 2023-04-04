package agent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuxiliaryRelationTypeConstructor(t *testing.T) {
	agent := newEmptyWorldAgent()
	art := agent.newAuxiliaryRelationType(auxiliaryTypeIdWant, false, nil)
	assert.Equal(t, art.auxiliaryTypeId, auxiliaryTypeIdWant)
}

func TestAuxiliaryRelationConstructor(t *testing.T) {
	agent := newEmptyWorldAgent()
	art := agent.newAuxiliaryRelationType(auxiliaryTypeIdWant, false, nil)

	tai := newTestActionInterface()
	tai.ReadyResult, tai.StepCount = true, 0
	aat := agent.newAtomicActionType(tai.instantiate(), nil)
	co := newTestContext(agent, 0)
	aa := agent.newAtomicAction(aat, agent.self, map[int]any{conceptArgContext: co})

	asct := agent.newActionStateChangeType(aat, nil)
	asc := agent.newActionStateChange(asct, aa, nil)
	ar := agent.newAuxiliaryRelation(art, agent.self, aa, map[int]any{
		conceptArgRelationAuxiliaryWantChange: asc,
	})

	ar.interpret()
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

func TestAuxiliaryRelationVersioning(t *testing.T) {
	agent := newEmptyWorldAgent()
	artP := agent.newAuxiliaryRelationType(auxiliaryTypeIdWant, false, nil)
	artN := agent.newAuxiliaryRelationType(auxiliaryTypeIdWant, true, nil)
	aat := agent.newAtomicActionType(newTestActionInterface().instantiate(), nil)
	co := newTestContext(agent, 0)
	aa := agent.newAtomicAction(aat, agent.self, map[int]any{conceptArgContext: co})
	asct := agent.newActionStateChangeType(aat, nil)
	asc := agent.newActionStateChange(asct, aa, nil)
	tpos, tsos := generateTime(agent, 0, 6)
	arT := agent.newAuxiliaryRelation(artP, agent.self, aa, map[int]any{
		conceptArgTime:                        tsos[2][4],
		conceptArgRelationAuxiliaryWantChange: asc,
	})
	assert.Equal(t, arT.time().start(), tpos[2])
	assert.Equal(t, arT.time().end(), tpos[4])

	arN := agent.newAuxiliaryRelation(artN, agent.self, aa, map[int]any{
		conceptArgTime: tsos[3][5],
	})
	assert.Equal(t, arT.time().start(), tpos[2])
	assert.Equal(t, arT.time().end(), tpos[3])
	assert.Equal(t, arN.time().start(), tpos[3])
	assert.Equal(t, arN.time().end(), tpos[5])
}

func TestAuxiliaryRelationVersioningInterrupt(t *testing.T) {
	agent := newEmptyWorldAgent()
	artP := agent.newAuxiliaryRelationType(auxiliaryTypeIdWant, false, nil)
	artN := agent.newAuxiliaryRelationType(auxiliaryTypeIdWant, true, nil)
	aat := agent.newAtomicActionType(newTestActionInterface().instantiate(), nil)
	co := newTestContext(agent, 0)
	aa := agent.newAtomicAction(aat, agent.self, map[int]any{conceptArgContext: co})
	asct := agent.newActionStateChangeType(aat, nil)
	asc := agent.newActionStateChange(asct, aa, nil)
	tpos, tsos := generateTime(agent, 0, 6)
	agent.newAuxiliaryRelation(artP, agent.self, aa, map[int]any{
		conceptArgTime:                        tsos[2][5],
		conceptArgRelationAuxiliaryWantChange: asc,
	})
	assert.Len(t, aa.relations(nil), 1)
	agent.newAuxiliaryRelation(artN, agent.self, aa, map[int]any{
		conceptArgTime: tsos[3][4],
	})
	assert.Len(t, aa.relations(nil), 4)
	assert.Len(t, aa.relations(map[int]any{conceptArgTime: tpos[5]}), 2)
}
