package agent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConceptCpntSyncInstShareParts(t *testing.T) {
	agent := newEmptyWorldAgent()
	tc := agent.newTestConcept(1)
	assert.NotNil(t, tc.instShareParts())

	amt := agent.newAspectModifierType(agent.aspect.find([]string{"info1", "info2"}...))
	am := amt.instantiate(tc, conceptSourceObservation)
	assert.Equal(t, partIdModifierTarget, am.instShareParts()[tc.id()])
}

func TestConceptCpntSyncTypeUpdateSync(t *testing.T) {
	agent := newEmptyWorldAgent()
	amt := agent.newAspectModifierType(agent.aspect.find([]string{"info1", "info2"}...))
	aat := agent.newAtomicActionType(newTestActionInterface().instantiate())
	assert.Empty(t, amt.syncMap)
	assert.Empty(t, aat.syncMap)

	so := agent.newSimpleObject(1)
	am := amt.instantiate(so, conceptSourceObservation)
	aa := aat.instantiate()
	aa.setReceiver(so)
	amt.typeUpdateSync(aat, am.instShareParts(), aa.instShareParts())
	assert.Len(t, amt.syncMap, 1)
	assert.Len(t, aat.syncMap, 1)
	assert.Equal(t, 1, amt.syncMap[aat.id()].data[partIdModifierTarget][partIdActionReceiver])
	assert.Equal(t, 1, aat.syncMap[amt.id()].data[partIdActionReceiver][partIdModifierTarget])

	amt.typeUpdateSync(aat, am.instShareParts(), aa.instShareParts())
	assert.Len(t, amt.syncMap, 1)
	assert.Len(t, aat.syncMap, 1)
	assert.Equal(t, 2, amt.syncMap[aat.id()].data[partIdModifierTarget][partIdActionReceiver])
	assert.Equal(t, 2, aat.syncMap[amt.id()].data[partIdActionReceiver][partIdModifierTarget])
}

func TestConceptCpntSyncTypeLockSync(t *testing.T) {
	agent := newEmptyWorldAgent()
	amt := agent.newAspectModifierType(agent.aspect.find([]string{"info1", "info2"}...))
	aat := agent.newAtomicActionType(newTestActionInterface().instantiate())
	so := agent.newSimpleObject(1)
	am := amt.instantiate(so, conceptSourceObservation)
	aa := aat.instantiate()
	aa.setReceiver(so)
	amt.typeUpdateSync(aat, am.instShareParts(), aa.instShareParts())
	assert.Empty(t, amt.lockMap)
	assert.Empty(t, aat.lockMap)

	amt.typeLockSync(aat, map[int]concept{partIdActionReceiver: aa.receiver()})
	assert.Len(t, amt.lockMap, 1)
	assert.Equal(t, aa.receiver(), amt.lockMap[partIdModifierTarget])
	assert.Empty(t, aat.lockMap)

	amt.typeUnlockSync()
	assert.Empty(t, amt.lockMap)
}
