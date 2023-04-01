package agent

type conceptCpntDecorator interface {
	modifiers() map[int]modifier
	addModifier(m modifier)
	relations() map[int]relation
	genPartRelations() map[int]relation
	applyPartRelation(relation *partRelation)
	genIdentityRelations() map[int]relation
	applyIdentityRelation(r *identityRelation) // support for identityRelation
	ctx() *contextObject
	setCtx(ctx *contextObject)
	time() temporalObject
	setTime(time temporalObject)
}

type conceptImplDecorator struct {
	abs        *abstractConcept
	_modifiers map[int]*memReference
	_relations map[int]*memReference
	_ctx       *memReference
	_time      *memReference
}

func (d *conceptImplDecorator) modifiers() map[int]modifier {
	return parseRefs[modifier](d.abs.agent, d._modifiers)
}

func (d *conceptImplDecorator) addModifier(m modifier) {
	if _, seen := d._modifiers[m.id()]; seen {
		return
	}

	if m.target() != d.abs._self {
		return
	}

	for _, oldModifierRef := range d._modifiers {
		// old modifier is inconsistent with memory (already deleted), continue for now, clean up should happen soon
		if oldModifierRef.c == nil {
			continue
		}

		mergedModifier := m.override(oldModifierRef.c)
		// self merged into existing - exit
		if mergedModifier == oldModifierRef.c {
			return
		}

		// existing merged into self - delete existing and keep merging
		if mergedModifier != nil {
			d.abs.agent.memory.remove(oldModifierRef.c)
			m = mergedModifier.(modifier)
		}
	}

	d._modifiers[m.id()] = m.createReference(d.abs._self, false)
}

func (d *conceptImplDecorator) relations() map[int]relation {
	return parseRefs[relation](d.abs.agent, d._relations)
}

// to be implemented per class
func (d *conceptImplDecorator) genPartRelations() map[int]relation {
	return map[int]relation{}
}

// to be implemented per class
func (d *conceptImplDecorator) applyPartRelation(_ *partRelation) {}

// to be implemented per class
func (d *conceptImplDecorator) genIdentityRelations() map[int]relation {
	return map[int]relation{}
}

// to be implemented per class
func (d *conceptImplDecorator) applyIdentityRelation(_ *identityRelation) {}

func (d *conceptImplDecorator) ctx() *contextObject {
	return parseRef[*contextObject](d.abs.agent, d._ctx)
}

func (d *conceptImplDecorator) setCtx(ctx *contextObject) {
	d._ctx = ctx.createReference(d.abs._self, true)
}

func (d *conceptImplDecorator) time() temporalObject {
	return parseRef[temporalObject](d.abs.agent, d._time)
}

func (d *conceptImplDecorator) setTime(time temporalObject) {
	d._time = time.createReference(d.abs._self, false)
}

func (d *conceptImplDecorator) clean(r *memReference) {
	if _, seen := d._modifiers[r.c.id()]; seen {
		delete(d._modifiers, r.c.id())
	}
	if _, seen := d._relations[r.c.id()]; seen {
		delete(d._relations, r.c.id())
	}
}

// returns args (in case input is nil and output is not nil) and whether there IS a collision
func injectConceptArg(args map[int]any, key int, val concept) (map[int]any, bool) {
	if args == nil {
		args = map[int]any{}
	}
	if ctx, seen := conceptArg[*contextObject](args, key); seen {
		if matchConcepts(ctx, val) == false {
			return args, true
		}
	} else if isNil(val) == false {
		args[conceptArgContext] = val
	}

	return args, false
}

func (a *Agent) newConceptImplDecorator(abs *abstractConcept) {
	abs.conceptImplDecorator = &conceptImplDecorator{
		abs:        abs,
		_modifiers: map[int]*memReference{},
		_relations: map[int]*memReference{},
	}
}
