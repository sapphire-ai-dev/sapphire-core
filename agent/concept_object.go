package agent

type object interface {
    concept
    addType(t objectType)
    types() map[int]objectType
}

type objectType interface {
    concept
    source() int
    modifTypes() map[int]modifierType
}

type abstractObject struct {
    *abstractConcept
    _types map[int]*memReference
}

func (o *abstractObject) match(n *abstractObject) bool {
    return o.abstractConcept.match(n.abstractConcept)
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
    _modifTypes map[int]*memReference
}

func (t *abstractObjectType) _match(o *abstractObjectType) bool {
    // if we got the new objectType through observation, it must _match all observed modifiers of an existing objectType
    //  in order to _match
    if t.matchModifTypes(o, conceptSourceObservation) == false {
        return false
    }
    return t.abstractConcept.match(o.abstractConcept)
}

func (t *abstractObjectType) matchModifTypes(o *abstractObjectType, source int) bool {
    tModifTypes, oModifTypes := t.modifTypes(), o.modifTypes()
    for tModifId, tModif := range tModifTypes {
        if _, seen := tModif.sources()[source]; !seen {
            continue
        }

        // if we have a modifier type that came from the source, but they do not, do not _match
        if oModif, seen := oModifTypes[tModifId]; !seen || tModif != oModif {
            return false
        }
    }

    for oModifId, oModif := range oModifTypes {
        if _, seen := oModif.sources()[source]; !seen {
            continue
        }

        // if they have a modifier type that came from the source, but we do not, do not _match
        if tModif, seen := tModifTypes[oModifId]; !seen || tModif != tModif {
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

func (t *abstractObjectType) modifTypes() map[int]modifierType {
    return parseRefs[modifierType](t.agent, t._modifTypes)
}

func (a *Agent) newAbstractObjectType(self concept, source int, args map[int]any, out **abstractObjectType) {
    *out = &abstractObjectType{
        _source:     source,
        _modifTypes: map[int]*memReference{},
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
