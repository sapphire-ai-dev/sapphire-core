package agent

type createContextAction struct {
	*abstractAction
	contextId int // todo: to be replaced
}

func (a *createContextAction) match(other concept) bool {
	o, ok := other.(*createContextAction)
	return ok && a.abstractAction.match(o.abstractAction) && a.contextId == o.contextId
}

func (a *Agent) newCreateContextAction(t *createContextActionType, performer object, contextId int) *createContextAction {
	result := &createContextAction{contextId: contextId}
	args := map[int]any{}
	args[partIdActionT] = t
	args[partIdActionPerformer] = performer
	a.newAbstractAction(result, args, &result.abstractAction)
	return result.memorize().(*createContextAction)
}

type createContextActionType struct {
	*abstractActionType
}

func (t *createContextActionType) match(other concept) bool {
	o, ok := other.(*createContextActionType)
	return ok && t.abstractActionType.match(o.abstractActionType)
}

func (a *Agent) newCreateContextActionType() *createContextActionType {
	result := &createContextActionType{}
	a.newAbstractActionType(result, nil, &result.abstractActionType)
	return result.memorize().(*createContextActionType)
}
