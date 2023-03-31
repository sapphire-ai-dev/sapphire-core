package agent

type imaginaryObject struct {
	*abstractObject
}

func (o *imaginaryObject) match(_ concept) bool {
	return false
}

func (o *imaginaryObject) isImaginary() bool {
	return true
}

func (a *Agent) newImaginaryObject(args map[int]any) concept {
	result := &imaginaryObject{}
	a.newAbstractObject(result, args, &result.abstractObject)
	return result.memorize()
}
