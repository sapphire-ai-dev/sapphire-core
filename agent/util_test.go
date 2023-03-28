package agent

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
	"github.com/sapphire-ai-dev/sapphire-core/world/empty"
	"github.com/sapphire-ai-dev/sapphire-core/world/text"
	"runtime"
)

var massTest = true
var massSize = int(1e5)

func newEmptyWorldAgent() *Agent {
	empty.Init()
	world.Reset()
	result := NewAgent()
	return result
}

func newTextWorldAgent() *Agent {
	text.Init()
	world.Reset()
	result := NewAgent()
	return result
}

type TestConcept struct {
	val int
	*abstractConcept
	friends map[int]*memReference
}

func (c *TestConcept) match(other concept) bool {
	o, ok := other.(*TestConcept)
	return ok && c.val == o.val && c.abstractConcept.match(o.abstractConcept)
}

func (c *TestConcept) clean(r *memReference) {
	delete(c.friends, r.c.id())
}

func (c *TestConcept) debugArgs() map[string]any {
	args := c.abstractConcept.debugArgs()
	args["val"] = c.val
	args["friends"] = c.friends
	return args
}

func (a *Agent) newTestConcept(val int, args map[int]any) *TestConcept {
	result := &TestConcept{
		val:     val,
		friends: map[int]*memReference{},
	}

	a.newAbstractConcept(result, args, &result.abstractConcept)
	return result.memorize().(*TestConcept)
}

func memUsage() int {
	// https://golang.org/pkg/runtime/#MemStats
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	return int(m.Alloc)
}
