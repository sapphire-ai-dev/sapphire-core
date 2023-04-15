package agent

type object interface {
	concept
	addType(t objectType)
	types() map[int]objectType
}

type objectType interface {
	concept
	source() int
	modifTypes(source *int) map[int]modifierType
}

type abstractObject struct {
	*abstractConcept
	_types map[int]*memReference
}

func (o *abstractObject) match(n *abstractObject) bool {
	return o.abstractConcept.match(n.abstractConcept)
}

func (o *abstractObject) part(partId int) concept {
	if partId == partIdObjectT {
		for _, t := range o.types() {
			if _, ok := t.(objectType); ok {
				return t
			}
		}
	}

	return o.abstractConcept.part(partId)
}

func (o *abstractObject) clean(r *memReference) {
	o.abstractConcept.clean(r)
	delete(o._types, r.c.id())
}

func (o *abstractObject) addType(t objectType) {
	if _, seen := o._types[t.id()]; seen {
		return
	}
	o._types[t.id()] = t.createReference(o._self, false)
}

func (o *abstractObject) types() map[int]objectType {
	return parseRefs[objectType](o.agent, o._types)
}

func (o *abstractObject) debugArgs() map[string]any {
	args := o.abstractConcept.debugArgs()
	args["types"] = o._types
	return args
}

func (a *Agent) newAbstractObject(self concept, args map[int]any, out **abstractObject) {
	*out = &abstractObject{
		_types: map[int]*memReference{},
	}

	a.newAbstractConcept(self, args, &(*out).abstractConcept)
}

type abstractObjectType struct {
	*abstractConcept
	_source     int
	_modifTypes map[int]map[int]*memReference
}

func (t *abstractObjectType) _match(o *abstractObjectType) bool {
	// if we got the new objectType through observation, it must _match all observed modifiers of an existing objectType
	//  in order to _match
	if t.matchModifTypes(o, &t._source) == false {
		return false
	}
	return t.abstractConcept.match(o.abstractConcept)
}

func (t *abstractObjectType) matchModifTypes(o *abstractObjectType, source *int) bool {
	tModifTypes, oModifTypes := t.modifTypes(source), o.modifTypes(source)
	for tModifId := range tModifTypes {
		// if we have a modifier type that came from the source, but they do not, do not _match
		if _, seen := oModifTypes[tModifId]; !seen {
			return false
		}
	}

	return true
}

func (t *abstractObjectType) debugArgs() map[string]any {
	args := t.abstractConcept.debugArgs()
	args["modifTypes"] = t._modifTypes
	return args
}

func (t *abstractObjectType) source() int {
	return t._source
}

func (t *abstractObjectType) modifTypes(source *int) map[int]modifierType {
	if source != nil {
		return parseRefs[modifierType](t.agent, t._modifTypes[*source])
	}

	result := map[int]modifierType{}
	for _, mts := range t._modifTypes {
		for _, mt := range parseRefs[modifierType](t.agent, mts) {
			result[mt.id()] = mt
		}
	}

	return result
}

func (a *Agent) newAbstractObjectType(self concept, source int, args map[int]any, out **abstractObjectType) {
	*out = &abstractObjectType{
		_source:     source,
		_modifTypes: map[int]map[int]*memReference{},
	}
	a.newAbstractConcept(self, args, &(*out).abstractConcept)
}

func getObjType[T objectType](types map[int]objectType) map[int]T {
	result := map[int]T{}
	for id, t := range types {
		if match, ok := t.(T); ok {
			result[id] = match
		}
	}

	return result
}
