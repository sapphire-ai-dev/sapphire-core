package agent

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestAuxiliaryRelationTypeConstructor(t *testing.T) {
    agent := newEmptyWorldAgent()
    sot := agent.newSimpleObjectType(conceptSourceObservation, map[int]modifierType{}, nil)
    aat := agent.newAtomicActionType(newTestActionInterface().instantiate(), nil)
    art := agent.newAuxiliaryRelationType(auxiliaryTypeIdWant, false, sot, aat, nil)
    assert.Equal(t, art.auxiliaryTypeId, auxiliaryTypeIdWant)
    assert.Equal(t, art.lType.c, sot)
    assert.Equal(t, art.rType.c, aat)

    art2 := agent.newAuxiliaryRelationType(auxiliaryTypeIdWant, false, nil, nil, nil)
    assert.Nil(t, art2.lType)
    assert.Nil(t, art2.rType)
}

func TestAuxiliaryRelationConstructor(t *testing.T) {
    agent := newEmptyWorldAgent()
    tai := newTestActionInterface()
    tai.ReadyResult, tai.StepCount = true, 0
    aat := agent.newAtomicActionType(tai.instantiate(), nil)
    sot := agent.newSimpleObjectType(conceptSourceObservation, map[int]modifierType{}, nil)
    agent.self.addType(sot)
    art := agent.newAuxiliaryRelationType(auxiliaryTypeIdWant, false, sot, aat, nil)

    co := newTestContext(agent, 0)
    aa := agent.newAtomicAction(aat, agent.self, map[int]any{conceptArgContext: co})

    asct := agent.newActionStateChangeType(aat, nil)
    asct.addValue(10.0)
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
    artP := agent.newAuxiliaryRelationType(auxiliaryTypeIdWant, false, nil, nil, nil)
    artN := agent.newAuxiliaryRelationType(auxiliaryTypeIdWant, true, nil, nil, nil)
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
    artP := agent.newAuxiliaryRelationType(auxiliaryTypeIdWant, false, nil, nil, nil)
    artN := agent.newAuxiliaryRelationType(auxiliaryTypeIdWant, true, nil, nil, nil)
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

func TestAuxiliaryRelationWantCancel(t *testing.T) {
    agent := newEmptyWorldAgent()
    tai := newTestActionInterface()
    tai.ReadyResult, tai.StepCount = true, 0
    aat := agent.newAtomicActionType(tai.instantiate(), nil)
    sot := agent.newSimpleObjectType(conceptSourceObservation, map[int]modifierType{}, nil)
    agent.self.addType(sot)
    art := agent.newAuxiliaryRelationType(auxiliaryTypeIdWant, false, sot, aat, nil)
    co := newTestContext(agent, 0)
    aa := agent.newAtomicAction(aat, agent.self, map[int]any{conceptArgContext: co})
    asct := agent.newActionStateChangeType(aat, nil)
    asct.addValue(10.0)
    asc := agent.newActionStateChange(asct, aa, nil)
    ar := agent.newAuxiliaryRelation(art, agent.self, aa, map[int]any{
        conceptArgRelationAuxiliaryWantChange: asc,
    })

    ar.interpret()
    assert.True(t, tai.ReadyResult)
    assert.Zero(t, tai.StepCount)

    agent.mind.add(aat)
    agent.cycle()
    assert.True(t, tai.ReadyResult)
    assert.Equal(t, tai.StepCount, 1)

    agent.cycle()
    assert.True(t, tai.ReadyResult)
    assert.Equal(t, tai.StepCount, 2)

    assert.Nil(t, ar.time())
    artN := agent.newAuxiliaryRelationType(auxiliaryTypeIdWant, true, sot, aat, nil)
    arN := agent.newAuxiliaryRelation(artN, agent.self, aa, map[int]any{
        conceptArgTime: agent.newTimeSegmentObject(agent.time.now.start(), nil, nil),
    })

    assert.Equal(t, ar.time().end(), arN.time().start())
    agent.cycle()
    assert.True(t, tai.ReadyResult)
    assert.Equal(t, tai.StepCount, 2)
}
