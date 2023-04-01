package agent

type mentionedLangCondition struct{}

func (c *mentionedLangCondition) match(other langCond) bool {
	_, ok := other.(*mentionedLangCondition)
	return ok
}

func (c *mentionedLangCondition) satisfied(root, _ concept, ctx *sntcCtx) *bool {
	if root == nil {
		return ternary(false)
	}

	for _, mentioned := range ctx.convCtx.mentioned {
		if root.match(mentioned) {
			return ternary(true)
		}
	}

	return ternary(false)
}

func (c *mentionedLangCondition) interpret(root, _ concept, truth *bool, ctx *sntcCtx) (concept, map[int]concept) {
	if truth == nil || *truth == false || (root != nil && root.isImaginary() == false) {
		return root, map[int]concept{}
	}

	var matches []concept
	for _, mentioned := range ctx.convCtx.mentioned {
		if root == nil || root.imaginaryFit(mentioned) {
			matches = append(matches, mentioned)
		}
	}

	if len(matches) == 1 {
		if root != nil {
			root.replace(matches[0])
		}
		return matches[0], map[int]concept{}
	}
	// todo raise confusion if there are more than 1 matches
	return root, map[int]concept{}
}
