package agent

type groupRelation struct {
	*abstractRelation
	members map[int]*memReference
}

func (r *groupRelation) interpret() {
	for _, member := range parseRefs[relation](r.agent, r.members) {
		member.interpret()
	}
}

func (r *groupRelation) match(other concept) bool {
	o, ok := other.(*groupRelation)
	if !ok || !r.abstractRelation.match(o.abstractRelation) || len(r.members) != len(o.members) {
		return false
	}

	for _, mm := range r.members {
		if om, seen := o.members[mm.c.id()]; !seen || mm.c != om.c {
			return false
		}
	}

	return true
}

func (a *Agent) newGroupRelation(members map[int]relation, args map[int]any) *groupRelation {
	result := &groupRelation{members: map[int]*memReference{}}
	var lTarget, rTarget concept
	var t relationType
	for _, member := range members {
		if lTarget == nil {
			lTarget = member.lTarget()
			rTarget = member.rTarget()
			t = member._type()
		} else if member.lTarget() != lTarget || member.rTarget() != rTarget {
			return nil
		}

		if t != member._type() {
			t.generalize(member._type())
			t = t.lowestCommonGeneralization(member._type()).(relationType)
		}
	}

	a.newAbstractRelation(result, t, lTarget, rTarget, args, &result.abstractRelation)
	for _, member := range members {
		result.members[member.id()] = member.createReference(result, false)
	}

	result = result.memorize().(*groupRelation)
	lTarget.addRelation(result)
	rTarget.addRelation(result)
	return result
}
