package agent

// comparing two numerical intervals to produce a categorical result
// cases:
//
//	[1, 2] less than [2, 3]
//	[1, 2] less than [3, 4]
//	[1, 2] equals [1, 2]
//	[1, 4] contains [2, 3]
type intervalRelation struct {
	*abstractRelation
}

func (r *intervalRelation) match(other concept) bool {
	o, ok := other.(*intervalRelation)
	if !ok {
		return false
	}

	if r.abstractRelation.match(o.abstractRelation) {
		return true
	}

	reverse := map[int]int{
		intervalIdEQ: intervalIdEQ,
		intervalIdGT: intervalIdLT,
		intervalIdLT: intervalIdGT,
	}

	if reverse[r._type().(*intervalRelationType).intervalId] ==
		o._type().(*intervalRelationType).intervalId {
		return r.abstractRelation.matchReverse(o.abstractRelation)
	}

	return false
}

func (r *intervalRelation) interpret() {
	//TODO implement me
	panic("implement me")
}

func (a *Agent) newIntervalRelation(t *intervalRelationType, lTarget, rTarget object,
	args map[int]any) *intervalRelation {
	result := &intervalRelation{}
	a.newAbstractRelation(result, t, lTarget, rTarget, args, &result.abstractRelation)
	return result.memorize().(*intervalRelation)
}

const (
	intervalIdGT = iota
	intervalIdEQ
	intervalIdLT
	intervalIdCT
)

type intervalRelationType struct {
	*abstractRelationType
	intervalId int
}

func (t *intervalRelationType) match(other concept) bool {
	o, ok := other.(*intervalRelationType)
	return ok && t.abstractRelationType.match(o.abstractRelationType) &&
		t.intervalId == o.intervalId
}

func (t *intervalRelationType) verify(args ...any) *bool {
	//TODO implement me
	panic("implement me")
}

func (a *Agent) newIntervalRelationType(intervalId int, args map[int]any) *intervalRelationType {
	result := &intervalRelationType{intervalId: intervalId}
	a.newAbstractRelationType(result, args, &result.abstractRelationType)
	return result.memorize().(*intervalRelationType)
}
