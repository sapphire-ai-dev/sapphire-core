package agent

type symbolicObjectType interface {
	objectType
	str() string
	//breakdowns() []symbolicObjectType
	//addBreakdown(b *exprObjectType)
	symbolicLinks() map[int]concept
	addSymbolicLink(c concept)
}

type abstractSymbolicObjectType struct {
	*abstractObjectType
	_str           string
	_symbolicLinks map[int]*memReference
}

func (o abstractSymbolicObjectType) str() string {
	return o._str
}

func (o abstractSymbolicObjectType) symbolicLinks() map[int]concept {
	return parseRefs[concept](o.agent, o._symbolicLinks)
}

func (o abstractSymbolicObjectType) addSymbolicLink(c concept) {
	if isNil(c) {
		return
	}

	if _, seen := o._symbolicLinks[c.id()]; seen {
		return
	}

	o._symbolicLinks[c.id()] = c.createReference(o._self, false)
}

func (a *Agent) newAbstractSymbolicObjectType(self concept, source int, str string,
	args map[int]any, out **abstractSymbolicObjectType) {
	*out = &abstractSymbolicObjectType{
		_str:           str,
		_symbolicLinks: map[int]*memReference{},
	}

	a.newAbstractObjectType(self, source, args, &(*out).abstractObjectType)
}

type symbolObjectType struct {
	*abstractSymbolicObjectType
	symbol string
}

func (o *symbolObjectType) match(other concept) bool {
	n, ok := other.(*symbolObjectType)
	return ok && o.symbol == n.symbol && o.abstractObjectType._match(n.abstractObjectType)
}

func (a *Agent) newSymbolObjectType(source int, symbol string, args map[int]any) *symbolObjectType {
	result := &symbolObjectType{}
	a.newAbstractSymbolicObjectType(result, source, symbol, args, &result.abstractSymbolicObjectType)
	return result.memorize().(*symbolObjectType)
}

type exprObjectType struct {
	*abstractSymbolicObjectType
	children []*memReference
}

func (o *exprObjectType) match(other concept) bool {
	n, ok := other.(*symbolObjectType)
	return ok && o.abstractObjectType._match(n.abstractObjectType)
}

func (a *Agent) newExpressionObject(source int, expr string, children []symbolicObjectType,
	args map[int]any) *exprObjectType {
	result := &exprObjectType{}
	a.newAbstractSymbolicObjectType(result, source, expr, args, &result.abstractSymbolicObjectType)
	for _, child := range children {
		result.children = append(result.children, child.createReference(result, true))
	}
	return result.memorize().(*exprObjectType)
}
