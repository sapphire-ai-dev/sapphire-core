package agent

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestTypeMemoryConstructor(t *testing.T) {
	agent := newEmptyWorldAgent()
	typeId, capacity, itemBits := 1, 2, 2
	tm := agent.memory.newTypeMemory(typeId, capacity, itemBits)
	assert.Equal(t, tm.typeId, typeId)
	assert.Equal(t, tm.capacity, capacity)
	assert.Equal(t, tm.itemBits, itemBits)
	assert.Empty(t, tm.items)
	assert.NotNil(t, tm.items)
	assert.NotNil(t, tm.dllHead)
	assert.NotNil(t, tm.dllTail)
	assert.Equal(t, tm.dllHead.next, tm.dllTail)
	assert.Equal(t, tm.dllTail.prev, tm.dllHead)
}

func TestTypeMemoryFind(t *testing.T) {
	agent := newEmptyWorldAgent()
	agent.memory.newType(toReflect[*TestConcept](), 2)
	tm := agent.memory.types[toReflect[*TestConcept]()]
	assert.Nil(t, tm.find(0))
	assert.Nil(t, tm.find(nil))
	tc := &TestConcept{}
	agent.newAbstractConcept(tc, nil, &tc.abstractConcept)
	assert.Nil(t, tm.find(tc))
}

func TestTypeMemoryAdd(t *testing.T) {
	agent := newEmptyWorldAgent()
	agent.memory.newType(toReflect[*TestConcept](), 2)
	tm := agent.memory.types[toReflect[*TestConcept]()]

	tc1 := &TestConcept{val: 1}
	agent.newAbstractConcept(tc1, nil, &tc1.abstractConcept)
	id1 := tc1.id()

	assert.Equal(t, tm.add(tc1), tc1)
	assert.Len(t, tm.items, 1)
	assert.NotEqual(t, tc1.id(), id1)

	assignedId1 := tc1.id()
	assert.Equal(t, tm.items[assignedId1].c, tc1)
	assert.Equal(t, tm.dllHead.next, tm.items[assignedId1])
	assert.Equal(t, tm.items[assignedId1].prev, tm.dllHead)
	assert.Equal(t, tm.items[assignedId1].next, tm.dllTail)
	assert.Equal(t, tm.dllTail.prev, tm.items[assignedId1])

	assert.Equal(t, tm.add(tc1), tc1)
	assert.Len(t, tm.items, 1)
	assert.Equal(t, tc1.id(), assignedId1)
	assert.Equal(t, tm.items[assignedId1].c, tc1)
	assert.Equal(t, tm.dllHead.next, tm.items[assignedId1])
	assert.Equal(t, tm.items[assignedId1].prev, tm.dllHead)
	assert.Equal(t, tm.items[assignedId1].next, tm.dllTail)
	assert.Equal(t, tm.dllTail.prev, tm.items[assignedId1])

	tc1Copy := &TestConcept{val: 1}
	agent.newAbstractConcept(tc1Copy, nil, &tc1Copy.abstractConcept)
	assert.Equal(t, tm.add(tc1Copy), tc1)
	assert.Len(t, tm.items, 1)

	id2 := -2
	tc2 := &TestConcept{val: 2}
	agent.newAbstractConcept(tc2, nil, &tc2.abstractConcept)
	tc2.cid = id2
	assert.Equal(t, tm.add(tc2), tc2)
	assert.Len(t, tm.items, 2)

	assignedId2 := tc2.id()
	assert.Equal(t, tm.items[assignedId1].c, tc1)
	assert.Equal(t, tm.items[assignedId2].c, tc2)
	assert.Equal(t, tm.dllHead.next, tm.items[assignedId2])
	assert.Equal(t, tm.items[assignedId2].prev, tm.dllHead)
	assert.Equal(t, tm.items[assignedId2].next, tm.items[assignedId1])
	assert.Equal(t, tm.items[assignedId1].prev, tm.items[assignedId2])
	assert.Equal(t, tm.items[assignedId1].next, tm.dllTail)
	assert.Equal(t, tm.dllTail.prev, tm.items[assignedId1])

	id3 := -3
	tc3 := &TestConcept{val: 3}
	agent.newAbstractConcept(tc3, nil, &tc3.abstractConcept)
	tc3.cid = id3
	assert.Equal(t, tm.add(tc3), tc3)
	assert.Len(t, tm.items, 2)
	assignedId3 := tc3.id()
	assert.NotContains(t, tm.items, assignedId1)
	assert.Equal(t, tm.items[assignedId2].c, tc2)
	assert.Equal(t, tm.items[assignedId3].c, tc3)
	assert.Equal(t, tm.dllHead.next, tm.items[assignedId3])
	assert.Equal(t, tm.items[assignedId3].prev, tm.dllHead)
	assert.Equal(t, tm.items[assignedId3].next, tm.items[assignedId2])
	assert.Equal(t, tm.items[assignedId2].prev, tm.items[assignedId3])
	assert.Equal(t, tm.items[assignedId2].next, tm.dllTail)
	assert.Equal(t, tm.dllTail.prev, tm.items[assignedId2])
}

func TestTypeMemoryAddUnique(t *testing.T) {
	agent := newEmptyWorldAgent()
	agent.memory.newType(toReflect[*TestConcept](), 2)
	tm := agent.memory.types[toReflect[*TestConcept]()]

	tc1 := &TestConcept{val: 1}
	agent.newAbstractConcept(tc1, nil, &tc1.abstractConcept)
	tm.add(tc1)
	assert.Len(t, tm.items, 1)

	tc1Copy := &TestConcept{val: 1}
	agent.newAbstractConcept(tc1Copy, nil, &tc1Copy.abstractConcept)
	assert.Equal(t, tm.add(tc1Copy), tc1)
	assert.Len(t, tm.items, 1)

	tc1Copy.unique = true
	assert.NotEqual(t, tm.add(tc1Copy), tc1)
	assert.Len(t, tm.items, 2)
}

func TestTypeMemoryRemove(t *testing.T) {
	agent := NewAgent()
	tc := agent.newTestConcept(1, nil)
	assert.Equal(t, agent.memory.types[reflect.TypeOf(tc)].find(tc.id()), tc)

	agent.memory.types[reflect.TypeOf(tc)].remove(-1)
	assert.Equal(t, agent.memory.types[reflect.TypeOf(tc)].find(tc.id()), tc)

	agent.memory.types[reflect.TypeOf(tc)].remove(tc.id())
	assert.Nil(t, agent.memory.types[reflect.TypeOf(tc)].find(tc.id()))
}

func TestTypeMemoryFreeId(t *testing.T) {
	agent := NewAgent()
	agent.memory.capacity = 3
	agent.memory.itemBits = 2
	tc := agent.newTestConcept(1, nil)
	agent.newTestConcept(2, nil)
	agent.newTestConcept(3, nil)
	for i := 0; i < 100; i++ {
		agent.memory.types[reflect.TypeOf(tc)].freeId()
	}
}

func TestTypeMemoryTouch(t *testing.T) {
	agent := newEmptyWorldAgent()
	agent.memory.newType(toReflect[*TestConcept](), 2)
	tm := agent.memory.types[toReflect[*TestConcept]()]

	tc1 := &TestConcept{val: 1}
	tc2 := &TestConcept{val: 2}
	tc3 := &TestConcept{val: 3}
	agent.newAbstractConcept(tc1, nil, &tc1.abstractConcept)
	agent.newAbstractConcept(tc2, nil, &tc2.abstractConcept)
	agent.newAbstractConcept(tc3, nil, &tc3.abstractConcept)
	tc1.cid = 1
	tc2.cid = 2
	tc3.cid = 3

	// nothing should happen
	tm.touch(tc1.id())

	tm.add(tc1)
	tm.add(tc2)
	assert.Contains(t, tm.items, tc1.id())
	assert.Contains(t, tm.items, tc2.id())

	// tc1 has been used earlier than tc2, therefore adding tc3 evicts tc2 instead of tc1
	tm.touch(tc1.id())
	tm.add(tc3)
	assert.Contains(t, tm.items, tc1.id())
	assert.Contains(t, tm.items, tc3.id())
	assert.NotContains(t, tm.items, tc2.id())
}

func TestAgentMemoryConstructor(t *testing.T) {
	agent := newEmptyWorldAgent()
	agent.newAgentMemory()
	am := agent.memory
	assert.Equal(t, am.agent, agent)
	assert.NotNil(t, am.types)
	assert.NotNil(t, am.typeId)
	assert.Empty(t, am.types)
	assert.Empty(t, am.typeId)
}

func TestAgentMemoryNewType(t *testing.T) {
	agent := newEmptyWorldAgent()
	agent.newAgentMemory()
	am := agent.memory
	tc := &TestConcept{val: 1}
	agent.newAbstractConcept(tc, nil, &tc.abstractConcept)
	tc.cid = 1
	capacity := 2

	am.newType(reflect.TypeOf(tc), capacity)
	assert.Len(t, am.types, 1)
	assert.Len(t, am.typeId, 1)

	for i, rt := range am.typeId {
		assert.Equal(t, am.types[rt].typeId, i)
	}
}

func TestAgentMemoryFind(t *testing.T) {
	agent := newEmptyWorldAgent()
	agent.newAgentMemory()
	am := agent.memory
	tc := &TestConcept{val: 1}
	agent.newAbstractConcept(tc, nil, &tc.abstractConcept)
	tc.cid = 1
	assert.Nil(t, am.find(tc))
	assert.Nil(t, am.find(tc.id()))
	am.add(tc)
	assert.Equal(t, am.find(tc), tc)
	assert.Equal(t, am.find(tc.id()), tc)
}

func TestAgentMemoryFindNil(t *testing.T) {
	agent := newEmptyWorldAgent()
	agent.newAgentMemory()
	am := agent.memory
	var nc concept = nil
	assert.Nil(t, am.find(nc))
}

func TestAgentMemoryRemove(t *testing.T) {
	agent := newEmptyWorldAgent()
	agent.newAgentMemory()
	am := agent.memory
	tc := &TestConcept{val: 1}
	agent.newAbstractConcept(tc, nil, &tc.abstractConcept)
	tc.cid = 1
	am.add(tc)
	assert.Equal(t, am.find(tc), tc)
	am.remove(tc)
	assert.Nil(t, am.find(tc))
}

func TestAgentMemoryTouch(t *testing.T) {
	agent := newEmptyWorldAgent()
	agent.newAgentMemory()
	am := agent.memory

	tc1 := &TestConcept{val: 1}
	tc2 := &TestConcept{val: 2}
	tc3 := &TestConcept{val: 3}
	tc4 := &TestConcept{val: 4}
	agent.newAbstractConcept(tc1, nil, &tc1.abstractConcept)
	agent.newAbstractConcept(tc2, nil, &tc2.abstractConcept)
	agent.newAbstractConcept(tc3, nil, &tc3.abstractConcept)
	agent.newAbstractConcept(tc3, nil, &tc4.abstractConcept)
	tc1.cid = 1
	tc2.cid = 2
	tc3.cid = 3
	tc4.cid = 4

	am.capacity = 2
	am.add(tc1)
	am.add(tc2)
	assert.Equal(t, am.find(tc1), tc1)
	assert.Equal(t, am.find(tc2), tc2)

	// tc1 has been used later than tc2, therefore adding tc3 evicts tc2 instead of tc1
	am.touch(tc1)
	am.add(tc3)
	assert.Equal(t, am.find(tc1), tc1)
	assert.Equal(t, am.find(tc3), tc3)
	assert.Nil(t, am.find(tc2))

	am.touch(tc3.id())
	am.add(tc4)
	assert.Nil(t, am.find(tc1))
	assert.Nil(t, am.find(tc2))
	assert.Equal(t, am.find(tc3), tc3)
	assert.Equal(t, am.find(tc4), tc4)
}
