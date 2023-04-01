package agent

// comparing two numerical quantities to produce a categorical result (greater than / less than)
type comparativeRelation struct {
	*abstractRelation
}

func (r *comparativeRelation) match(other concept) bool {
	o, ok := other.(*comparativeRelation)
	if !ok {
		return false
	}

	if r.abstractRelation.match(o.abstractRelation) {
		return true
	}

	reverse := map[int]int{
		comparativeIdEQ: comparativeIdEQ,
		comparativeIdGT: comparativeIdLT,
		comparativeIdGE: comparativeIdLE,
		comparativeIdLE: comparativeIdGE,
		comparativeIdLT: comparativeIdGT,
	}

	if reverse[r._type().(*comparativeRelationType).comparativeId] ==
		o._type().(*comparativeRelationType).comparativeId {
		return r.abstractRelation.matchReverse(o.abstractRelation)
	}

	return false
}

func (r *comparativeRelation) interpret() {
	//TODO implement me
	panic("implement me")
}

func (a *Agent) newComparativeRelation(t *comparativeRelationType, lTarget, rTarget object,
	args map[int]any) *comparativeRelation {
	result := &comparativeRelation{}
	a.newAbstractRelation(result, t, lTarget, rTarget, args, &result.abstractRelation)
	return result.memorize().(*comparativeRelation)
}

const (
	comparativeIdGT = iota
	comparativeIdGE
	comparativeIdEQ
	comparativeIdLE
	comparativeIdLT
)

type comparativeRelationType struct {
	*abstractRelationType
	comparativeId int
}

func (t *comparativeRelationType) match(other concept) bool {
	o, ok := other.(*comparativeRelationType)
	return ok && t.abstractRelationType.match(o.abstractRelationType) &&
		t.comparativeId == o.comparativeId
}

func (t *comparativeRelationType) verify(args ...any) *bool {
	//TODO implement me
	panic("implement me")
}

func (a *Agent) newComparativeRelationType(comparativeId int, args map[int]any) *comparativeRelationType {
	result := &comparativeRelationType{comparativeId: comparativeId}
	a.newAbstractRelationType(result, args, &result.abstractRelationType)
	return result.memorize().(*comparativeRelationType)
}
