package agent

import (
	"errors"
	"fmt"
)

type conceptCpntMemory interface {
	id() int
	match(other concept) bool
	override(old concept) (merged concept)
	createReference(user concept, hardDependency bool) *memReference
	deprecate()
	clean(r *memReference) // prevent memory leak from soft dependencies
	replace(replacement concept)
}

type conceptImplMemory struct {
	agent     *Agent
	abs       *abstractConcept
	cid       int
	producers map[*memReference]bool
	consumers map[*memReference]bool
	unique    bool
}

func (m *conceptImplMemory) id() int {
	return m.cid
}

var errSelfReference = errors.New("_self reference")

func (m *conceptImplMemory) memorize() conceptCpntMemory {
	memoryCopy := m.agent.memory.add(m.abs._self)
	if memoryCopy != m.abs._self {
		m.deprecate()
	}

	return memoryCopy
}

func (m *conceptImplMemory) createReference(consumer concept, hardDependency bool) *memReference {
	if m.abs._self == consumer {
		panic(errSelfReference)
	}

	r := &memReference{
		c:        m.abs._self,
		consumer: consumer,
		hard:     hardDependency,
	}

	m.consumers[r] = true
	consumer.abs().producers[r] = true
	return r
}

func (m *conceptImplMemory) deprecate() {
	for ref := range m.producers {
		ref.delete()
	}

	for ref := range m.consumers {
		ref.delete()
	}
}

// to be implemented per class
func (m *conceptImplMemory) clean(_ *memReference) {}

// to be implemented per class
func (m *conceptImplMemory) replace(replacement concept) {
	if replacement == m.abs._self {
		return
	}

	// migrate all [references to _self] to the replacement
	replacementAbs := replacement.abs()
	for ref := range m.consumers {
		ref.c = replacement
		replacementAbs.consumers[ref] = true
		delete(m.consumers, ref)
	}

	m.agent.memory.remove(m.abs._self)

	// remove _self-connections
	rc := replacement.abs()
	for ref := range rc.producers {
		if rc._self == ref.c {
			delete(rc.producers, ref)
		}
	}

	for ref := range rc.consumers {
		if rc._self == ref.consumer {
			delete(rc.consumers, ref)
		}
	}
}

func (c *abstractConcept) match(_ *abstractConcept) bool {
	return true
}

func (c *abstractConcept) override(_ concept) concept {
	return nil
}

func (a *Agent) newConceptImplMemory(abs *abstractConcept) {
	result := &conceptImplMemory{
		agent:     a,
		abs:       abs,
		cid:       -1,
		producers: map[*memReference]bool{},
		consumers: map[*memReference]bool{},
		unique:    false,
	}

	abs.conceptImplMemory = result
}

type memReference struct {
	c        concept
	consumer concept
	hard     bool
}

func (r *memReference) delete() {
	r.consumer.clean(r)
	delete(r.c.abs().consumers, r)
	delete(r.consumer.abs().producers, r)
	if r.hard {
		r.consumer.abs().agent.memory.remove(r.consumer)
	}
	r.c = nil
	r.consumer = nil
}

func (r *memReference) debug(indent string, depth int) string {
	if r == nil {
		return "nil"
	}
	return fmt.Sprintf("%s", r.c.abs()._self.debug(indent, depth))
}

func parseRef[T any](a *Agent, ref *memReference) T {
	var result T
	if ref == nil {
		return result
	}

	c := a.memory.find(ref.c)
	if c == nil {
		return result
	}

	if ct, ok := c.(T); ok {
		return ct
	}

	return result
}

func parseRefs[T any](a *Agent, refs map[int]*memReference) map[int]T {
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

func matchRefs(l, r *memReference) bool {
	if l == nil && r == nil {
		return true
	}
	if l == nil || r == nil {
		return false
	}
	return l.c == r.c
}
