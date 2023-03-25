package agent

type imaginaryObject struct {
	*abstractObject
}

func (o *imaginaryObject) match(_ concept) bool {
	return false
}
