package agent

type actionStateChange struct {
    *abstractChange
    target *memReference
}

func (c *actionStateChange) match(other concept) bool {
    o, ok := other.(*actionStateChange)
    return ok && c.abstractChange.match(o.abstractChange) && o.target.c.match(c.target.c)
}

func (a *Agent) newActionStateChange(t *actionStateChangeType, target performableAction,
    args map[int]any) *actionStateChange {
    if target.state() != actionStateDone {
        return nil
    }

    result := &actionStateChange{}
    a.newAbstractChange(result, t, nil, nil, &result.abstractChange)
    result.target = target.createReference(result, true)
    return result.memorize().(*actionStateChange)
}

type actionStateChangeType struct {
    *abstractChangeType
    target *memReference
}

func (t *actionStateChangeType) match(other concept) bool {
    o, ok := other.(*actionStateChangeType)
    return ok && t.abstractChangeType.match(o.abstractChangeType) && t.target.c.match(o.target.c)
}

func (a *Agent) newActionStateChangeType(target performableActionType, args map[int]any) *actionStateChangeType {
    result := &actionStateChangeType{}
    a.newAbstractChangeType(result, nil, &result.abstractChangeType)
    result.target = target.createReference(result, true)
    return result.memorize().(*actionStateChangeType)
}
