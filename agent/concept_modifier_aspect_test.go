package agent

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestAspectModifierTypeConstructor(t *testing.T) {
	agent := newEmptyWorldAgent()
	info1, info2 := "info1", "info2"
	info := []string{info1, info2}
	amt := agent.newAspectModifierType(agent.aspect.find(info...))
	assert.Equal(t, amt, agent.memory.find(amt))
	assert.Equal(t, amt.aspect, agent.aspect.find(info...))

	amtCopy := agent.newAspectModifierType(agent.aspect.find(info...))
	assert.Equal(t, amt, amtCopy)
}

func TestAspectModifierTypeToString(t *testing.T) {
	agent := newEmptyWorldAgent()
	info1, info2 := "info1", "info2"
	info := []string{info1, info2}
	amt := agent.newAspectModifierType(agent.aspect.find(info...))
	amtStr := amt.debug("", 2)
	for _, infoStr := range info {
		assert.Contains(t, amtStr, infoStr)
	}
}

func TestAspectModifierTypeInstantiate(t *testing.T) {
	agent := newEmptyWorldAgent()
	info1, info2 := "info1", "info2"
	info := []string{info1, info2}
	amt := agent.newAspectModifierType(agent.aspect.find(info...))
	tc := agent.newTestConcept(1)
	assert.Empty(t, amt.sources())
	am := amt.instantiate(tc, conceptSourceObservation)
	assert.Len(t, amt.sources(), 1)
	assert.Contains(t, amt.sources(), conceptSourceObservation)
	assert.Equal(t, amt, am._type())
	assert.Equal(t, tc, am.target())
	assert.Equal(t, conceptSourceObservation, am.source())
	assert.Equal(t, am, amt.instantiate(tc, conceptSourceObservation))
	amStr := am.debug("", 2)
	for _, infoStr := range info {
		assert.Contains(t, amStr, infoStr)
	}
	assert.Contains(t, amStr, modifierSourceNames[conceptSourceObservation])

	val := 123
	valuedAm := amt.instantiate(tc, conceptSourceObservation, val)
	assert.NotEqual(t, am, valuedAm)
	amStr = valuedAm.debug("", 2)
	assert.Contains(t, amStr, strconv.Itoa(val))
}

func TestAspectModifierOverride(t *testing.T) {
	agent := newEmptyWorldAgent()
	info11, info12, info21, info22 := "info11", "info12", "info21", "info22"
	infoList11, infoList12, infoList21 := []string{info11, info21}, []string{info11, info22}, []string{info12, info21}
	amt11 := agent.newAspectModifierType(agent.aspect.find(infoList11...))
	amt12 := agent.newAspectModifierType(agent.aspect.find(infoList12...))
	amt21 := agent.newAspectModifierType(agent.aspect.find(infoList21...))
	tc := agent.newTestConcept(1)

	assert.Len(t, tc.modifiers(), 0)
	am11 := amt11.instantiate(tc, conceptSourceObservation)
	assert.Len(t, tc.modifiers(), 1)
	assert.Equal(t, tc.modifiers()[am11.id()], am11)
	am12 := amt12.instantiate(tc, conceptSourceObservation)
	assert.Len(t, tc.modifiers(), 1)
	assert.Equal(t, tc.modifiers()[am12.id()], am12)
	am21 := amt21.instantiate(tc, conceptSourceObservation)
	assert.Len(t, tc.modifiers(), 2)
	assert.Equal(t, tc.modifiers()[am12.id()], am12)
	assert.Equal(t, tc.modifiers()[am21.id()], am21)
}

func TestAspectModifierOverrideErrorHandling(t *testing.T) {
	agent := newEmptyWorldAgent()
	info1, info2 := "info1", "info2"
	info := []string{info1, info2}
	amt := agent.newAspectModifierType(agent.aspect.find(info...))
	tc1, tc2 := agent.newTestConcept(1), agent.newTestConcept(2)
	am1, am2 := amt.instantiate(tc1, conceptSourceObservation), amt.instantiate(tc2, conceptSourceObservation)
	assert.Nil(t, am1.override(am2))
}

func TestAspectModifierTypeVerify(t *testing.T) {
	agent := newEmptyWorldAgent()
	info11, info12, info21, info22 := "info11", "info12", "info21", "info22"
	infoList11, infoList12, infoList21 := []string{info11, info21}, []string{info11, info22}, []string{info12, info21}
	amt11 := agent.newAspectModifierType(agent.aspect.find(infoList11...))
	amt12 := agent.newAspectModifierType(agent.aspect.find(infoList12...))
	amt21 := agent.newAspectModifierType(agent.aspect.find(infoList21...))
	tc := agent.newTestConcept(1)
	amt11.instantiate(tc, conceptSourceObservation)
	assert.Nil(t, amt11.verify())
	amt11.lockMap[partIdModifierTarget] = tc
	amt12.lockMap[partIdModifierTarget] = tc
	amt21.lockMap[partIdModifierTarget] = tc
	assert.True(t, *amt11.verify())
	assert.False(t, *amt12.verify())
	assert.Nil(t, amt21.verify())
}

func TestAspectModifierShareSync(t *testing.T) {
	agent := newEmptyWorldAgent()
	info1, info2 := "info1", "info2"
	info := []string{info1, info2}
	amt := agent.newAspectModifierType(agent.aspect.find(info...))
	tc := agent.newTestConcept(1)
	am := amt.instantiate(tc, conceptSourceObservation)
	sync := am.instShareParts()
	assert.Len(t, sync, 1)
	assert.Equal(t, partIdModifierTarget, sync[tc.id()])
}
