package agent

type groupModifier struct {
	*abstractModifier
	members map[int]*memReference
}

func (m *groupModifier) match(other concept) bool {
	//TODO implement me
	panic("implement me")
}

func (a *Agent) newGroupModifier(t modifierType, members map[int]modifier, target concept,
	source int, args map[int]any) *groupModifier {
	result := &groupModifier{members: map[int]*memReference{}}
	a.newAbstractModifier(result, t, target, source, args, &result.abstractModifier)
	for _, member := range members {
		result.members[member.id()] = member.createReference(result, false)
	}
	return result.memorize().(*groupModifier)
}
