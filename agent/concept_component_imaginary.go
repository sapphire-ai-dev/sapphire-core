package agent

import "reflect"

type conceptCpntImaginary interface {
	imagineReflect() reflect.Type
}

type conceptImplImaginary struct {
	abs *abstractConcept
}

func (i *conceptImplImaginary) imagineReflect() reflect.Type {
	t, seen := i.abs.agent.record.imagineReflects[reflect.TypeOf(i.abs.self())]
	if seen {
		return t
	}

	return nil
}

func (a *Agent) newConceptImplImaginary(abs *abstractConcept) {
	result := &conceptImplImaginary{
		abs: abs,
	}

	abs.conceptImplImaginary = result
}
