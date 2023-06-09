package agent

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
)

type atomicAction struct {
	*abstractPerformableAction
}

func (a *atomicAction) match(other concept) bool {
	o, ok := other.(*atomicAction)
	return ok && a.abstractPerformableAction.match(o.abstractPerformableAction)
}

func (a *atomicAction) debugArgs() map[string]any {
	args := a.abstractPerformableAction.debugArgs()
	return args
}

func (a *atomicAction) start() bool {
	if a._state != actionStateIdle {
		a.complete()
		return false
	}

	if a._type().(*atomicActionType).actionInterface.Ready() == false {
		a.complete()
		return false
	}

	a._state = actionStateActive
	a.snapshot(snapshotTimingPrev, nil)
	return true
}

func (a *atomicAction) step() bool {
	defer a.complete()
	if a._state != actionStateActive {
		return false
	}

	if a.t.c.(*atomicActionType).actionInterface.Ready() == false {
		return false
	}

	a.t.c.(*atomicActionType).actionInterface.Step()
	a._state = actionStateDone
	return true
}

func (a *Agent) newAtomicAction(args map[int]any) *atomicAction {
	result := &atomicAction{}
	a.newAbstractPerformableAction(result, args, &result.abstractPerformableAction)
	return result.memorize().(*atomicAction)
}

type atomicActionType struct {
	*abstractPerformableActionType
	actionInterface *world.ActionInterface
}

func (t *atomicActionType) match(other concept) bool {
	o, ok := other.(*atomicActionType)
	return ok && t.abstractPerformableActionType.match(o.abstractPerformableActionType) &&
		t.actionInterface == o.actionInterface
}

func (t *atomicActionType) instantiate(args map[int]any) map[int]performableAction {
	if args == nil {
		args = map[int]any{}
	}
	args[partIdActionT] = t
	args[partIdActionPerformer] = t.agent.self
	inst := t.agent.newAtomicAction(args)
	return map[int]performableAction{inst.cid: inst}
}

func (t *atomicActionType) debugArgs() map[string]any {
	args := t.abstractPerformableActionType.debugArgs()
	args["actionInterface"] = t.actionInterface.Name
	return args
}

func (a *Agent) newAtomicActionType(actionInterface *world.ActionInterface, args map[int]any) *atomicActionType {
	result := &atomicActionType{
		actionInterface: actionInterface,
	}
	a.newAbstractPerformableActionType(result, nil, args, &result.abstractPerformableActionType)
	return result.memorize().(*atomicActionType)
}
