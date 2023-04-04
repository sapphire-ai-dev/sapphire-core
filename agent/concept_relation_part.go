package agent

// lTarget's PartID is rTarget
type partRelation struct {
	*abstractRelation
}

func (r *partRelation) match(other concept) bool {
	o, ok := other.(*partRelation)
	return ok && r.abstractRelation.match(o.abstractRelation)
}

func (r *partRelation) interpret() {
	r.lTarget().setPart(r.t.c.(*partRelationType).partId, r.rTarget())
}

func (a *Agent) newPartRelation(t *partRelationType, lTarget, rTarget concept, args map[int]any) *partRelation {
	result := &partRelation{}
	a.newAbstractRelation(result, t, lTarget, rTarget, args, &result.abstractRelation)
	return result.memorize().(*partRelation)
}

type partRelationType struct {
	*abstractRelationType
	partId int
}

func (t *partRelationType) instRejectsCondition(inst concept) bool {
	//TODO implement me
	panic("implement me")
}

func (t *partRelationType) instVerifiesCondition(inst concept) bool {
	//TODO implement me
	panic("implement me")
}

func (t *partRelationType) match(other concept) bool {
	o, ok := other.(*partRelationType)
	return ok && t.partId == o.partId && t.abstractRelationType.match(o.abstractRelationType)
}

func (a *Agent) newPartRelationType(partId int, args map[int]any) *partRelationType {
	result := &partRelationType{
		partId: partId,
	}

	a.newAbstractRelationType(result, args, &result.abstractRelationType)
	return result
}
