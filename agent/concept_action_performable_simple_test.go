package agent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleActionTypeConstructor(t *testing.T) {
	agent := newEmptyWorldAgent()
	aat := agent.newAtomicActionType(newTestActionInterface().instantiate(), nil)
	sot := agent.newSimpleObjectType(conceptSourceObservation, map[int]modifierType{}, nil)
	sat := agent.newSimpleActionType(sot, aat, nil)
	assert.Equal(t, aat, sat.child())
	assert.Equal(t, sat, agent.newSimpleActionType(sot, aat, nil))

	satNoReceiver := agent.newSimpleActionType(nil, aat, nil)
	assert.NotEqual(t, sat, satNoReceiver)
}

func TestSimpleActionTypeInstantiate(t *testing.T) {
	agent := newEmptyWorldAgent()
	aat := agent.newAtomicActionType(newTestActionInterface().instantiate(), nil)
	sot := agent.newSimpleObjectType(conceptSourceObservation, map[int]modifierType{}, nil)
	sat := agent.newSimpleActionType(sot, aat, nil)
	var sa *simpleAction
	for _, inst := range sat.instantiate(nil) {
		sa = inst.(*simpleAction)
	}
	assert.Equal(t, sat, sa.part(partIdActionT))
	assert.Equal(t, agent.self, sa.part(partIdActionPerformer))
	assert.Nil(t, sa.part(partIdActionReceiver))

	aa := sa.part(partIdActionSimpleChild).(action)
	assert.Equal(t, aat, aa._type())
	assert.Nil(t, aa.part(partIdActionReceiver))

	so := agent.newSimpleObject(1, nil)
	sa.setReceiver(so)
	assert.Equal(t, so, sa.part(partIdActionReceiver))
	assert.Equal(t, so, aa.part(partIdActionReceiver))

	sa.setReceiver(nil)
	assert.Equal(t, so, sa.part(partIdActionReceiver))
	assert.Equal(t, so, aa.part(partIdActionReceiver))
}

func TestSimpleActionTypeDebug(t *testing.T) {
	agent := newEmptyWorldAgent()
	aat := agent.newAtomicActionType(newTestActionInterface().instantiate(), nil)
	sat := agent.newSimpleActionType(nil, aat, nil)
	var sa *simpleAction
	for _, inst := range sat.instantiate(nil) {
		sa = inst.(*simpleAction)
	}
	assert.Contains(t, sat.debug("", 2), toReflect[*atomicActionType]().Name())
	assert.Contains(t, sa.debug("", 2), toReflect[*atomicAction]().Name())
}

func TestSimpleActionLifecycle(t *testing.T) {
	agent := newEmptyWorldAgent()
	tai := newTestActionInterface()
	aat := agent.newAtomicActionType(tai.instantiate(), nil)
	sat := agent.newSimpleActionType(nil, aat, nil)
	var sa *simpleAction
	for _, inst := range sat.instantiate(nil) {
		sa = inst.(*simpleAction)
	}

	assert.Equal(t, sa.state(), actionStateIdle)
	assert.False(t, sa.step())
	assert.Equal(t, sa.state(), actionStateIdle)
	assert.False(t, sa.start())
	assert.Equal(t, sa.state(), actionStateIdle)

	tai.ReadyResult = true
	assert.True(t, sa.start())
	assert.Equal(t, sa.state(), actionStateActive)
	assert.False(t, sa.start())
	assert.Equal(t, sa.state(), actionStateActive)
	assert.Zero(t, tai.StepCount)

	tai.ReadyResult = false
	assert.False(t, sa.step())
	tai.ReadyResult = true
	assert.True(t, sa.step())
	assert.Equal(t, sa.state(), actionStateDone)
	assert.Equal(t, 1, tai.StepCount)

	assert.False(t, sa.start())
	assert.False(t, sa.step())
}
