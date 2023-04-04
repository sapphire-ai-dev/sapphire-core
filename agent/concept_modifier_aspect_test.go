package agent

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestAspectModifierTypeConstructor(t *testing.T) {
	agent := newEmptyWorldAgent()
	info1, info2 := "info1", "info2"
	info := []string{info1, info2}
	amt := agent.newAspectModifierType(agent.aspect.find(info...), nil)
	assert.Equal(t, amt, agent.memory.find(amt))
	assert.Equal(t, amt.aspect, agent.aspect.find(info...))

	amtCopy := agent.newAspectModifierType(agent.aspect.find(info...), nil)
	assert.Equal(t, amt, amtCopy)
}

func TestAspectModifierTypeToString(t *testing.T) {
	agent := newEmptyWorldAgent()
	info1, info2 := "info1", "info2"
	info := []string{info1, info2}
	amt := agent.newAspectModifierType(agent.aspect.find(info...), nil)
	amtStr := amt.debug("", 2)
	for _, infoStr := range info {
		assert.Contains(t, amtStr, infoStr)
	}
}

func TestAspectModifierTypeInstantiate(t *testing.T) {
	agent := newEmptyWorldAgent()
	info1, info2 := "info1", "info2"
	info := []string{info1, info2}
	amt := agent.newAspectModifierType(agent.aspect.find(info...), nil)
	tc := agent.newTestConcept(1, nil)
	assert.Empty(t, amt.sources())
	am := amt.instantiate(tc, conceptSourceObservation, nil)
	assert.Len(t, amt.sources(), 1)
	assert.Contains(t, amt.sources(), conceptSourceObservation)
	assert.Equal(t, amt, am._type())
	assert.Equal(t, tc, am.target())
	assert.Equal(t, conceptSourceObservation, am.source())
	assert.Equal(t, am, amt.instantiate(tc, conceptSourceObservation, nil))
	amStr := am.debug("", 2)
	for _, infoStr := range info {
		assert.Contains(t, amStr, infoStr)
	}
	assert.Contains(t, amStr, modifierSourceNames[conceptSourceObservation])

	val := 123
	valuedAm := amt.instantiate(tc, conceptSourceObservation, nil, val)
	assert.NotEqual(t, am, valuedAm)
	amStr = valuedAm.debug("", 2)
	assert.Contains(t, amStr, strconv.Itoa(val))
}

func TestAspectModifierOverride(t *testing.T) {
	agent := newEmptyWorldAgent()
	info11, info12, info21, info22 := "info11", "info12", "info21", "info22"
	infoList11, infoList12, infoList21 := []string{info11, info21}, []string{info11, info22}, []string{info12, info21}
	amt11 := agent.newAspectModifierType(agent.aspect.find(infoList11...), nil)
	amt12 := agent.newAspectModifierType(agent.aspect.find(infoList12...), nil)
	amt21 := agent.newAspectModifierType(agent.aspect.find(infoList21...), nil)
	tc := agent.newTestConcept(1, nil)

	assert.Len(t, tc.modifiers(nil), 0)
	am11 := amt11.instantiate(tc, conceptSourceObservation, nil)
	assert.Len(t, tc.modifiers(nil), 1)
	assert.Equal(t, tc.modifiers(nil)[am11.id()], am11)
	am12 := amt12.instantiate(tc, conceptSourceObservation, nil)
	assert.Len(t, tc.modifiers(nil), 1)
	assert.Equal(t, tc.modifiers(nil)[am12.id()], am12)
	am21 := amt21.instantiate(tc, conceptSourceObservation, nil)
	assert.Len(t, tc.modifiers(nil), 2)
	assert.Equal(t, tc.modifiers(nil)[am12.id()], am12)
	assert.Equal(t, tc.modifiers(nil)[am21.id()], am21)
}

func TestAspectModifierTypeVerify(t *testing.T) {
	agent := newEmptyWorldAgent()
	info11, info12, info21, info22 := "info11", "info12", "info21", "info22"
	infoList11, infoList12, infoList21 := []string{info11, info21}, []string{info11, info22}, []string{info12, info21}
	amt11 := agent.newAspectModifierType(agent.aspect.find(infoList11...), nil)
	amt12 := agent.newAspectModifierType(agent.aspect.find(infoList12...), nil)
	amt21 := agent.newAspectModifierType(agent.aspect.find(infoList21...), nil)
	tc := agent.newTestConcept(1, nil)
	amt11.instantiate(tc, conceptSourceObservation, nil)
	assert.Nil(t, amt11.verify(nil))
	amt11.lockMap[partIdModifierTarget] = tc
	amt12.lockMap[partIdModifierTarget] = tc
	amt21.lockMap[partIdModifierTarget] = tc
	assert.True(t, *amt11.verify(nil))
	assert.False(t, *amt12.verify(nil))
	assert.Nil(t, amt21.verify(nil))
}

func TestAspectModifierShareSync(t *testing.T) {
	agent := newEmptyWorldAgent()
	info1, info2 := "info1", "info2"
	info := []string{info1, info2}
	amt := agent.newAspectModifierType(agent.aspect.find(info...), nil)
	tc := agent.newTestConcept(1, nil)
	am := amt.instantiate(tc, conceptSourceObservation, nil)
	_, sync := am.instShareParts()
	assert.Len(t, sync, 1)
	assert.Equal(t, partIdModifierTarget, sync[tc.id()])
}

func TestAspectModifierGeneralize(t *testing.T) {
	agent := newEmptyWorldAgent()
	amt1 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "red"), nil)
	amt2 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "blue"), nil)
	amtG := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color"), nil)
	amt1.generalize(amt2)
	assert.Equal(t, amtG, amt1.lowestCommonGeneralization(amt2))
}

func TestAspectModifierGroup(t *testing.T) {
	agent := newEmptyWorldAgent()
	amt1 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "red"), nil)
	amt2 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "blue"), nil)
	so := agent.newSimpleObject(123, nil)
	tpos, _ := generateTime(agent, 0, 2)
	am1 := amt1.instantiate(so, conceptSourceGeneralization, map[int]any{conceptArgTime: tpos[0]})
	am2 := amt2.instantiate(so, conceptSourceGeneralization, map[int]any{conceptArgTime: tpos[1]})
	amG := am1.buildGroup(map[int]concept{am2.id(): am2}).(modifier)
	assert.NotNil(t, amG)
	assert.Equal(t, amG._type(), amt1.lowestCommonGeneralization(amt2))
}

func TestAspectModifierVersionInterrupt(t *testing.T) {
	agent := newEmptyWorldAgent()
	amt1 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "red"), nil)
	amt2 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "blue"), nil)
	so := agent.newSimpleObject(123, nil)
	_, tsos := generateTime(agent, 0, 7)
	amt1.instantiate(so, conceptSourceGeneralization, map[int]any{conceptArgTime: tsos[1][6]})
	assert.Len(t, so.modifiers(nil), 1)
	amt2.instantiate(so, conceptSourceGeneralization, map[int]any{conceptArgTime: tsos[3][4]})
	assert.Len(t, so.modifiers(nil), 4)
}
