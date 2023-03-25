package agent

type concept interface {
	abs() *abstractConcept
	self() concept
	conceptCpntCore
	conceptCpntMemory
	conceptCpntSync
	conceptCpntLang
	conceptCpntDecorator
	conceptCpntImaginary
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

func (a *Agent) newAbstractConcept(self concept, out **abstractConcept) {
	*out = &abstractConcept{
		agent: a,
		_self: self,
	}

	a.newConceptImplCore(*out)
	a.newConceptImplMemory(*out)
	a.newConceptImplSync(*out)
	a.newConceptImplLang(*out)
	a.newConceptImplDecorator(*out)
}

const (
	conceptSourceObservation = iota
	conceptSourceGeneralization
	conceptSourceLanguage
)
