package agent

type relationChange struct {
	*abstractChange
	before *memReference
	after  *memReference
}

func (c *relationChange) match(other concept) bool {
	o, ok := other.(*relationChange)
	return ok && c.abstractChange.match(o.abstractChange) &&
		c.before.c.match(o.before.c) && c.after.c.match(o.after.c)
}

func (a *Agent) newRelationChange(t *relationChangeType, before, after relation,
	args map[int]any) *relationChange {
	result := &relationChange{}
	a.newAbstractChange(result, t, nil, args, &result.abstractChange)
	if before != nil {
		result.before = before.createReference(result, true)
	}
	if after != nil {
		result.after = after.createReference(result, true)
	}
	return result.memorize().(*relationChange)
}

type relationChangeType struct {
	*abstractChangeType
	beforeType *memReference
	afterType  *memReference
}

func (t *relationChangeType) match(other concept) bool {
	o, ok := other.(*relationChangeType)
	return ok && t.abstractChangeType.match(o.abstractChangeType) &&
		t.beforeType.c.match(o.beforeType.c) && t.afterType.c.match(o.afterType.c)
}

func (a *Agent) newRelationChangeType(beforeType, afterType relationType, args map[int]any) *relationChangeType {
	result := &relationChangeType{}
	a.newAbstractChangeType(result, args, &result.abstractChangeType)
	if beforeType != nil {
		result.beforeType = beforeType.createReference(result, true)
	}
	if afterType != nil {
		result.afterType = afterType.createReference(result, true)
	}
	return result.memorize().(*relationChangeType)
}
