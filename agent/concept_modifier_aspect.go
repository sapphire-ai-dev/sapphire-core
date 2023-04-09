package agent

type aspectModifier struct {
	*abstractModifier
	params map[string]any
}

func (m *aspectModifier) match(other concept) bool {
	o, ok := other.(*aspectModifier)
	return ok && m.abstractModifier.match(o.abstractModifier) && matchParams(m.params, o.params)
}

func (m *aspectModifier) versionCollides(other concept) bool {
	o, ok := other.(*aspectModifier)
	if !ok || m.target() != o.target() {
		return false
	}

	mType := m._type().(*aspectModifierType)
	oType := o._type().(*aspectModifierType)
	return mType.aspect == oType.aspect || mType.aspect.parent == oType.aspect.parent
}

// also disjoints self from target to prevent infinite recursion on versioning component
func (m *aspectModifier) versioningReplicate() concept {
	delete(m.target().abs()._modifiers, m.cid)
	result := &aspectModifier{params: m.params}
	args := map[int]any{}
	if m.ctx() != nil {
		args[partIdConceptContext] = m.ctx()
	}

	m.agent.newAbstractModifier(result, m._type(), m.target(), m.source(), args, &result.abstractModifier)
	return result
}

func (m *aspectModifier) debugArgs() map[string]any {
	args := m.abstractModifier.debugArgs()
	for paramName, param := range m.params {
		args[paramName] = param
	}
	return args
}

func (a *Agent) newAspectModifier(t *aspectModifierType, target concept, source int,
	args map[int]any, params map[string]any) *aspectModifier {
	result := &aspectModifier{
		params: params,
	}

	a.newAbstractModifier(result, t, target, source, args, &result.abstractModifier)
	result = result.memorize().(*aspectModifier)
	target.addModifier(result)
	return result
}

type aspectModifierType struct {
	*abstractModifierType
	aspect *aspectNode
}

func (t *aspectModifierType) match(other concept) bool {
	o, ok := other.(*aspectModifierType)
	return ok && t.abstractModifierType.match(o.abstractModifierType) && t.aspect == o.aspect
}

func (t *aspectModifierType) debugArgs() map[string]any {
	args := t.abstractModifierType.debugArgs()
	args["aspect"] = t.aspect.toString()
	return args
}

func (t *aspectModifierType) generalize(other concept) {
	o, ok := generalizeHeader[*aspectModifierType](t, other, t.abstractModifierType)
	if !ok {
		return
	}

	gAsp := t.agent.aspect.lowestCommonAncestor(t.aspect, o.aspect)
	args := map[int]any{}
	if ctx, ctxMatch := t.agent.commonCtx(t, o); ctxMatch {
		args[partIdConceptContext] = ctx
	}

	gen := t.agent.newAspectModifierType(gAsp, args)
	gen._linkGeneralization(t, o)
}

func (t *aspectModifierType) instRejectsCondition(inst concept) bool {
	m, ok := inst.(*aspectModifier)
	return ok && m._type() != t._self && m._type().(*aspectModifierType).aspect.parent == t.aspect.parent
}

func (t *aspectModifierType) instVerifiesCondition(inst concept) bool {
	m, ok := inst.(*aspectModifier)
	return ok && m._type() == t._self
}

func (t *aspectModifierType) instantiate(target concept, source int, args map[int]any, modifArgs ...any) modifier {
	params := map[string]any{}
	if len(modifArgs) > 0 {
		params["totalValue"] = modifArgs[0]
	}

	return t.agent.newAspectModifier(t, target, source, args, params).memorize().(*aspectModifier)
}

func (a *Agent) newAspectModifierType(aspect *aspectNode, args map[int]any) *aspectModifierType {
	result := &aspectModifierType{
		aspect: aspect,
	}

	a.newAbstractModifierType(result, args, &result.abstractModifierType)
	return result.memorize().(*aspectModifierType)
}
