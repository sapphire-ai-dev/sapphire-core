package agent

type createContextAction struct {
    *abstractAction
}

func (a *createContextAction) match(other concept) bool {
    o, ok := other.(*createContextAction)
    return ok && a.abstractAction.match(o.abstractAction)
}

func (a *Agent) newCreateContextAction(t *createContextActionType, performer object) *createContextAction {
    result := &createContextAction{}
    a.newAbstractAction(result, t, performer, nil, &result.abstractAction)
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
