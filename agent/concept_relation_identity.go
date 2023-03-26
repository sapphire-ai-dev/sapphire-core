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
	if r.lTarget().imagineReflect() != nil && r.rTarget().imagineReflect() == nil {
		r.lTarget().replace(r.rTarget())
	} else if r.lTarget().imagineReflect() == nil && r.rTarget().imagineReflect() != nil {
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

func (t identityRelationType) match(other concept) bool {
	o, ok := other.(*identityRelationType)
	return ok && t.abstractRelationType.match(o.abstractRelationType)
}

func (t identityRelationType) verify(_ ...any) *bool {
	if t.lockMap == nil {
		return nil
	}

	lTarget, lSeen := t.lockMap[partIdRelationLTarget]
	rTarget, rSeen := t.lockMap[partIdRelationRTarget]
	if !lSeen || !rSeen {
		return nil
	}

	lTarget.genIdentityRelations()
	rTarget.genIdentityRelations()
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

func (a *Agent) newIdentityRelationType(args map[int]any) *identityRelationType {
	result := &identityRelationType{}

	a.newAbstractRelationType(result, args, &result.abstractRelationType)
	return result.memorize().(*identityRelationType)
}
