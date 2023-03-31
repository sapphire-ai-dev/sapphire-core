package agent

// used to attach parts to a concept C without knowing what C is
type conceptCpntImaginary interface {
	isImaginary() bool
	imaginaryFit(_ concept) bool
}

type conceptImplImaginary struct {
	abs *abstractConcept
}

// imaginary concept classes are responsible for overriding this method to return true
func (i *conceptImplImaginary) isImaginary() bool {
	return false
}

func (i *conceptImplImaginary) imaginaryFit(_ concept) bool {
	return true
}

func (a *Agent) newConceptImplImaginary(abs *abstractConcept) {
	result := &conceptImplImaginary{
		abs: abs,
	}

	abs.conceptImplImaginary = result
}
