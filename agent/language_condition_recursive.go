package agent

type recursiveLangCondition struct {
	partId int // the target of the condition is this part of some parent concept
}

func (c *recursiveLangCondition) match(other langCond) bool {
	o, ok := other.(*recursiveLangCondition)
	return ok && c.partId == o.partId
}

func (c *recursiveLangCondition) satisfied(root, parent concept, _ *sntcCtx) *bool {
	part := parent.part(c.partId)
	if part == nil {
		return nil
	}

	result := part == root
	return &result
}

func (c *recursiveLangCondition) interpret(root, parent concept, truth *bool, _ *sntcCtx) map[int]concept {
	if truth == nil || *truth == false {
		return map[int]concept{}
	}

	prt := root.abs().agent.newPartRelationType(c.partId, nil)
	pr := root.abs().agent.newPartRelation(prt, parent, root, nil)
	return map[int]concept{pr.id(): pr}
}
