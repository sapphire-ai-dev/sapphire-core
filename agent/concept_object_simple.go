package agent

type simpleObject struct {
	*abstractObject
	worldId int
}

func (o *simpleObject) part(partId int) concept {
	if partId == partIdObjectGroupSize {
		return o.agent.symbolic.numerics.number1
	}

	return o.abstractObject.part(partId)
}

func (o *simpleObject) match(other concept) bool {
	n, ok := other.(*simpleObject)
	return ok && n.abstractObject.match(n.abstractObject) && o.worldId == n.worldId
}

func (o *simpleObject) debugArgs() map[string]any {
	args := o.abstractObject.debugArgs()
	args["worldId"] = o.worldId
	return args
}

func (a *Agent) newSimpleObject(worldId int, args map[int]any) *simpleObject {
	result := &simpleObject{worldId: worldId}
	a.newAbstractObject(result, args, &result.abstractObject)
	return result.memorize().(*simpleObject)
}

func (a *Agent) interpretSimpleObject(concepts map[int]concept, _ ...any) concept {
	result := &simpleObject{worldId: -1}
	a.newAbstractObject(result, nil, &result.abstractObject)
	for partId, part := range concepts {
		if partId == partIdObjectT {
			if t, ok := part.(objectType); ok {
				result.addType(t)
			}
		}
	}
	result.unique = true
	return result.memorize().(*simpleObject)
}

type simpleObjectType struct {
	*abstractObjectType
}

func (t *simpleObjectType) match(other concept) bool {
	o, ok := other.(*simpleObjectType)
	return ok && t.abstractObjectType._match(o.abstractObjectType)
}

func (t *simpleObjectType) debugArgs() map[string]any {
	args := t.abstractObjectType.debugArgs()
	return args
}

func (t *simpleObjectType) generalize(other concept) {
	o, ok := other.(*simpleObjectType)
	if !ok { // not a simple object type, elevate to abstract object type
		o.abstractObjectType.generalize(other)
		return
	}

	if !isNil(t.lowestCommonGeneralization(o)) {
		return
	}

	commonModifs := mapIntersection[modifierType](t.modifTypes(nil), o.modifTypes(nil))
	if len(commonModifs) == 0 { // if there is nothing in common, do not generalize (todo: is this correct?)
		return
	}

	args := map[int]any{}
	if ctx, ctxMatch := t.agent.commonCtx(t, o); ctxMatch {
		args[conceptArgContext] = ctx
	}

	gen := t.agent.newSimpleObjectType(conceptSourceGeneralization, commonModifs, args)
	gen._linkGeneralization(t, o)
}

func (a *Agent) newSimpleObjectType(source int, modifTypes map[int]modifierType, args map[int]any) *simpleObjectType {
	result := &simpleObjectType{}
	a.newAbstractObjectType(result, source, args, &result.abstractObjectType)
	if _, seen := result._modifTypes[source]; !seen {
		result._modifTypes[source] = map[int]*memReference{}
	}
	for modifTypeId, modifType := range modifTypes {
		result._modifTypes[source][modifTypeId] = modifType.createReference(result, false)
	}
	return result.memorize().(*simpleObjectType)
}

func (a *Agent) interpretSimpleObjectType(_ map[int]concept, _ ...any) concept {
	result := &simpleObjectType{}
	a.newAbstractObjectType(result, conceptSourceLanguage, nil, &result.abstractObjectType)
	result.unique = true
	return result.memorize().(*simpleObjectType)
}
