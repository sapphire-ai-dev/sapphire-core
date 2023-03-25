package agent

import (
	"sort"
)

type agentMind struct {
	agent       *Agent
	capacity    int
	thoughts    map[concept]bool
	newThoughts map[concept]int
}

func (m *agentMind) add(c concept) {
	m.newThoughts[c]++
}

//func (m *agentMind) cycle() {
//	m.generalize()
//	m.propagateThoughts()
//	m.printBreakdown()
//}
//
//func (m *agentMind) generalize() {
//	thoughtGens := map[int]map[int]concept{}
//	for c := range m.thoughts {
//		thoughtGens[c._id()] = c._generalizations()
//	}
//
//	for l := range m.thoughts {
//		for r := range m.thoughts {
//			if thoughtGens[l._id()][r._id()] != nil || thoughtGens[r._id()][l._id()] != nil {
//				continue
//			}
//
//			l._generalize(r)
//		}
//	}
//}
//
//func (m *agentMind) propagateThoughts() {
//	for c := range m.thoughts {
//		if _, ok := c.(object); ok {
//			continue
//		}
//		if _, ok := c.(change); ok {
//			continue
//		}
//		if _, ok := c.(modifier); ok {
//			continue
//		}
//		m.add(c)
//	}
//
//	for _, c := range m.agent.perception.visibleObjects {
//		m.add(c)
//	}
//
//	m.thoughts = m.filteredNewThoughts()
//	m.newThoughts = map[concept]int{}
//}

func (m *agentMind) filteredNewThoughts() map[concept]bool {
	result := map[concept]bool{}

	type pair struct {
		c concept
		n int
	}

	var thoughtList []pair
	for c, n := range m.newThoughts {
		if m.agent.memory.find(c) == c {
			thoughtList = append(thoughtList, pair{c, n})
		}
	}

	sort.SliceStable(thoughtList, func(i, j int) bool {
		return thoughtList[i].n > thoughtList[j].n
	})

	if len(thoughtList) > m.capacity {
		thoughtList = thoughtList[:m.capacity]
	}

	for _, p := range thoughtList {
		result[p.c] = true
	}

	return result
}

// match mind concepts by type
func mindConcepts[T concept](m *agentMind) map[int]T {
	result := map[int]T{}

	for c := range m.thoughts {
		if tc, ok := c.(T); ok {
			result[tc.id()] = tc
		}
	}

	return result
}

func (a *Agent) newAgentMind() {
	result := &agentMind{
		agent:       a,
		thoughts:    map[concept]bool{},
		newThoughts: map[concept]int{},
	}

	a.mind = result
}
