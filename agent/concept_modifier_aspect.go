package agent

type aspectModifier struct {
	*abstractModifier
	params map[string]any
}

func (m *aspectModifier) match(other concept) bool {
	o, ok := other.(*aspectModifier)
	return ok && m.abstractModifier.match(o.abstractModifier) && matchParams(m.params, o.params)
}

func (m *aspectModifier) override(other concept) concept {
	o, ok := other.(*aspectModifier)
	if !ok || m.target() != o.target() {
		return nil
	}

	mType := m._type().(*aspectModifierType)
	oType := o._type().(*aspectModifierType)
	if mType.aspect == oType.aspect || mType.aspect.parent == oType.aspect.parent {
		return m
	}

	return nil
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

func (t *aspectModifierType) instantiate(target concept, source int, args ...any) modifier {
	params := map[string]any{}
	if len(args) > 0 {
		params["totalValue"] = args[0]
	}

	return t.agent.newAspectModifier(t, target, source, nil, params).memorize().(*aspectModifier)
}

func (t *aspectModifierType) verify(_ ...any) *bool {
	if target, seen := t.lockMap[partIdModifierTarget]; seen {
		for _, m := range target.modifiers() {
			if m._type() == t._self {
				return ternary(true)
			}

			if o, ok := m._type().(*aspectModifierType); ok {
				if o.aspect.parent == t.aspect.parent {
					return ternary(false)
				}
			}
		}
	}

	return nil
}

func (a *Agent) newAspectModifierType(aspect *aspectNode, args map[int]any) *aspectModifierType {
	result := &aspectModifierType{
		aspect: aspect,
	}

	a.newAbstractModifierType(result, args, &result.abstractModifierType)
	return result.memorize().(*aspectModifierType)
}
