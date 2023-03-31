package agent

type participantLangCondition struct {
	participantTypeId int
}

func (c *participantLangCondition) match(other langCond) bool {
	o, ok := other.(*participantLangCondition)
	return ok && c.participantTypeId == o.participantTypeId
}

func (c *participantLangCondition) satisfied(root, _ concept, ctx *sntcCtx) *bool {
	if root == nil {
		return ternary(false)
	}

	var match concept
	if c.participantTypeId == participantTypeIdSpeaker {
		match = ctx.src
	} else if c.participantTypeId == participantTypeIdListener {
		match = ctx.dst
	}

	if match != nil && root.match(match) {
		return ternary(true)
	}

	return nil
}

func (c *participantLangCondition) interpret(root, _ concept, truth *bool, ctx *sntcCtx) (concept, map[int]concept) {
	if truth == nil || *truth == false {
		return root, map[int]concept{}
	}

	var match concept
	if c.participantTypeId == participantTypeIdSpeaker {
		match = ctx.src
	} else if c.participantTypeId == participantTypeIdListener {
		match = ctx.dst
	}

	if match != nil {
		root.replace(match)
	} else {
		panic("should not get here: condition was evaluated as true but match does not exist")
	}
	return match, map[int]concept{}
}

const (
	participantTypeIdSpeaker = iota
	participantTypeIdListener
)
