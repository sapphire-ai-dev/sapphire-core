package agent

// used between two concept type classes to record information of the form "my X is your Y"
// we use the terms:
//
//	type: the concept type, correlated with many instances, for example "cats", "mice"
//	inst: the concept instance, correlated with a type, for example "Tom", "Jerry"
//
// instShareSync:
// i.e. we observe that when an action of the agent eating an apple was taken, there is also a
//
//	relation where the apple has the same X position as the agent:
//	(action A).instShareParts() => {agent: performer, "an apple": receiver}
//	(relation R).instShareParts() => {agent: left target, "an apple": right target}
//
// typeUpdateSync:
// we observe that the performer of the action A match the left target of the relation R
// we also observe that the receiver of the action A match the right target of the relation R
// we therefore update this information into the synchronization components of the action TYPE
// "eat" as well as the relation TYPE "x-pos":
// eat.sync[x-pos]: {self-performer: match-lTarget, self-receiver: match-rTarget}
//
// lock / unlockSync:
// after two TYPEs are synchronized, there may not be any inherent links between newly created
// INSTs, we make these associations by locking parts on the TYPEs and search for INSTs
// i.e. given (action B).instShareParts() => {agent: performer, "a banana": receiver}
// we lock (relationType "x-pos") with {left target: agent, right target: "a banana"}
// which allows us to potentially find some (relation S) that describes the above
// afterwards, (relationType "x-pos") can simply unlock to reset its state, in order to lock
// onto some other set of parts in future
type conceptCpntSync interface {
	instShareParts() (map[int]concept, map[int]int) // partId -> concept, conceptId -> partId
	typeUpdateSync(matchType concept, selfInstParts, matchInstParts map[int]int)
	typeLockSync(sourceType concept, parts map[int]concept) // partId -> concept
	typeUnlockSync()
}

type conceptImplSync struct {
	abs *abstractConcept

	// matching type id -> sync entry
	syncMap map[int]*syncEntry
	lockMap map[int]concept
}

func (s *conceptImplSync) instShareParts() (map[int]concept, map[int]int) {
	return map[int]concept{}, map[int]int{}
}

func (s *conceptImplSync) typeUpdateSync(match concept, selfInstParts, matchInstParts map[int]int) {
	affectedEntries := map[*syncEntry]bool{}
	for partConceptId, selfInstPartId := range selfInstParts {
		if matchInstPartId, seen := matchInstParts[partConceptId]; seen {
			for entry := range s.typeUpdateSyncHelper(match, selfInstPartId, matchInstPartId) {
				affectedEntries[entry] = true
			}
		}
	}

	for entry := range affectedEntries {
		entry.count++
	}
}

func (s *conceptImplSync) typeUpdateSyncHelper(match concept, selfInstPartId, matchInstPartId int) map[*syncEntry]bool {
	affectedEntries := map[*syncEntry]bool{}
	sId, mId := s.abs.id(), match.id()
	if _, seen := s.syncMap[mId]; !seen {
		s.syncMap[mId] = s.newSyncEntry()
	}
	affectedEntries[s.syncMap[mId].add(selfInstPartId, matchInstPartId)] = true

	m := match.abs().conceptImplSync
	if _, seen := m.syncMap[sId]; !seen {
		m.syncMap[sId] = s.newSyncEntry()
	}
	affectedEntries[m.syncMap[sId].add(matchInstPartId, selfInstPartId)] = true
	return affectedEntries
}

func (s *conceptImplSync) typeLockSync(sourceType concept, sourceInstParts map[int]concept) {
	if entry, seen := s.syncMap[sourceType.id()]; seen {
		s.lockMap = entry.lockMap(sourceInstParts)
	}
}

func (s *conceptImplSync) typeUnlockSync() {
	s.lockMap = map[int]concept{}
}

func (a *Agent) newConceptImplSync(abs *abstractConcept) {
	abs.conceptImplSync = &conceptImplSync{
		abs:     abs,
		syncMap: map[int]*syncEntry{},
		lockMap: map[int]concept{},
	}
}

type syncEntry struct {
	data  map[int]map[int]int // my inst's part id -> their inst's part id -> occurrence
	count int                 // number of time this sync is observed
}

const syncTrustThreshold = 0.9

func (e *syncEntry) add(selfInstPartId, matchInstPartId int) *syncEntry {
	if _, seen := e.data[selfInstPartId]; !seen {
		e.data[selfInstPartId] = map[int]int{}
	}

	e.data[selfInstPartId][matchInstPartId]++
	return e
}

func (e *syncEntry) lockMap(matchInstParts map[int]concept) map[int]concept {
	result := map[int]concept{}
	for partId, data := range e.data {
		for matchId, count := range data {
			c, seen := matchInstParts[matchId]
			if seen && float64(count)/float64(e.count) > syncTrustThreshold {
				result[partId] = c
			}
		}
	}

	return result
}

func (s *conceptImplSync) newSyncEntry() *syncEntry {
	return &syncEntry{data: map[int]map[int]int{}}
}
