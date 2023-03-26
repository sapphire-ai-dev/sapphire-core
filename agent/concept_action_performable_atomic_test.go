package agent

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestAtomicActionTypeConstructor(t *testing.T) {
    agent := newEmptyWorldAgent()
    tai := newTestActionInterface()
    ai := tai.instantiate()
    aat := agent.newAtomicActionType(ai, nil)
    assert.Equal(t, ai, aat.actionInterface)
    assert.Equal(t, aat, agent.newAtomicActionType(ai, nil))

    assert.NotNil(t, aat.conditions())
    assert.NotNil(t, aat.causations())
    assert.Nil(t, aat.receiverType())
}

func TestAtomicActionTypeInstantiate(t *testing.T) {
    agent := newEmptyWorldAgent()
    aat := agent.newAtomicActionType(newTestActionInterface().instantiate(), nil)
    aa := aat.instantiate()
    assert.Equal(t, aa, aat.instantiate())
    assert.Equal(t, aat, aa.part(partIdActionT))
    assert.Equal(t, agent.self, aa.part(partIdActionPerformer))
    assert.Nil(t, aa.part(partIdActionReceiver))
    assert.Nil(t, aa.part(partIdStart))
}

func TestAtomicActionDebug(t *testing.T) {
    agent := newEmptyWorldAgent()
    aat := agent.newAtomicActionType(newTestActionInterface().instantiate(), nil)
    aa := aat.instantiate()
    assert.Contains(t, aat.debug("", 1), toReflect[*TestActionInterface]().Name())
    assert.Contains(t, aa.debug("", 2), toReflect[*TestActionInterface]().Name())
}

func TestAtomicActionLifecycle(t *testing.T) {
    agent := newEmptyWorldAgent()
    tai := newTestActionInterface()
    aat := agent.newAtomicActionType(tai.instantiate(), nil)
    aa := aat.instantiate()

    assert.Equal(t, aa.state(), actionStateIdle)
    assert.False(t, aa.step())
    assert.Equal(t, aa.state(), actionStateIdle)
    assert.False(t, aa.start())
    assert.Equal(t, aa.state(), actionStateIdle)

    tai.ReadyResult = true
    assert.True(t, aa.start())
    assert.Equal(t, aa.state(), actionStateActive)
    assert.False(t, aa.start())
    assert.Equal(t, aa.state(), actionStateActive)
    assert.Zero(t, tai.StepCount)

    tai.ReadyResult = false
    assert.False(t, aa.step())
    tai.ReadyResult = true
    assert.True(t, aa.step())
    assert.Equal(t, aa.state(), actionStateDone)
    assert.Equal(t, 1, tai.StepCount)

    assert.False(t, aa.start())
    assert.False(t, aa.step())
}

func TestAtomicActionInterpret(t *testing.T) {
    agent := newEmptyWorldAgent()
    aat := agent.newAtomicActionType(newTestActionInterface().instantiate(), nil)
    so1, so2 := agent.newSimpleObject(1, nil), agent.newSimpleObject(2, nil)
    concepts := map[int]concept{
        partIdActionT:         aat,
        partIdActionPerformer: so1,
        partIdActionReceiver:  so2,
    }
    aa := agent.interpretAtomicAction(concepts).(*atomicAction)
    for partId, c := range concepts {
        assert.Equal(t, c, aa.part(partId))
    }

    aa.setReceiver(nil)
    assert.Equal(t, so2, aa.receiver())

    assert.NotEqual(t, aa, agent.interpretAtomicAction(concepts))
}
