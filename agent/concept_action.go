package agent

type action interface {
	concept
	_type() actionType
	performer() object
}

type actionType interface {
	concept
}

type abstractAction struct {
	*abstractConcept
	t          *memReference
	_performer *memReference
}

func (a *abstractAction) match(o *abstractAction) bool {
	return a.abstractConcept.match(o.abstractConcept) && matchRefs(a.t, o.t) && matchRefs(a._performer, o._performer)
}

func (a *abstractAction) part(partId int) concept {
	if partId == partIdActionT {
		return a._type()
	}
	if partId == partIdActionPerformer {
		return a.performer()
	}
	return a.abstractConcept.part(partId)
}

func (a *abstractAction) _type() actionType {
	return parseRef[actionType](a.agent, a.t)
}

func (a *abstractAction) performer() object {
	return parseRef[object](a.agent, a._performer)
}

func (a *abstractAction) debugArgs() map[string]any {
	args := a.abstractConcept.debugArgs()
	args["type"] = a.t
	args["performer"] = a.performer
	return args
}

func (a *Agent) newAbstractAction(self concept, args map[int]any, out **abstractAction) {
	t, tOk := conceptArg[actionType](args, partIdActionT)
	performer, pOk := conceptArg[object](args, partIdActionPerformer)
	if !tOk {
		panic("action type not found")
	}

	*out = &abstractAction{}
	a.newAbstractConcept(self, args, &(*out).abstractConcept)
	(*out).t = t.createReference(self, true)
	if pOk {
		(*out)._performer = performer.createReference(self, true)
	}
}

type abstractActionType struct {
	*abstractConcept
}

func (t *abstractActionType) match(o *abstractActionType) bool {
	return t.abstractConcept.match(o.abstractConcept)
}

func (t *abstractActionType) debugArgs() map[string]any {
	args := t.abstractConcept.debugArgs()
	return args
}

func (a *Agent) newAbstractActionType(self concept, args map[int]any, out **abstractActionType) {
	*out = &abstractActionType{}
	a.newAbstractConcept(self, args, &(*out).abstractConcept)
}
