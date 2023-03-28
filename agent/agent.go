package agent

import "github.com/sapphire-ai-dev/sapphire-core/world"

type Agent struct {
	memory     *agentMemory
	mind       *agentMind
	language   *agentLanguage
	aspect     *agentAspect
	perception *agentPerception
	symbolic   *agentSymbolicRecord
	activity   *agentActivity

	self   *selfObject
	record *partRecord
}

func (a *Agent) cycle() {
	a.perception.cycle()
	a.mind.cycle()
	a.activity.cycle()
}

func NewAgent() *Agent {
	worldId, actionInterfaces := world.NewActor()
	result := &Agent{}
	result.newAgentMemory()
	result.newAgentMind()
	result.newAgentLanguage()
	result.newAgentAspect()
	result.newAgentPerception()
	result.newAgentSymbolicRecord()
	result.newAgentActivity(actionInterfaces)

	result.self = result.newSelfObject(worldId, nil)
	result.record = result.newPartRecord()
	return result
}
