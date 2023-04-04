package agent

type conceptCpntVersioning interface {
	versionCollides(other concept) bool
	prevVersion() concept
	setPrevVersion(prev concept)
	nextVersion() concept
	setNextVersion(next concept)
	updateVersions()
	collectVersions() map[int]concept
	updateSelfVersion(newest concept)
	versioningReplicate() concept // replicate self without memorizing
	versioningReplicaFinalize()   // after replicating, complete memorizing or interpreting here
}

type conceptImplVersioning struct {
	abs          *abstractConcept
	_nextVersion *memReference
	_prevVersion *memReference
}

func (v *conceptImplVersioning) versionCollides(_ concept) bool {
	return false
}

func (v *conceptImplVersioning) prevVersion() concept {
	return parseRef[concept](v.abs.agent, v._prevVersion)
}

func (v *conceptImplVersioning) setPrevVersion(prev concept) {
	if !isNil(prev) {
		v._prevVersion = prev.createReference(v.abs._self, false)
	}
}

func (v *conceptImplVersioning) nextVersion() concept {
	return parseRef[concept](v.abs.agent, v._nextVersion)
}

func (v *conceptImplVersioning) setNextVersion(next concept) {
	if !isNil(next) {
		v._nextVersion = next.createReference(v.abs._self, false)
	}
}

func (v *conceptImplVersioning) updateVersions() {
	// go through abs and self instead of directly calling v.collectVersions to avoid missing override implementations
	for _, version := range v.abs._self.collectVersions() {
		version.updateSelfVersion(v.abs._self)
	}
}

func (v *conceptImplVersioning) collectVersions() map[int]concept {
	return map[int]concept{}
}

func (v *conceptImplVersioning) updateSelfVersion(newest concept) {
	v.updateSelfTime(newest)
}

func (v *conceptImplVersioning) updateSelfTime(newest concept) {
	selfSO, selfEO, selfS, selfE := v.getSelfStartEndHelper()
	var newS, newE *int
	var newSO, newEO *timePointObject

	if !isNil(newest.time()) {
		if newest.time().start() != nil {
			newS = newest.time().start().clockTime
		}
		if newest.time().end() != nil {
			newE = newest.time().end().clockTime
		}
		newSO = newest.time().start()
		newEO = newest.time().end()
	}

	startOverlap := newS == nil || (selfS != nil && *selfS > *newS)
	endOverlap := newE == nil || (selfE != nil && *selfE < *newE)
	if startOverlap && endOverlap { // full overlap
		v.abs._self.replace(newest)
		return
	}

	if !startOverlap && !endOverlap { // full contain
		lRep := v.abs.self().versioningReplicate()
		rRep := v.abs.self().versioningReplicate()
		lRep.setTime(v.abs.agent.time.temporalObjJoin(selfSO, newSO, true))
		rRep.setTime(v.abs.agent.time.temporalObjJoin(newEO, selfEO, true))
		lRep.versioningReplicaFinalize()
		rRep.versioningReplicaFinalize()
		newGroup := lRep.buildGroup(map[int]concept{rRep.id(): rRep})
		newGroup.setTime(v.abs.time())
		v.abs._self.replace(newGroup)
		lRep.setNextVersion(newest)
		newest.setPrevVersion(lRep)
		newest.setNextVersion(rRep)
		rRep.setPrevVersion(newest)
		return
	}

	if startOverlap && (selfS == nil || *selfS < *newE) {
		v.abs.setTime(v.abs.agent.time.temporalObjJoin(newEO, selfEO, true))
		v.abs.setPrevVersion(newest)
		return
	}

	if endOverlap && (selfE == nil || *selfE > *newS) {
		v.abs.setTime(v.abs.agent.time.temporalObjJoin(selfSO, newSO, true))
		v.abs.setNextVersion(newest)
		return
	}
}

func (v *conceptImplVersioning) getSelfStartEndHelper() (*timePointObject, *timePointObject, *int, *int) {
	var selfS, selfE *int
	var selfSO, selfEO *timePointObject
	if !isNil(v.abs.time()) {
		if v.abs.time().start() != nil {
			selfS = v.abs.time().start().clockTime
			selfSO = v.abs.time().start()
		}

		if v.abs.time().end() != nil {
			selfE = v.abs.time().end().clockTime
			selfEO = v.abs.time().end()
		}
	}

	if !isNil(v.prevVersion()) && !isNil(v.prevVersion().time()) {
		prevE := v.prevVersion().time().end().clockTime
		if selfS == nil || (prevE != nil && *selfS < *prevE) {
			selfS = prevE
			selfSO = v.prevVersion().time().end()
		}
	}

	if !isNil(v.nextVersion()) && !isNil(v.nextVersion().time()) {
		nextS := v.nextVersion().time().start().clockTime
		if selfE == nil || (nextS != nil && *selfE > *nextS) {
			selfE = nextS
			selfEO = v.nextVersion().time().start()
		}
	}

	return selfSO, selfEO, selfS, selfE
}

func (v *conceptImplVersioning) versioningReplicate() concept {
	return nil
}

func (v *conceptImplVersioning) versioningReplicaFinalize() {}

func (a *Agent) newConceptImplVersioning(abs *abstractConcept) {
	abs.conceptImplVersioning = &conceptImplVersioning{
		abs: abs,
	}
}
