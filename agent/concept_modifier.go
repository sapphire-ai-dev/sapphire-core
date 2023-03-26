package agent

type modifier interface {
    concept
    _type() modifierType
    target() concept
    source() int
}

type modifierType interface {
    conditionType
    sources() map[int]bool
    addSource(source int)
    instantiate(target concept, source int, args ...any) modifier
}

type abstractModifier struct {
    *abstractConcept
    _target *memReference
    t       *memReference
    _source int
}

func (m *abstractModifier) match(o *abstractModifier) bool {
    return m.abstractConcept.match(o.abstractConcept) && m._target.c == o._target.c && m.t.c == o.t.c
}

func (m *abstractModifier) part(partId int) concept {
    if partId == partIdModifierT {
        return m._type()
    }

    if partId == partIdModifierTarget {
        return m.target()
    }

    return nil
}

func (m *abstractModifier) target() concept {
    return parseRef[concept](m.agent, m._target)
}

func (m *abstractModifier) _type() modifierType {
    return parseRef[modifierType](m.agent, m.t)
}

func (m *abstractModifier) instShareParts() map[int]int {
    return map[int]int{
        m._target.c.id(): partIdModifierTarget,
    }
}

func (m *abstractModifier) source() int {
    return m._source
}

var modifierSourceNames = map[int]string{
    conceptSourceObservation:    "[observation]",
    conceptSourceGeneralization: "[generalization]",
    conceptSourceLanguage:       "[language]",
}

func (m *abstractModifier) debugArgs() map[string]any {
    args := m.abstractConcept.debugArgs()
    args["type"] = m.t
    args["source"] = modifierSourceNames[m._source]
    return args
}

func (a *Agent) newAbstractModifier(self modifier, t modifierType,
    target concept, source int, args map[int]any, out **abstractModifier) {
    *out = &abstractModifier{_source: source}
    a.newAbstractConcept(self, args, &(*out).abstractConcept)
    t.addSource(source)
    (*out).t = t.createReference(self, true)
    (*out)._target = target.createReference(self, true)
}

type abstractModifierType struct {
    *abstractConcept
    _sources map[int]bool
}

func (t *abstractModifierType) sources() map[int]bool {
    return t._sources
}

func (t *abstractModifierType) addSource(source int) {
    t._sources[source] = true
}

func (t *abstractModifierType) match(o *abstractModifierType) bool {
    return t.abstractConcept.match(o.abstractConcept)
}

func (a *Agent) newAbstractModifierType(self concept, args map[int]any, out **abstractModifierType) {
    *out = &abstractModifierType{
        _sources: map[int]bool{},
    }
    a.newAbstractConcept(self, args, &(*out).abstractConcept)
}
