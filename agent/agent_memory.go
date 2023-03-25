package agent

import (
	"math/rand"
	"reflect"
)

type agentMemory struct {
	agent    *Agent
	types    map[reflect.Type]*typeMemory
	typeId   map[int]reflect.Type
	capacity int
	itemBits int
}

const defaultItemBits = 20
const defaultTypeCapacity = 30000

func (m *agentMemory) newType(t reflect.Type, capacity int) {
	if _, seen := m.types[t]; seen {
		return
	}

	m.typeId[len(m.types)] = t
	m.types[t] = m.newTypeMemory(len(m.types), capacity, m.itemBits)
}

func (m *agentMemory) add(c concept) concept {
	cType := reflect.TypeOf(c)
	m.newType(cType, m.capacity)
	return m.types[cType].add(c)
}

func (m *agentMemory) touch(c any) {
	switch c.(type) {
	case int:
		m.types[m.typeId[c.(int)>>m.itemBits]].touch(c.(int))
	case concept:
		cType := reflect.TypeOf(c)
		m.newType(cType, m.capacity)
		m.types[cType].touch(c.(concept).id())
	}
}

func (m *agentMemory) find(c any) concept {
	switch c.(type) {
	case int:
		return m.types[m.typeId[c.(int)>>m.itemBits]].find(c)
	case concept:
		cType := reflect.TypeOf(c)
		m.newType(cType, m.capacity)
		return m.types[cType].find(c)
	}

	return nil
}

func (m *agentMemory) remove(c concept) {
	m.types[reflect.TypeOf(c)].remove(c.id())
}

func (a *Agent) newAgentMemory() {
	a.memory = &agentMemory{
		agent:    a,
		types:    map[reflect.Type]*typeMemory{},
		typeId:   map[int]reflect.Type{},
		capacity: defaultTypeCapacity,
		itemBits: defaultItemBits,
	}
}

func filterConcepts[T any](a *Agent, refs map[int]*memReference) map[int]T {
	result := map[int]T{}
	for _, ref := range refs {
		c := a.memory.find(ref.c)
		if c == nil {
			continue
		}

		if ct, ok := c.(T); ok {
			result[c.id()] = ct
		}
	}

	return result
}

type typeMemory struct {
	typeId   int
	items    map[int]*memoryItem
	capacity int
	itemBits int
	dllHead  *memoryItem
	dllTail  *memoryItem
}

func (m *agentMemory) newTypeMemory(typeId, capacity, itemBits int) *typeMemory {
	dllHead := &memoryItem{}
	dllTail := &memoryItem{}
	dllHead.connect(dllTail)

	return &typeMemory{
		typeId:   typeId,
		items:    map[int]*memoryItem{},
		capacity: capacity,
		itemBits: itemBits,
		dllHead:  dllHead,
		dllTail:  dllTail,
	}
}

func (m *typeMemory) add(c concept) concept {
	result := m.find(c)

	if result == nil {
		result = m.match(c)
	}

	if result == nil {
		m.evict()

		freeId := m.freeId()
		c.abs().cid = freeId
		m.items[freeId] = &memoryItem{c: c}
		first := m.dllHead.next
		m.dllHead.connect(m.items[freeId])
		m.items[freeId].connect(first)
		result = c
	}

	return result
}

func (m *typeMemory) find(c any) concept {
	var result concept

	if cc, isInt := c.(int); isInt {
		if item, seen := m.items[cc]; seen {
			result = item.c
		}

		return result
	}

	if c == nil || reflect.ValueOf(c).IsNil() {
		return nil
	}

	if cc, isConcept := c.(concept); isConcept {
		if item, seen := m.items[cc.id()]; seen {
			result = item.c
		}
	}

	if result != nil {
		m.touch(result.id())
	}

	return result
}

func (m *typeMemory) match(c concept) concept {
	if c.abs().unique {
		return nil
	}

	for _, item := range m.items {
		if item.c.match(c) {
			return item.c
		}
	}

	return nil
}

func (m *typeMemory) remove(id int) {
	item, seen := m.items[id]
	if !seen {
		return
	}

	item.prev.connect(item.next)
	delete(m.items, id)
	item.c.deprecate()
}

var evicted int = 0

func (m *typeMemory) evict() {
	for len(m.items) >= m.capacity {
		last := m.dllTail.prev
		last.prev.connect(m.dllTail)
		evicted++
		delete(m.items, last.c.id())
		last.c.deprecate()
	}
}

func (m *typeMemory) freeId() int {
	itemId := (m.typeId << m.itemBits) + rand.Intn(1<<m.itemBits)
	for m.items[itemId] != nil {
		itemId = (m.typeId << m.itemBits) + rand.Intn(1<<m.itemBits)
	}
	return itemId
}

func (m *typeMemory) touch(id int) {
	item, seen := m.items[id]
	if !seen {
		return
	}

	item.prev.connect(item.next)
	first := m.dllHead.next
	m.dllHead.connect(item)
	item.connect(first)
}

// doubly linked list node for LRU cache implementation
type memoryItem struct {
	prev *memoryItem
	next *memoryItem
	c    concept
}

func (i *memoryItem) connect(o *memoryItem) {
	i.next = o
	o.prev = i
}
