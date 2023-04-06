package agent

type auxiliaryRelation struct {
	*abstractRelation
	_wantChange      *memReference // weird naming... this is the change that carries value for a [want] auxiliaryRelation
	_actionPerformer *memReference
	_actionReceiver  *memReference
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

func (r *auxiliaryRelation) instShareParts() (map[int]concept, map[int]int) {
	parts, sync := r.abstractRelation.instShareParts()
	if r._actionPerformer != nil {
		parts[partIdRelationAuxiliaryPerformer] = r._actionPerformer.c
		sync[r._actionPerformer.c.id()] = partIdRelationAuxiliaryPerformer
	}
	if r._actionReceiver != nil {
		parts[partIdRelationAuxiliaryReceiver] = r._actionReceiver.c
		sync[r._actionReceiver.c.id()] = partIdRelationAuxiliaryReceiver
	}
	return parts, sync
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
	if rTarget.(performableAction).performer() != nil {
		result._actionPerformer = rTarget.(performableAction).performer().createReference(result, true)
		rTarget.(performableAction).performer().addAuxiliary(result)
	}
	if rTarget.(performableAction).receiver() != nil {
		result._actionReceiver = rTarget.(performableAction).receiver().createReference(result, true)
		rTarget.(performableAction).receiver().addAuxiliary(result)
	}
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

func (t *auxiliaryRelationType) verifyCollectInsts(args map[int]any) map[int]concept {
	insts := t.abstractRelationType.verifyCollectInsts(args)
	rPerformer, pSeen := t.lockMap[partIdRelationAuxiliaryPerformer]
	if pSeen {
		for _, r := range rPerformer.auxiliaries(args) {
			insts[r.id()] = r
		}
	}
	rReceiver, rSeen := t.lockMap[partIdRelationAuxiliaryReceiver]
	if rSeen {
		for _, r := range rReceiver.relations(args) {
			insts[r.id()] = r
		}
	}

	return insts
}

func (t *auxiliaryRelationType) instMatch(r *auxiliaryRelation) bool {
	if isNil(r._lTarget) || isNil(r._rTarget) {
		return false
	}

	selfLTarget, lOk := t.lockMap[partIdRelationLTarget]
	selfRPerformer, rpOk := t.lockMap[partIdRelationAuxiliaryPerformer]
	selfRReceiver, rrOk := t.lockMap[partIdRelationAuxiliaryReceiver]
	if (lOk && r.lTarget() != selfLTarget) ||
		(rpOk && r.rTarget().(performableAction).performer() != selfRPerformer) ||
		(rrOk && r.rTarget().(performableAction).receiver() != selfRReceiver) {
		return false
	}

	if !lOk { // if self did not lock onto a left target, the instance must have the same type
		if _, seen := r.lTarget().(object).types()[t.lType.c.id()]; !seen {
			return false
		}
	}

	rTargetPA, rTargetIsPA := r.rTarget().(performableAction)
	if !rTargetIsPA || rTargetPA._type() != t.rType.c {
		return false
	}

	return true
}

func (t *auxiliaryRelationType) instVerifiesCondition(inst concept) bool {
	r, ok := inst.(*auxiliaryRelation)
	if !ok || !t.instMatch(r) {
		return false
	}

	rT, tOk := r._type().(*auxiliaryRelationType)
	return tOk && rT.auxiliaryTypeId == t.auxiliaryTypeId && rT.negative == t.negative
}

func (t *auxiliaryRelationType) instRejectsCondition(inst concept) bool {
	r, ok := inst.(*auxiliaryRelation)
	if !ok || !t.instMatch(r) {
		return false
	}

	rT, tOk := r._type().(*auxiliaryRelationType)
	return tOk && rT.auxiliaryTypeId == t.auxiliaryTypeId && rT.negative != t.negative
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
