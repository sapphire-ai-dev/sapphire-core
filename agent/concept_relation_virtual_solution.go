package agent

type virtualSolutionRelation struct {
	*abstractRelation
}

func (r *virtualSolutionRelation) match(other concept) bool {
	o, ok := other.(*virtualSolutionRelation)
	return ok && r.abstractRelation.match(o.abstractRelation)
}

func (r *virtualSolutionRelation) interpret() {
	lType := r.lTarget().(performableAction)._type().(performableActionType)
	rType := r.rTarget().(*virtualAction)._type().(*virtualActionType)
	rType.addSolution(lType)
}

func (a *Agent) newVirtualSolutionRelation(t *virtualSolutionRelationType, lTarget performableAction,
	rTarget *virtualAction, args map[int]any) *virtualSolutionRelation {
	result := &virtualSolutionRelation{}
	a.newAbstractRelation(result, t, lTarget, rTarget, args, &result.abstractRelation)
	return result.memorize().(*virtualSolutionRelation)
}

type virtualSolutionRelationType struct {
	*abstractRelationType
}

func (t *virtualSolutionRelationType) match(other concept) bool {
	o, ok := other.(*virtualSolutionRelationType)
	return ok && t.abstractRelationType.match(o.abstractRelationType)
}

func (t *virtualSolutionRelationType) instRejectsCondition(inst concept) bool {
	//TODO implement me
	panic("implement me")
}

func (t *virtualSolutionRelationType) instVerifiesCondition(inst concept) bool {
	//TODO implement me
	panic("implement me")
}

func (a *Agent) newVirtualSolutionRelationType(args map[int]any) *virtualSolutionRelationType {
	result := &virtualSolutionRelationType{}
	a.newAbstractRelationType(result, args, &result.abstractRelationType)
	return result.memorize().(*virtualSolutionRelationType)
}
