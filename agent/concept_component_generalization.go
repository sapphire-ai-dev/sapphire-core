package agent

// generalizations and specifications form a directed acyclic graph
type conceptCpntGeneralization interface {
	generalize(other concept)
	generalizations() map[int]concept                 // recursive
	specifications() map[int]concept                  // recursive
	generalizationsHelper(result map[int]concept)     // recursive
	specificationsHelper(result map[int]concept)      // recursive
	lowestCommonGeneralization(other concept) concept // todo should this return a collection?
}

type conceptImplGeneralization struct {
	abs              *abstractConcept
	_generalizations map[int]*memReference
	_specifications  map[int]*memReference
}

func (g *conceptImplGeneralization) generalize(_ concept) {}

func (g *conceptImplGeneralization) generalizations() map[int]concept {
	result := map[int]concept{}
	g.generalizationsHelper(result)
	delete(result, g.abs.cid)
	return result
}

func (g *conceptImplGeneralization) generalizationsHelper(result map[int]concept) {
	if _, seen := result[g.abs.id()]; seen {
		return
	}

	result[g.abs.id()] = g.abs._self
	for _, gen := range parseRefs[concept](g.abs.agent, g._generalizations) {
		gen.generalizationsHelper(result)
	}
}

func (g *conceptImplGeneralization) specifications() map[int]concept {
	result := map[int]concept{}
	g.specificationsHelper(result)
	delete(result, g.abs.cid)
	return result
}

func (g *conceptImplGeneralization) specificationsHelper(result map[int]concept) {
	if _, seen := result[g.abs.id()]; seen {
		return
	}

	result[g.abs.id()] = g.abs._self
	for _, spc := range parseRefs[concept](g.abs.agent, g._specifications) {
		spc.specificationsHelper(result)
	}
}

func (g *conceptImplGeneralization) lowestCommonGeneralization(other concept) concept {
	self := g.abs._self
	if self == other {
		return self
	}

	sGens, oGens, commonGens := self.generalizations(), other.generalizations(), map[int]concept{}
	for _, sGen := range sGens {
		if oGen, seen := oGens[sGen.id()]; seen && oGen == sGen {
			commonGens[sGen.id()] = sGen
		}
	}

	for _, cGen := range commonGens {
		hasCommonGenDescendent := false
		for _, cSpc := range cGen.specifications() {
			if _, seen := commonGens[cSpc.id()]; seen {
				hasCommonGenDescendent = true
				break
			}
		}

		if !hasCommonGenDescendent {
			return cGen
		}
	}

	return nil
}

func (g *conceptImplGeneralization) _linkGeneralization(l, r concept) {
	if g.abs._self != l {
		g._specifications[l.id()] = l.createReference(g.abs._self, false)
		l.abs()._generalizations[g.abs.cid] = g.abs.createReference(l, false)
	}

	if g.abs._self != r {
		g._specifications[r.id()] = r.createReference(g.abs._self, false)
		r.abs()._generalizations[g.abs.cid] = g.abs.createReference(r, false)
	}
}

func (a *Agent) newConceptImplGeneralization(abs *abstractConcept) {
	abs.conceptImplGeneralization = &conceptImplGeneralization{
		abs:              abs,
		_generalizations: map[int]*memReference{},
		_specifications:  map[int]*memReference{},
	}
}
