package agent

// lTarget is rTarget
type identityRelation struct {
	*abstractRelation
}

func (r *identityRelation) match(other concept) bool {
	o, ok := other.(*identityRelation)
	return ok && r.abstractRelation.match(o.abstractRelation)
}

func (r *identityRelation) interpret() {
	if r.lTarget().isImaginary() && r.rTarget().isImaginary() == false {
		r.lTarget().replace(r.rTarget())
	} else if r.lTarget().isImaginary() == false && r.rTarget().isImaginary() {
		r.rTarget().replace(r.lTarget())
	} else {
		r.lTarget().applyIdentityRelation(r)
		r.rTarget().applyIdentityRelation(r)
	}
}

func (a *Agent) newIdentityRelation(t *identityRelationType, lTarget, rTarget concept,
	args map[int]any) *identityRelation {
	result := &identityRelation{}
	a.newAbstractRelation(result, t, lTarget, rTarget, args, &result.abstractRelation)
	return result.memorize().(*identityRelation)
}

type identityRelationType struct {
	*abstractRelationType
}

func (t *identityRelationType) instRejectsCondition(inst concept) bool {
	//TODO implement me
	panic("implement me")
}

func (t *identityRelationType) instVerifiesCondition(inst concept) bool {
	//TODO implement me
	panic("implement me")
}

func (t *identityRelationType) match(other concept) bool {
	o, ok := other.(*identityRelationType)
	return ok && t.abstractRelationType.match(o.abstractRelationType)
}

func (a *Agent) newIdentityRelationType(args map[int]any) *identityRelationType {
	result := &identityRelationType{}

	a.newAbstractRelationType(result, args, &result.abstractRelationType)
	return result.memorize().(*identityRelationType)
}
