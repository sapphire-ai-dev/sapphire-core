package agent

type auxiliaryRelation struct {
	*abstractRelation
	_wantChange *memReference // weird naming... this is the change that carries value for a [want] auxiliaryRelation
}

func (r *auxiliaryRelation) match(other concept) bool {
	o, ok := other.(*auxiliaryRelation)
	return ok && r.abstractRelation.match(o.abstractRelation)
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
	result.interpret()
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
}

func (t *auxiliaryRelationType) match(other concept) bool {
	o, ok := other.(*auxiliaryRelationType)
	return ok && t.abstractRelationType.match(o.abstractRelationType) && t.auxiliaryTypeId == o.auxiliaryTypeId
}

func (t *auxiliaryRelationType) verify(_ ...any) *bool {
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

func (a *Agent) newAuxiliaryRelationType(auxiliaryTypeId int, args map[int]any) *auxiliaryRelationType {
	result := &auxiliaryRelationType{
		auxiliaryTypeId: auxiliaryTypeId,
	}
	a.newAbstractRelationType(result, args, &result.abstractRelationType)
	return result.memorize().(*auxiliaryRelationType)
}
