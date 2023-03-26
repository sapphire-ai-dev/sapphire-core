package agent

type change interface {
    action
}

type changeType interface {
    actionType
    value() float64
    addValue(value float64)
}

type abstractChange struct {
    *abstractAction
    params map[string]any
}

func (a *abstractChange) match(o *abstractChange) bool {
    if len(a.params) != len(o.params) {
        return false
    }

    for aParamKey := range a.params {
        if oVal, seen := o.params[aParamKey]; !seen || a.params[aParamKey] != oVal {
            return false
        }
    }
    return a.abstractAction.match(o.abstractAction)
}

func (a *abstractChange) debugArgs() map[string]any {
    args := a.abstractAction.debugArgs()
    for pName, pVal := range a.params {
        args[pName] = pVal
    }
    return args
}

func (a *Agent) newAbstractChange(self concept, t actionType, performer object, args map[int]any,
    out **abstractChange) {
    *out = &abstractChange{
        params: map[string]any{},
    }
    a.newAbstractAction(self, t, performer, args, &(*out).abstractAction)
}

type abstractChangeType struct {
    *abstractActionType
    instances int
    _value    float64
}

func (t *abstractChangeType) match(o *abstractChangeType) bool {
    return t.abstractActionType.match(o.abstractActionType)
}

func (t *abstractChangeType) value() float64 {
    return t._value
}

func (t *abstractChangeType) addValue(value float64) {
    t.instances++
    t._value = (t._value*float64(t.instances-1) + value) / float64(t.instances)
}

func (t *abstractChangeType) debugArgs() map[string]any {
    args := t.abstractActionType.debugArgs()
    return args
}

func (a *Agent) newAbstractChangeType(self concept, args map[int]any, out **abstractChangeType) {
    *out = &abstractChangeType{}
    a.newAbstractActionType(self, args, &(*out).abstractActionType)
}
