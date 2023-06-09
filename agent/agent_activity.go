package agent

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
	"math"
	"math/rand"
)

type agentActivity struct {
	agent                  *Agent
	atomicActionTypes      map[int]*atomicActionType
	activeAction           performableAction
	prevAction             performableAction
	completedActions       []performableAction
	prevActionValues       map[int]float64
	currActionValues       map[int]float64
	atomicActionInterfaces []*world.ActionInterface
}

func (a *agentActivity) cycle() {
	a.reflect()

	if a.activeAction == nil {
		a.startAction()
	}

	if a.activeAction != nil {
		a.propagateAction()
	}
}

var valChangeThreshold = 1.0

func (a *agentActivity) buildSequentialActions() {
	if a.prevAction == nil {
		return
	}

	// temporary workaround prevent sequential actions to build on each other TODO REMOVE
	firstActionType := a.prevAction._type().(performableActionType)
	if _, ok := firstActionType.(*simpleActionType); !ok {
		return
	}

	// note: a.prevAction is one time step earlier than a.prevActionValues
	for prevActionTypeId, prevVal := range a.prevActionValues {
		currVal := a.currActionValues[prevActionTypeId]

		if currVal-prevVal > valChangeThreshold {
			nextActionType := a.agent.memory.find(prevActionTypeId).(performableActionType)
			if _, ok := nextActionType.(*simpleActionType); !ok {
				return
			}
			sat := a.agent.newSequentialActionType(nextActionType.receiverType(), firstActionType, nextActionType, nil)
			a.agent.mind.add(sat)
		}
	}
}

func (a *agentActivity) reflect() {
	for _, ac := range a.completedActions {
		a.reflectSingle(ac)
	}

	a.completedActions = []performableAction{}
}

func (a *agentActivity) reflectSingle(inst performableAction) {
	if inst == nil {
		return
	}

	inst.snapshot(snapshotTimingPost, nil)

	var receiverType objectType
	if inst.receiver() != nil {
		// there should be at most 1 simpleObjectType
		receiverTypeCandidates := getObjType[*simpleObjectType](inst.receiver().types())
		for _, t := range receiverTypeCandidates {
			receiverType = t
		}
	}

	instType := inst._type().(performableActionType)
	prevApat, isAtomic := instType.(*atomicActionType)
	if isAtomic {
		instType.update(inst)
		instType = a.agent.newSimpleActionType(receiverType, prevApat, nil)
		a.agent.mind.add(instType)
	}

	instType.update(inst)
}

func (a *agentActivity) startAction() {
	var bestActions []performableAction
	bestVal := 0.0
	a.clearActionValues()
	for _, pat := range mindConcepts[performableActionType](a.agent.mind) {
		for pa, paVal := range pat.predictValue(map[int]any{partIdConceptTime: a.agent.time.now}) {
			a.currActionValues[pat.id()] = math.Max(a.currActionValues[pat.id()], paVal)

			if paVal > bestVal {
				bestActions = []performableAction{}
				bestVal = paVal
			}

			if paVal == bestVal {
				bestActions = append(bestActions, a.agent.memory.find(pa).(performableAction))
			}
		}
	}

	var bestAction performableAction
	if len(bestActions) > 0 && bestVal > 0 {
		bestAction = bestActions[rand.Intn(len(bestActions))]
	}

	if bestAction != nil {
		a.agent.mind.add(bestAction._type())
		a.activeAction = bestAction
		//if bestActionType.receiverType() != nil {
		//	for _, candidateReceiver := range mindConcepts[object](a.agent.mind) {
		//		if _, seen := candidateReceiver.types()[bestActionType.receiverType().id()]; seen {
		//			a.activeAction.setReceiver(candidateReceiver)
		//			break
		//		}
		//	}
		//}
	}
}

func (a *agentActivity) propagateAction() {
	if a.activeAction.state() == actionStateIdle {
		a.buildSequentialActions()

		success := a.activeAction.start()
		if !success {
			a.clearActiveAction()
			return
		}
	}

	if a.activeAction.state() == actionStateActive {
		success := a.activeAction.step()
		if !success {
			a.clearActiveAction()
			return
		}
	}

	if a.activeAction.state() == actionStateDone {
		a.clearActiveAction()
	}
}

func (a *agentActivity) clearActiveAction() {
	a.prevAction = a.activeAction
	a.activeAction = nil
}

func (a *agentActivity) clearActionValues() {
	a.prevActionValues = a.currActionValues
	a.currActionValues = map[int]float64{}
}

func (a *agentActivity) logTouch(obj object) {
	a.prevAction.setReceiver(obj)
}

func (a *Agent) newAgentActivity(interfaces []*world.ActionInterface) {
	a.activity = &agentActivity{
		agent:             a,
		atomicActionTypes: map[int]*atomicActionType{},
	}

	for _, actionInterface := range interfaces {
		aat := a.newAtomicActionType(actionInterface, nil)
		a.activity.atomicActionTypes[aat.id()] = aat
	}

	a.activity.atomicActionInterfaces = interfaces
}
