package agent

type virtualAction struct {
	*abstractPerformableAction
	_solution *memReference
}

func (a *virtualAction) match(other concept) bool {
	o, ok := other.(*virtualAction)
	return ok && a.abstractPerformableAction.match(o.abstractPerformableAction) && matchRefs(a._solution, o._solution)
}

func (a *virtualAction) part(partId int) concept {
	if partId == partIdActionVirtualSolution {
		return a.solution()
	}

	return a.abstractPerformableAction.part(partId)
}

func (a *virtualAction) debugArgs() map[string]any {
	args := a.abstractPerformableAction.debugArgs()
	args["child"] = a._solution
	return args
}

func (a *virtualAction) solution() performableAction {
	return parseRef[performableAction](a.agent, a._solution)
}

func (a *virtualAction) setReceiver(o object) {
	if a._receiver != nil {
		return
	}
	a._receiver = o.createReference(a._self, false)
	if !isNil(a.solution()) {
		a.solution().setReceiver(o)
	}
}

func (a *virtualAction) start() bool {
	if a._state != actionStateIdle {
		a.complete()
		return false
	}

	if a.solution().start() == false {
		a.complete()
		return false
	}

	a._state = actionStateActive
	a.snapshot(snapshotTimingPrev, nil)
	return true
}

func (a *virtualAction) step() bool {
	if a._state != actionStateActive {
		a.complete()
		return false
	}

	if a.solution().step() == false {
		a.complete()
		return false
	}

	a._state = a.solution().state()
	if a._state == actionStateDone {
		a.complete()
	}
	return true
}

func (a *Agent) newVirtualAction(args map[int]any) *virtualAction {
	result := &virtualAction{}
	a.newAbstractPerformableAction(result, args, &result.abstractPerformableAction)
	if solution, ok := conceptArg[performableAction](args, partIdActionVirtualSolution); ok {
		result._solution = solution.createReference(result, true)
	}
	return result.memorize().(*virtualAction)
}

// example:
// virtualActionType "eat an apple"
// core: eat
// receiverType: "apple"
// solutions: {sequentialActionType ["pick up an apple", "bite"]}
// virtualActionType "eat"
// core: nil
type virtualActionType struct {
	*abstractPerformableActionType
	_core      *memReference
	_solutions map[int]*memReference
}

func (t *virtualActionType) match(other concept) bool {
	o, ok := other.(*virtualActionType)
	return ok && t.abstractPerformableActionType.match(o.abstractPerformableActionType) && matchRefs(t._core, o._core)
}

func (t *virtualActionType) part(partId int) concept {
	if partId == partIdActionVirtualTypeCore {
		return t.core()
	}

	return t.abstractPerformableActionType.part(partId)
}

func (t *virtualActionType) instantiate(args map[int]any) map[int]performableAction {
	result := map[int]performableAction{}
	for _, solution := range t.solutions() {
		for _, solInst := range solution.instantiate(args) {
			args[partIdActionT] = t
			args[partIdActionPerformer] = t.agent.self
			args[partIdActionVirtualSolution] = solInst
			inst := t.agent.newVirtualAction(args)
			result[inst.id()] = inst
		}
	}

	return result
}

func (t *virtualActionType) core() *virtualActionType {
	return parseRef[*virtualActionType](t.agent, t._core)
}

func (t *virtualActionType) addSolution(solution performableActionType) {
	if _, seen := t._solutions[solution.id()]; seen {
		return
	}

	t._solutions[solution.id()] = solution.createReference(t, false)
}

func (t *virtualActionType) solutions() map[int]performableActionType {
	return parseRefs[performableActionType](t.agent, t._solutions)
}

func (a *Agent) newVirtualActionType(args map[int]any) *virtualActionType {
	result := &virtualActionType{
		_solutions: map[int]*memReference{},
	}
	core, _ := conceptArg[*virtualActionType](args, partIdActionVirtualTypeCore)
	receiver, _ := conceptArg[objectType](args, partIdActionReceiver)

	a.newAbstractPerformableActionType(result, receiver, args, &result.abstractPerformableActionType)
	if core != nil {
		result._core = core.createReference(result, true)
	}

	result = result.memorize().(*virtualActionType)
	if core != nil {
		core.addSolution(result)
	}

	return result
}
