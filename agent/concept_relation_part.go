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

func (a *Agent) newPartRelation(t *partRelationType, lTarget, rTarget concept) *partRelation {
	result := &partRelation{}
	a.newAbstractRelation(result, t, lTarget, rTarget, &result.abstractRelation)
	return result.memorize().(*partRelation)
}

type partRelationType struct {
	*abstractRelationType
	partId int
}

func (t *partRelationType) match(other concept) bool {
	o, ok := other.(*partRelationType)
	return ok && t.partId == o.partId && t.abstractRelationType.match(o.abstractRelationType)
}

func (t *partRelationType) verify(_ ...any) *bool {
	if t.lockMap == nil {
		return nil
	}

	lTarget, lSeen := t.lockMap[partIdRelationLTarget]
	rTarget, rSeen := t.lockMap[partIdRelationRTarget]
	if !lSeen || !rSeen {
		return nil
	}

	lTarget.genPartRelations()
	rTarget.genPartRelations()
	insts, certainFalse := t.abstractRelationType.verifyInsts()
	if certainFalse != nil {
		return certainFalse
	}

	for _, inst := range insts {
		if inst.lTarget() == lTarget && inst.rTarget() == rTarget && inst._type() == t {
			return ternary(true)
		}
	}

	return nil
}

func (a *Agent) newPartRelationType(partId int) *partRelationType {
	result := &partRelationType{
		partId: partId,
	}

	a.newAbstractRelationType(result, &result.abstractRelationType)
	return result
}
