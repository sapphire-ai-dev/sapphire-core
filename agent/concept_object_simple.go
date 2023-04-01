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

//func (t *simpleObjectType) _generalize(other concept) concept {
//	o, ok := other.(*simpleObjectType)
//	if !ok {
//		return nil
//	}
//
//	// skip common generalizations (TODO MOVE THIS TO CALLER)
//	tGens, oGens := t.generalizations(), o.generalizations()
//	if tGens[o.id] != nil || oGens[t.id] != nil || len(mapIntersection[concept](tGens, oGens)) > 0 {
//		return nil
//	}
//
//	commonModifs := mapIntersection[modifierType](t._modifTypes(), o._modifTypes())
//	if len(commonModifs) == 0 {
//		return nil
//	}
//
//	gen := t.agent.newSimpleObjectType(nil, conceptSourceGeneralization, commonModifs)
//	gen._linkGeneralization(t, o)
//	return gen
//}

func (a *Agent) newSimpleObjectType(source int, modifTypes map[int]modifierType, args map[int]any) *simpleObjectType {
	result := &simpleObjectType{}
	a.newAbstractObjectType(result, source, args, &result.abstractObjectType)
	for modifTypeId, modifType := range modifTypes {
		result._modifTypes[modifTypeId] = modifType.createReference(result, false)
	}
	return result.memorize().(*simpleObjectType)
}

func (a *Agent) interpretSimpleObjectType(_ map[int]concept, _ ...any) concept {
	result := &simpleObjectType{}
	a.newAbstractObjectType(result, conceptSourceLanguage, nil, &result.abstractObjectType)
	result.unique = true
	return result.memorize().(*simpleObjectType)
}
