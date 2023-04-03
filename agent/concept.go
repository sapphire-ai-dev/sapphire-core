package agent

import "reflect"

type concept interface {
	abs() *abstractConcept
	self() concept
	conceptCpntCore
	conceptCpntMemory
	conceptCpntSync
	conceptCpntLang
	conceptCpntDecorator
	conceptCpntImaginary
	conceptCpntGroup
	conceptCpntVersioning
}

type abstractConcept struct {
	agent *Agent
	_self concept
	*conceptImplCore
	*conceptImplMemory
	*conceptImplSync
	*conceptImplLang
	*conceptImplDecorator
	*conceptImplImaginary
	*conceptImplGroup
	*conceptImplVersioning
}

func (c *abstractConcept) abs() *abstractConcept {
	return c
}

func (c *abstractConcept) self() concept {
	return c._self
}

func (c *abstractConcept) clean(r *memReference) {
	c.conceptImplDecorator.clean(r)
}

func (c *abstractConcept) part(partId int) concept {
	if partId == partIdConceptContext {
		return c.ctx()
	}

	if partId == partIdConceptTime {
		return c.time()
	}

	return nil
}

func (a *Agent) newAbstractConcept(self concept, args map[int]any, out **abstractConcept) {
	*out = &abstractConcept{
		agent: a,
		_self: self,
	}

	a.newConceptImplCore(*out)
	a.newConceptImplMemory(*out)
	a.newConceptImplSync(*out)
	a.newConceptImplLang(*out)
	a.newConceptImplDecorator(*out)
	a.newConceptImplImaginary(*out)
	a.newConceptImplGroup(*out)
	a.newConceptImplVersioning(*out)

	if ctx, seen := conceptArg[*contextObject](args, conceptArgContext); seen {
		(*out).setCtx(ctx)
	}

	if temporal, seen := conceptArg[temporalObject](args, conceptArgTime); seen {
		(*out).setTime(temporal)
	}
}

const (
	conceptSourceObservation = iota
	conceptSourceGeneralization
	conceptSourceLanguage
)

// if struct S implements an interface, a variable s with type S would produce s != nil therefore must use reflect
func isNil(c any) bool {
	return c == nil || reflect.ValueOf(c).IsNil()
}
