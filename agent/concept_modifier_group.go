package agent

type groupModifier struct {
	*abstractModifier
	members map[int]*memReference
}

func (m *groupModifier) match(other concept) bool {
	o, ok := other.(*groupModifier)
	if !ok || !m.abstractModifier.match(o.abstractModifier) || len(m.members) != len(o.members) {
		return false
	}

	for _, mm := range m.members {
		if om, seen := o.members[mm.c.id()]; !seen || mm.c != om.c {
			return false
		}
	}

	return true
}

func (a *Agent) newGroupModifier(members map[int]modifier, args map[int]any) *groupModifier {
	result := &groupModifier{members: map[int]*memReference{}}
	var target concept
	var t modifierType
	for _, member := range members {
		if target == nil {
			target = member.target()
			t = member._type()
		} else if member.target() != target {
			return nil
		}

		if t != member._type() {
			t.generalize(member._type())
			t = t.lowestCommonGeneralization(member._type()).(modifierType)
		}
	}

	a.newAbstractModifier(result, t, target, conceptSourceGeneralization, args, &result.abstractModifier)
	for _, member := range members {
		result.members[member.id()] = member.createReference(result, false)
	}

	result = result.memorize().(*groupModifier)
	target.addModifier(result)
	return result
}
