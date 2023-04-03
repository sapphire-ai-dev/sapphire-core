package agent

type conceptCpntGroup interface {
	buildGroup(others map[int]concept) concept
}

type conceptImplGroup struct {
	abs *abstractConcept
}

func (g *conceptImplGroup) buildGroup(_ map[int]concept) concept {
	return nil
}

func (a *Agent) newConceptImplGroup(abs *abstractConcept) {
	abs.conceptImplGroup = &conceptImplGroup{
		abs: abs,
	}
}
