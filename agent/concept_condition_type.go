package agent

type conditionType interface {
	concept
	verify(args map[int]any) *bool
	verifyCollectInsts(args map[int]any) map[int]concept
	instRejectsCondition(inst concept) bool
	instVerifiesCondition(inst concept) bool
}

type abstractConditionType struct {
	*abstractConcept
}

func (t *abstractConditionType) verify(args map[int]any) *bool {
	self := t._self.(conditionType)
	insts := self.verifyCollectInsts(args)
	verified, rejected := false, false
	for _, inst := range insts {
		if self.instVerifiesCondition(inst) {
			verified = true
		}
		if self.instRejectsCondition(inst) {
			rejected = true
		}
	}

	if verified && rejected {
		return nil
	} else if verified {
		return ternary(true)
	} else if rejected {
		return ternary(false)
	}
	return nil
}

func (a *Agent) newAbstractConditionType(self concept, args map[int]any, out **abstractConditionType) {
	*out = &abstractConditionType{}
	a.newAbstractConcept(self, args, &(*out).abstractConcept)
}
