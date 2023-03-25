package agent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSequentialActionTypeConstructor(t *testing.T) {
	agent := newEmptyWorldAgent()
	tai1, tai2 := newTestActionInterface(), newTestActionInterface()
	ai1, ai2 := tai1.instantiate(), tai2.instantiate()
	aat1, aat2 := agent.newAtomicActionType(ai1), agent.newAtomicActionType(ai2)
	sot := agent.newSimpleObjectType(conceptSourceObservation, map[int]modifierType{})
	sat := agent.newSequentialActionType(sot, aat1, aat2)
	assert.Equal(t, aat1, sat.first())
	assert.Equal(t, aat2, sat.next())
	assert.Equal(t, sot, sat.receiverType())
	assert.Equal(t, sat, agent.newSequentialActionType(sot, aat1, aat2))
}

func TestSequentialActionTypeInstantiate(t *testing.T) {
	agent := newEmptyWorldAgent()
	tai1, tai2 := newTestActionInterface(), newTestActionInterface()
	ai1, ai2 := tai1.instantiate(), tai2.instantiate()
	aat1, aat2 := agent.newAtomicActionType(ai1), agent.newAtomicActionType(ai2)
	sot := agent.newSimpleObjectType(conceptSourceObservation, map[int]modifierType{})
	sat := agent.newSequentialActionType(sot, aat1, aat2)
	sa := sat.instantiate()
	assert.Equal(t, sa, sat.instantiate())
	assert.Equal(t, sat, sa.part(partIdActionT))
	assert.Equal(t, agent.self, sa.part(partIdActionPerformer))
	aa1, aa2 := sa.part(partIdActionSequentialFirst).(action), sa.part(partIdActionSequentialNext).(action)
	assert.Equal(t, aat1, aa1._type())
	assert.Equal(t, aat2, aa2._type())

	so := agent.newSimpleObject(1)
	sa.setReceiver(so)
	assert.Equal(t, so, sa.receiver())
	sa.setReceiver(nil)
	assert.Equal(t, so, sa.receiver())
}

func TestSequentialActionDebug(t *testing.T) {
	agent := newEmptyWorldAgent()
	tai1, tai2 := newTestActionInterface(), newTestActionInterface()
	ai1, ai2 := tai1.instantiate(), tai2.instantiate()
	aat1, aat2 := agent.newAtomicActionType(ai1), agent.newAtomicActionType(ai2)
	sat := agent.newSequentialActionType(nil, aat1, aat2)
	sa := sat.instantiate()
	assert.Contains(t, sat.debug("", 2), toReflect[*atomicActionType]().Name())
	assert.Contains(t, sa.debug("", 2), toReflect[*atomicAction]().Name())
}

func TestSequentialActionLifecycle(t *testing.T) {
	agent := newEmptyWorldAgent()
	tai1, tai2 := newTestActionInterface(), newTestActionInterface()
	ai1, ai2 := tai1.instantiate(), tai2.instantiate()
	aat1, aat2 := agent.newAtomicActionType(ai1), agent.newAtomicActionType(ai2)
	sat := agent.newSequentialActionType(nil, aat1, aat2)
	sa := sat.instantiate().(*sequentialAction)
	aa1, aa2 := sa.first(), sa.next()

	assert.Equal(t, sa.state(), actionStateIdle)
	assert.False(t, sa.step())
	assert.Equal(t, sa.state(), actionStateIdle)
	assert.False(t, sa.start())
	assert.Equal(t, sa.state(), actionStateIdle)
	assert.Equal(t, aa1.state(), actionStateIdle)
	assert.Equal(t, aa2.state(), actionStateIdle)

	tai1.ReadyResult = true
	assert.True(t, sa.start())
	assert.Equal(t, sa.state(), actionStateActive)
	assert.Equal(t, aa1.state(), actionStateActive)
	assert.False(t, sa.start())
	assert.Equal(t, sa.state(), actionStateActive)
	assert.Equal(t, aa1.state(), actionStateActive)
	assert.Zero(t, tai1.StepCount)

	tai1.ReadyResult = false
	assert.False(t, sa.step())
	tai1.ReadyResult = true
	assert.True(t, sa.step())
	assert.Equal(t, sa.state(), actionStateActive)
	assert.Equal(t, aa1.state(), actionStateDone)
	assert.Equal(t, aa2.state(), actionStateIdle)
	assert.Equal(t, 1, tai1.StepCount)
	assert.Equal(t, 0, tai2.StepCount)

	assert.False(t, sa.step())
	tai2.ReadyResult = true
	assert.Equal(t, aa2.state(), actionStateIdle)
	tai2.ReadyResult = false
	assert.False(t, sa.step())
	tai2.ReadyResult = true
	assert.Equal(t, aa2.state(), actionStateIdle)
	assert.True(t, sa.step())
	assert.Equal(t, aa2.state(), actionStateDone)
	assert.Equal(t, sa.state(), actionStateDone)

	assert.False(t, sa.step())
	assert.False(t, sa.start())
	assert.Equal(t, 1, tai1.StepCount)
	assert.Equal(t, 1, tai2.StepCount)
}

func TestSequentialActionLifecycleThreeStep(t *testing.T) {
	agent := newEmptyWorldAgent()
	tai1, tai2, tai3 := newTestActionInterface(), newTestActionInterface(), newTestActionInterface()
	ai1, ai2, ai3 := tai1.instantiate(), tai2.instantiate(), tai3.instantiate()
	aat1, aat2, aat3 := agent.newAtomicActionType(ai1), agent.newAtomicActionType(ai2), agent.newAtomicActionType(ai3)
	sat := agent.newSequentialActionType(nil, aat1, agent.newSequentialActionType(nil, aat2, aat3))
	sa := sat.instantiate()
	tai1.ReadyResult = true
	tai2.ReadyResult = true
	tai3.ReadyResult = false
	assert.True(t, sa.start())
	assert.True(t, sa.step())
	assert.True(t, sa.step())
	assert.False(t, sa.step())
	tai3.ReadyResult = true
	assert.True(t, sa.step())
	assert.Equal(t, sa.state(), actionStateDone)
	aa1 := sa.part(partIdActionSequentialFirst).(performableAction)
	sa2 := sa.part(partIdActionSequentialNext).(performableAction)
	aa2 := sa2.part(partIdActionSequentialFirst).(performableAction)
	aa3 := sa2.part(partIdActionSequentialNext).(performableAction)
	assert.Equal(t, aa1.state(), actionStateDone)
	assert.Equal(t, sa2.state(), actionStateDone)
	assert.Equal(t, aa2.state(), actionStateDone)
	assert.Equal(t, aa3.state(), actionStateDone)
}
