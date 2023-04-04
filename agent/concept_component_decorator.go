package agent

type conceptCpntDecorator interface {
	modifiers(args map[int]any) map[int]modifier
	addModifier(m modifier)
	relations(args map[int]any) map[int]relation
	addRelation(r relation)
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

func (d *conceptImplDecorator) modifiers(args map[int]any) map[int]modifier {
	result := parseRefs[modifier](d.abs.agent, d._modifiers)
	if temporal, seen := conceptArg[temporalObject](args, conceptArgTime); seen {
		result = filterOverlapTemporal[modifier](d.abs.agent.time, result, temporal)
	}

	return result
}

func (d *conceptImplDecorator) addModifier(m modifier) {
	if _, seen := d._modifiers[m.id()]; seen {
		return
	}

	if m.target() != d.abs._self {
		return
	}

	d._modifiers[m.id()] = m.createReference(d.abs._self, false)
}

func (d *conceptImplDecorator) relations(args map[int]any) map[int]relation {
	result := parseRefs[relation](d.abs.agent, d._relations)
	if temporal, seen := conceptArg[temporalObject](args, conceptArgTime); seen {
		result = filterOverlapTemporal[relation](d.abs.agent.time, result, temporal)
	}

	return result
}

func (d *conceptImplDecorator) addRelation(r relation) {
	if _, seen := d._relations[r.id()]; seen {
		return
	}

	if r.lTarget() != d.abs._self && r.rTarget() != d.abs._self {
		return
	}

	d._relations[r.id()] = r.createReference(d.abs._self, false)
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
	if !isNil(ctx) {
		d._ctx = ctx.createReference(d.abs._self, true)
	}
}

func (d *conceptImplDecorator) time() temporalObject {
	return parseRef[temporalObject](d.abs.agent, d._time)
}

func (d *conceptImplDecorator) setTime(time temporalObject) {
	if time != nil {
		d._time = time.createReference(d.abs._self, false)
	}
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
