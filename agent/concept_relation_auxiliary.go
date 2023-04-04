package agent

type auxiliaryRelation struct {
	*abstractRelation
	_wantChange *memReference // weird naming... this is the change that carries value for a [want] auxiliaryRelation
}

func (r *auxiliaryRelation) match(other concept) bool {
	o, ok := other.(*auxiliaryRelation)
	return ok && r.abstractRelation.match(o.abstractRelation)
}

func (r *auxiliaryRelation) versionCollides(other concept) bool {
	o, ok := other.(*auxiliaryRelation)
	return ok && r.lTarget() == o.lTarget() && r.rTarget() == o.rTarget() &&
		r._type().(*auxiliaryRelationType).auxiliaryTypeId == o._type().(*auxiliaryRelationType).auxiliaryTypeId
}

// also disjoints self from target to prevent infinite recursion on versioning component
func (r *auxiliaryRelation) versioningReplicate() concept {
	delete(r.lTarget().abs()._relations, r.cid)
	delete(r.rTarget().abs()._relations, r.cid)
	result := &auxiliaryRelation{}

	args := map[int]any{}
	if r.ctx() != nil {
		args[conceptArgContext] = r.ctx()
	}

	r.agent.newAbstractRelation(result, r._type(), r.lTarget(), r.rTarget(), args, &result.abstractRelation)
	if r.wantChange() != nil {
		result._wantChange = r.wantChange().createReference(result, true)
	}

	return result
}

func (r *auxiliaryRelation) lObject() object {
	return parseRef[object](r.agent, r._lTarget)
}

func (r *auxiliaryRelation) rAction() performableAction {
	return parseRef[performableAction](r.agent, r._rTarget)
}

func (r *auxiliaryRelation) wantChange() *actionStateChange {
	return parseRef[*actionStateChange](r.agent, r._wantChange)
}

func (r *auxiliaryRelation) interpret() {
	if r._type().(*auxiliaryRelationType).auxiliaryTypeId == auxiliaryTypeIdWant {
		rAction := r.rAction()
		if rAction.performer() != r.agent.self {
			return
		}

		prevMind, postMind := map[int]concept{r.cid: r}, map[int]concept{}
		if wantChange := r.wantChange(); wantChange != nil {
			postMind[wantChange.cid] = wantChange
		}

		rAction.snapshot(snapshotTimingPrev, prevMind)
		rAction.snapshot(snapshotTimingPost, postMind)
		rAction._type().(performableActionType).update(rAction)
	}
}

func (a *Agent) newAuxiliaryRelation(t *auxiliaryRelationType, lTarget, rTarget concept,
	args map[int]any) *auxiliaryRelation {
	_, lTargetIsObject := lTarget.(object)
	_, rTargetIsAction := rTarget.(performableAction)
	if !lTargetIsObject || !rTargetIsAction {
		return nil
	}

	result := &auxiliaryRelation{}
	a.newAbstractRelation(result, t, lTarget, rTarget, args, &result.abstractRelation)
	if wantChange, seen := conceptArg[*actionStateChange](args, conceptArgRelationAuxiliaryWantChange); seen {
		result._wantChange = wantChange.createReference(result, true)
	}

	result = result.memorize().(*auxiliaryRelation)
	lTarget.addRelation(result)
	rTarget.addRelation(result)
	return result
}

const (
	auxiliaryTypeIdStart = iota
	auxiliaryTypeIdBelieve
	auxiliaryTypeIdWant
)

var auxiliaryTypeIdNames = map[string]int{
	"believe": auxiliaryTypeIdBelieve,
	"want":    auxiliaryTypeIdWant,
}

type auxiliaryRelationType struct {
	*abstractRelationType
	auxiliaryTypeId int
	negative        bool
	lType           *memReference
	rType           *memReference
}

func (t *auxiliaryRelationType) instVerifiesCondition(inst concept) bool {
	r, ok := inst.(*auxiliaryRelation)
	lTarget, lOk := t.lockMap[partIdRelationLTarget]
	if !ok || !lOk || r.lTarget() != lTarget || isNil(r._rTarget) {
		return false
	}

	rTargetPA, rTargetIsPA := r.rTarget().(performableAction)
	rT, tOk := r._type().(*auxiliaryRelationType)
	return tOk && rTargetIsPA && rT.auxiliaryTypeId == t.auxiliaryTypeId &&
		rT.negative == t.negative && rTargetPA._type() == t.rType.c
}

func (t *auxiliaryRelationType) instRejectsCondition(inst concept) bool {
	r, ok := inst.(*auxiliaryRelation)
	lTarget, lOk := t.lockMap[partIdRelationLTarget]
	if !ok || !lOk || r.lTarget() != lTarget || isNil(r._rTarget) {
		return false
	}

	rTargetPA, rTargetIsPA := r.rTarget().(performableAction)
	rT, tOk := r._type().(*auxiliaryRelationType)
	return tOk && rTargetIsPA && rT.auxiliaryTypeId == t.auxiliaryTypeId &&
		rT.negative != t.negative && rTargetPA._type() == t.rType.c
}

func (t *auxiliaryRelationType) match(other concept) bool {
	o, ok := other.(*auxiliaryRelationType)
	return ok && t.abstractRelationType.match(o.abstractRelationType) && t.auxiliaryTypeId == o.auxiliaryTypeId &&
		t.negative == o.negative && matchRefs(t.lType, o.lType) && matchRefs(t.rType, o.rType)
}

func (a *Agent) newAuxiliaryRelationType(auxiliaryTypeId int, negative bool,
	lType objectType, rType performableActionType, args map[int]any) *auxiliaryRelationType {
	result := &auxiliaryRelationType{
		auxiliaryTypeId: auxiliaryTypeId,
		negative:        negative,
	}
	a.newAbstractRelationType(result, args, &result.abstractRelationType)
	if lType != nil {
		result.lType = lType.createReference(result, true)
	}
	if rType != nil {
		result.rType = rType.createReference(result, true)
	}
	return result.memorize().(*auxiliaryRelationType)
}
