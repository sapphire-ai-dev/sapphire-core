package agent

type auxiliaryRelation struct {
	*abstractRelation
}

func (r *auxiliaryRelation) match(other concept) bool {
	o, ok := other.(*auxiliaryRelation)
	return ok && r.abstractRelation.match(o.abstractRelation)
}

func (r *auxiliaryRelation) interpret() {

}

func (a *Agent) newAuxiliaryRelation(t *auxiliaryRelationType, lTarget, rTarget concept) *auxiliaryRelation {
	result := &auxiliaryRelation{}
	a.newAbstractRelation(result, t, lTarget, rTarget, &result.abstractRelation)
	return result.memorize().(*auxiliaryRelation)
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

func (t *auxiliaryRelationType) link() {

}

func (a *Agent) newAuxiliaryRelationType(auxiliaryTypeId int) *auxiliaryRelationType {
	result := &auxiliaryRelationType{
		auxiliaryTypeId: auxiliaryTypeId,
	}
	a.newAbstractRelationType(result, &result.abstractRelationType)
	result = result.memorize().(*auxiliaryRelationType)
	result.link()
	return result
}
