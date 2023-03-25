package agent

type imaginaryRelation struct {
	*abstractRelation
}

func (r imaginaryRelation) match(_ concept) bool {
	return false
}

func (r imaginaryRelation) interpret() {}
