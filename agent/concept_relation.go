package agent

type relation interface {
	concept
	_type() relationType
	lTarget() concept
	rTarget() concept
	interpret()
}

type relationType interface {
	conditionType
}

type abstractRelation struct {
	*abstractConcept
	t        *memReference
	params   map[string]any
	_lTarget *memReference
	_rTarget *memReference
}

func (r *abstractRelation) match(o *abstractRelation) bool {
	if r.t.c != o.t.c || r._lTarget.c != o._lTarget.c || r._rTarget.c != o._rTarget.c ||
		len(r.params) != len(o.params) {
		return false
	}

	for rParamKey := range r.params {
		if oVal, seen := o.params[rParamKey]; !seen || r.params[rParamKey] != oVal {
			return false
		}
	}

	return r.abstractConcept.match(o.abstractConcept)
}

func (r *abstractRelation) matchReverse(o *abstractRelation) bool {
	if r._lTarget.c != o._rTarget.c || r._rTarget.c != o._lTarget.c || len(r.params) != len(o.params) {
		return false
	}

	for rParamKey := range r.params {
		if oVal, seen := o.params[rParamKey]; !seen || r.params[rParamKey] != oVal {
			return false
		}
	}

	return r.abstractConcept.match(o.abstractConcept)
}

func (r *abstractRelation) part(partId int) concept {
	if partId == partIdRelationT {
		return r._type()
	}

	if partId == partIdRelationLTarget {
		return r.lTarget()
	}

	if partId == partIdRelationRTarget {
		return r.rTarget()
	}

	return nil
}

func (r *abstractRelation) _type() relationType {
	return parseRef[relationType](r.agent, r.t)
}

func (r *abstractRelation) lTarget() concept {
	return parseRef[concept](r.agent, r._lTarget)
}

func (r *abstractRelation) rTarget() concept {
	return parseRef[concept](r.agent, r._rTarget)
}

func (r *abstractRelation) shareSync() map[int]int {
	return map[int]int{
		r._lTarget.c.id(): partIdRelationLTarget,
		r._rTarget.c.id(): partIdRelationRTarget,
	}
}

func (r *abstractRelation) debugArgs() map[string]any {
	args := r.abstractConcept.debugArgs()
	args["type"] = r.t
	args["lTarget"] = r._lTarget
	args["rTarget"] = r._rTarget
	for pName, pVal := range r.params {
		args[pName] = pVal
	}
	return args
}

func (a *Agent) newAbstractRelation(self concept, t relationType, lTarget, rTarget concept, args map[int]any,
	out **abstractRelation) {
	*out = &abstractRelation{}
	a.newAbstractConcept(self, args, &(*out).abstractConcept)
	(*out).t = t.createReference(self, true)
	(*out)._lTarget = lTarget.createReference(self, true)
	(*out)._rTarget = rTarget.createReference(self, true)
}

type abstractRelationType struct {
	*abstractConcept
}

func (t *abstractRelationType) match(o *abstractRelationType) bool {
	return t.abstractConcept.match(o.abstractConcept)
}

// scan for relations to potentially match the relationType t.self
//
//	if a relation exist with matching lTarget but different rTarget or matching rTarget but
//	  different lTarget, return ternary false if t.self is exclusive
//	if no such relation found, we are currently unsure, pass all possible candidates to caller
//	  for further checks
func (t *abstractRelationType) verifyInsts() (map[int]relation, *bool) {
	insts := map[int]relation{}
	if t.lockMap == nil {
		return insts, nil
	}

	lTarget, lSeen := t.lockMap[partIdRelationLTarget]
	rTarget, rSeen := t.lockMap[partIdRelationRTarget]
	if !lSeen || !rSeen {
		return insts, nil
	}

	for _, r := range lTarget.relations(nil) {
		if r._type() == t._self {
			insts[r.id()] = r
			// todo: this is not correct for all cases, the logic here is, if there is a relation
			// between L and S of type T, then there isn't a relation between L and R of type T
			// something like if adam's father is bob then adam's father cannot be charlie, this
			// only applies if the relation type is exclusive
			if r.rTarget() != rTarget {
				return map[int]relation{}, ternary(false)
			}
		}
	}

	for _, r := range rTarget.relations(nil) {
		if r._type() == t._self {
			insts[r.id()] = r
			// todo: same as above
			if r.lTarget() != lTarget {
				return map[int]relation{}, ternary(false)
			}
		}
	}

	return insts, nil
}

func (a *Agent) newAbstractRelationType(self concept, args map[int]any, out **abstractRelationType) {
	*out = &abstractRelationType{}
	a.newAbstractConcept(self, args, &(*out).abstractConcept)
}
