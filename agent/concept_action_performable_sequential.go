package agent

type sequentialAction struct {
	*abstractPerformableAction
	_first *memReference
	_next  *memReference
}

func (a *sequentialAction) match(other concept) bool {
	o, ok := other.(*sequentialAction)
	return ok && a.abstractPerformableAction.match(o.abstractPerformableAction) &&
		a._first.c == o._first.c && a._next.c == o._next.c
}

func (a *sequentialAction) debugArgs() map[string]any {
	args := a.abstractPerformableAction.debugArgs()
	args["first"] = a._first
	args["next"] = a._next
	return args
}

func (a *sequentialAction) part(partId int) concept {
	if partId == partIdActionSequentialFirst {
		return a.first()
	}
	if partId == partIdActionSequentialNext {
		return a.next()
	}

	return a.abstractPerformableAction.part(partId)
}

func (a *sequentialAction) setReceiver(o object) {
	if a._receiver != nil {
		return
	}
	a._receiver = o.createReference(a._self, false)
	//a.next().setReceiver(o)
}

func (a *sequentialAction) first() performableAction {
	return parseRef[performableAction](a.agent, a._first)
}

func (a *sequentialAction) next() performableAction {
	return parseRef[performableAction](a.agent, a._next)
}

func (a *sequentialAction) start() bool {
	if a._state != actionStateIdle {
		a.complete()
		return false
	}

	if a.first().start() == false {
		a.complete()
		return false
	}

	a._state = actionStateActive
	return true
}

func (a *sequentialAction) step() bool {
	if a.first().state() == actionStateIdle {
		a.complete()
		return false
	}

	if a.first().state() == actionStateActive {
		success := a.first().step()
		if !success {
			a.complete()
		}
		return success
	}

	if a.next().state() == actionStateIdle {
		if a.next().start() == false {
			a.complete()
			return false
		}
	}

	if a.next().state() == actionStateDone {
		a.complete()
		return false
	}

	if a.next().state() == actionStateActive {
		if !a.next().step() {
			a.complete()
			return false
		}
	}

	if a.next().state() == actionStateDone {
		a.complete()
		a._state = actionStateDone
	}

	return true
}

func (a *Agent) newSequentialAction(t *sequentialActionType, performer object,
	firstChild, nextChild performableAction) *sequentialAction {
	result := &sequentialAction{}
	a.newAbstractPerformableAction(result, t, performer, &result.abstractPerformableAction)
	result._first = firstChild.createReference(result, true)
	result._next = nextChild.createReference(result, true)
	return result.memorize().(*sequentialAction)
}

type sequentialActionType struct {
	*abstractPerformableActionType
	_first *memReference
	_next  *memReference
}

func (t *sequentialActionType) match(other concept) bool {
	o, ok := other.(*sequentialActionType)

	return ok && t.abstractPerformableActionType.match(o.abstractPerformableActionType) &&
		matchRefs(t._first, o._first) && t._next.c == o._next.c
}

func (t *sequentialActionType) debugArgs() map[string]any {
	args := t.abstractPerformableActionType.debugArgs()
	args["first"] = t._first
	args["next"] = t._next
	return args
}

func (t *sequentialActionType) first() performableActionType {
	return parseRef[performableActionType](t.agent, t._first)
}

func (t *sequentialActionType) next() performableActionType {
	return parseRef[performableActionType](t.agent, t._next)
}

func (t *sequentialActionType) instantiate(args ...any) performableAction {
	return t.agent.newSequentialAction(t, t.agent.self, t.first().instantiate(args), t.next().instantiate(args))
}

func (a *Agent) newSequentialActionType(receiverType objectType, firstChild,
	nextChild performableActionType) *sequentialActionType {
	result := &sequentialActionType{}
	a.newAbstractPerformableActionType(result, receiverType, &result.abstractPerformableActionType)
	result._first = firstChild.createReference(result, true)
	result._next = nextChild.createReference(result, true)
	return result.memorize().(*sequentialActionType)
}
