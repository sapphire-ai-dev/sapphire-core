package agent

// context: a setting where a concept exists, i.e. reality, imagination
// typically any description that starts with "if" or "suppose" creates a context
// i.e.
// "bob ate your apple" -> reality context
// "what would you do if bob ate your apple" -> assumption context
type contextObject struct {
	*abstractObject
	creation *memReference
}

func (o *contextObject) match(other concept) bool {
	n, ok := other.(*contextObject)
	return ok && o.abstractObject.match(n.abstractObject) && o.creation.c.match(n.creation.c)
}

func (a *Agent) newContextObject(creation action) *contextObject {
	result := &contextObject{}
	a.newAbstractObject(result, nil, &result.abstractObject)
	result.creation = creation.createReference(result, true) // todo not sure if this should be true
	return result.memorize().(*contextObject)
}

type contextObjectType struct {
	*abstractObjectType
}

func (t *contextObjectType) match(other concept) bool {
	o, ok := other.(*contextObjectType)
	return ok && t.abstractObjectType._match(o.abstractObjectType)
}

func (a *Agent) newContextObjectType(source int) *contextObjectType {
	result := &contextObjectType{}
	a.newAbstractObjectType(result, source, nil, &result.abstractObjectType)
	return result.memorize().(*contextObjectType)
}
