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

func (r *abstractRelation) collectVersions() map[int]concept {
	result := map[int]concept{}
	for _, c := range r.lTarget().relations(map[int]any{
		conceptArgContext: r.ctx(),
		conceptArgTime:    r.time(),
	}) {
		if r._self.versionCollides(c) {
			result[c.id()] = c
		}
	}

	return result
}

func (r *abstractRelation) versioningReplicaFinalize() {
	r.memorize()
	r.lTarget().addRelation(r._self.(relation))
	r.rTarget().addRelation(r._self.(relation))
}

func (r *abstractRelation) buildGroup(others map[int]concept) concept {
	members := map[int]relation{r.cid: r._self.(relation)}
	for _, other := range others {
		if otherRelation, ok := other.(relation); !ok {
			return nil
		} else {
			members[otherRelation.id()] = otherRelation
		}
	}

	return r.agent.newGroupRelation(members, nil)
}

func (r *abstractRelation) instShareParts() (map[int]concept, map[int]int) {
	parts, sync := map[int]concept{}, map[int]int{}
	parts[partIdRelationLTarget] = r.lTarget()
	sync[r.lTarget().id()] = partIdRelationLTarget
	parts[partIdRelationRTarget] = r.rTarget()
	sync[r.rTarget().id()] = partIdRelationRTarget
	return parts, sync
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
	*abstractConditionType
}

func (t *abstractRelationType) match(o *abstractRelationType) bool {
	return t.abstractConcept.match(o.abstractConcept)
}

// do not check if type or other target matches, as it is possible for relations of different types or other targets
// to verify or reject each other
func (t *abstractRelationType) verifyCollectInsts(args map[int]any) map[int]concept {
	insts := map[int]concept{}
	lTarget, lSeen := t.lockMap[partIdRelationLTarget]
	if lSeen {
		for _, r := range lTarget.relations(args) {
			insts[r.id()] = r
		}
	}
	rTarget, rSeen := t.lockMap[partIdRelationRTarget]
	if rSeen {
		for _, r := range rTarget.relations(args) {
			insts[r.id()] = r
		}
	}

	return insts
}

func (a *Agent) newAbstractRelationType(self concept, args map[int]any, out **abstractRelationType) {
	*out = &abstractRelationType{}
	a.newAbstractConditionType(self, args, &(*out).abstractConditionType)
}
