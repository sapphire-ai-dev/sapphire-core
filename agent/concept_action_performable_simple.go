package agent

// just a wrapper on atomicAction, as it would be inappropriate to set receivers on atomicActionType singletons
type simpleAction struct {
	*abstractPerformableAction
	_child *memReference
}

func (a *simpleAction) match(other concept) bool {
	o, ok := other.(*simpleAction)
	return ok && a.abstractPerformableAction.match(o.abstractPerformableAction) &&
		a._child.c == o._child.c
}

func (a *simpleAction) part(partId int) concept {
	if partId == partIdActionSimpleChild {
		return a.child()
	}

	return a.abstractPerformableAction.part(partId)
}

func (a *simpleAction) debugArgs() map[string]any {
	args := a.abstractPerformableAction.debugArgs()
	args["child"] = a._child
	return args
}

func (a *simpleAction) child() performableAction {
	return parseRef[performableAction](a.agent, a._child)
}

func (a *simpleAction) setReceiver(o object) {
	if a._receiver != nil {
		return
	}
	a._receiver = o.createReference(a._self, false)
	a.child().setReceiver(o)
}

func (a *simpleAction) start() bool {
	if a._state != actionStateIdle {
		a.complete()
		return false
	}

	if a.child().start() == false {
		a.complete()
		return false
	}

	a._state = actionStateActive
	a.snapshot(snapshotTimingPrev, nil)
	return true
}

func (a *simpleAction) step() bool {
	if a._state != actionStateActive {
		a.complete()
		return false
	}

	if a.child().step() == false {
		a.complete()
		return false
	}

	a._state = a.child().state()
	if a._state == actionStateDone {
		a.complete()
	}
	return true
}

func (a *Agent) newSimpleAction(t *simpleActionType, performer object, child performableAction,
	args map[int]any) *simpleAction {
	result := &simpleAction{}
	a.newAbstractPerformableAction(result, t, performer, args, &result.abstractPerformableAction)
	result._child = child.createReference(result, true)
	return result.memorize().(*simpleAction)
}

type simpleActionType struct {
	*abstractPerformableActionType
	_child *memReference
}

func (t *simpleActionType) match(other concept) bool {
	o, ok := other.(*simpleActionType)
	return ok && t.abstractPerformableActionType.match(o.abstractPerformableActionType) &&
		t._child.c == o._child.c
}

func (t *simpleActionType) debugArgs() map[string]any {
	args := t.abstractPerformableActionType.debugArgs()
	args["child"] = t._child
	return args
}

func (t *simpleActionType) instantiate(args map[int]any) map[int]performableAction {
	result := map[int]performableAction{}
	for _, child := range t.child().instantiate(args) {
		inst := t.agent.newSimpleAction(t, t.agent.self, child, nil)
		result[inst.cid] = inst
	}

	return result
}

func (t *simpleActionType) child() performableActionType {
	return parseRef[performableActionType](t.agent, t._child)
}

func (a *Agent) newSimpleActionType(receiverType objectType, child *atomicActionType,
	args map[int]any) *simpleActionType {
	result := &simpleActionType{}
	a.newAbstractPerformableActionType(result, receiverType, args, &result.abstractPerformableActionType)
	result._child = child.createReference(result, true)
	return result.memorize().(*simpleActionType)
}
