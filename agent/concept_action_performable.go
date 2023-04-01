package agent

import "math"

type performableAction interface {
	action
	start() bool
	step() bool
	state() int
	receiver() object
	setReceiver(o object)
	snapshot(timing int, overrideMind map[int]concept)
	getSnapshot(timing int) *snapshot
}

type performableActionType interface {
	actionType
	value() float64
	conditions() map[int]conditionType
	instantiate(args ...any) performableAction
	receiverType() objectType
	update(inst performableAction)
}

type abstractPerformableAction struct {
	*abstractAction
	_state     int
	_receiver  *memReference
	_snapshots map[int]*snapshot
}

func (a *abstractPerformableAction) match(o *abstractPerformableAction) bool {
	return a.abstractAction.match(o.abstractAction)
}

func (a *abstractPerformableAction) debugArgs() map[string]any {
	args := a.abstractAction.debugArgs()
	args["state"] = actionStateNames[a._state]
	args["receiver"] = a._receiver
	return args
}

func (a *abstractPerformableAction) part(partId int) concept {
	if partId == partIdActionReceiver {
		return a.receiver()
	}
	return a.abstractAction.part(partId)
}

func (a *abstractPerformableAction) state() int {
	return a._state
}

func (a *abstractPerformableAction) complete() {
	a.agent.activity.completedActions = append(a.agent.activity.completedActions, a._self.(performableAction))
}

func (a *abstractPerformableAction) instShareParts() map[int]int {
	sync := map[int]int{}
	if a._performer != nil {
		sync[a._performer.c.id()] = partIdActionPerformer
	}
	if a._receiver != nil {
		sync[a._receiver.c.id()] = partIdActionReceiver
	}
	return sync
}

func (a *abstractPerformableAction) receiver() object {
	return parseRef[object](a.agent, a._receiver)
}

func (a *abstractPerformableAction) setReceiver(o object) {
	if a._receiver != nil {
		return
	}
	a._receiver = o.createReference(a._self, false)
}

func (a *Agent) newAbstractPerformableAction(self concept, t actionType, performer object, args map[int]any,
	out **abstractPerformableAction) {
	*out = &abstractPerformableAction{_state: actionStateIdle, _snapshots: map[int]*snapshot{}}
	a.newAbstractAction(self, t, performer, args, &(*out).abstractAction)
}

type snapshot struct {
	mind      map[int]*memReference
	condTruth map[int]*bool
}

const (
	snapshotTimingPrev = iota
	snapshotTimingPost
)

func (a *abstractPerformableAction) snapshot(timing int, overrideMind map[int]concept) {
	a._snapshots[timing] = &snapshot{
		mind:      map[int]*memReference{},
		condTruth: map[int]*bool{},
	}

	if overrideMind != nil {
		overrideMind = mindConcepts[concept](a.agent.mind)
	}

	for _, c := range overrideMind {
		a._snapshots[timing].mind[c.id()] = c.createReference(a._self, false)
	}

	instParts := map[int]concept{}
	instParts[partIdActionPerformer] = a.performer()
	instParts[partIdActionReceiver] = a.receiver()

	for _, condType := range a.t.c.(performableActionType).conditions() {
		condType.typeLockSync(a.t.c, instParts)
		a._snapshots[timing].condTruth[condType.id()] = condType.verify()
		condType.typeUnlockSync()
	}
}

func (a *abstractPerformableAction) getSnapshot(timing int) *snapshot {
	return a._snapshots[timing]
}

type abstractPerformableActionType struct {
	*abstractActionType
	_conditions      map[int]*memReference
	_causations      map[int]*memReference
	causationRecords map[int]*causationRecord
	_receiverType    *memReference
}

func (t *abstractPerformableActionType) conditions() map[int]conditionType {
	return parseRefs[conditionType](t.agent, t._conditions)
}

func (t *abstractPerformableActionType) causations() map[int]changeType {
	return parseRefs[changeType](t.agent, t._causations)
}

func (t *abstractPerformableActionType) receiverType() objectType {
	return parseRef[objectType](t.agent, t._receiverType)
}

func (t *abstractPerformableActionType) match(o *abstractPerformableActionType) bool {
	if t._receiverType == nil && o._receiverType != nil || t._receiverType != nil && o._receiverType == nil {
		return false
	}

	return t.abstractActionType.match(o.abstractActionType) &&
		((t._receiverType == nil && o._receiverType == nil) || (t._receiverType.c == o._receiverType.c))
}

func (t *abstractPerformableActionType) update(inst performableAction) {
	instPrevSnapshot, instPostSnapshot := inst.getSnapshot(snapshotTimingPrev), inst.getSnapshot(snapshotTimingPost)
	actionInstPartIds := inst.instShareParts()
	log := instPrevSnapshot.condTruth

	for _, modif := range filterConcepts[modifier](t.agent, instPrevSnapshot.mind) {
		condType := modif._type()
		if _, seen := t._conditions[condType.id()]; !seen {
			t._conditions[condType.id()] = condType.createReference(t._self, false)
		}
		log[condType.id()] = ternary(true)
		t.typeUpdateSync(condType, actionInstPartIds, modif.instShareParts())
	}

	for _, rel := range filterConcepts[relation](t.agent, instPrevSnapshot.mind) {
		condType := rel._type()
		if _, seen := t._conditions[condType.id()]; !seen {
			t._conditions[condType.id()] = condType.createReference(t._self, false)
		}
		log[condType.id()] = ternary(true)
		t.typeUpdateSync(condType, actionInstPartIds, rel.instShareParts())
	}

	for _, causation := range filterConcepts[change](t.agent, instPostSnapshot.mind) {
		causationType := causation._type().(changeType)
		t.newCausationRecord(causationType).addLog(log)
	}
}

func (t *abstractPerformableActionType) value() float64 {
	actionInstParts := t.searchInstParts()
	if actionInstParts == nil {
		return math.Inf(-1)
	}

	result := 0.0
	conditionTruth := t.conditionTruth(actionInstParts)
	for causationId, record := range t.causationRecords {
		prediction := record.predict(conditionTruth)
		if prediction != nil && *prediction {
			result += t._causations[causationId].c.(changeType).value()
		}
	}

	return result
}

func (t *abstractPerformableActionType) searchInstParts() map[int]concept {
	syncMap := map[int]concept{}
	syncMap[partIdActionPerformer] = t.agent.self

	if t._receiverType != nil {
		for _, objInst := range t.agent.perception.visibleObjects {
			if _, seen := objInst.types()[t.receiverType().id()]; seen {
				syncMap[partIdActionReceiver] = objInst
				break
			}
		}

		// action type has a receiver type, action does not have a receiver, mismatch
		if _, seen := syncMap[partIdActionReceiver]; !seen {
			return nil
		}
	}

	return syncMap
}

func (t *abstractPerformableActionType) conditionTruth(syncMap map[int]concept) map[int]*bool {
	result := map[int]*bool{}
	for conditionId, condition := range t.conditions() {
		condition.typeLockSync(t._self, syncMap)
		result[conditionId] = condition.verify()
		condition.typeUnlockSync()
	}

	return result
}

func (t *abstractPerformableActionType) updateHelperAdd(m map[int]*memReference, occ map[int]int, c concept) {
	if _, seen := m[c.id()]; !seen {
		m[c.id()] = c.createReference(t._self, false)
		occ[c.id()] = 0
	}
	occ[c.id()]++
}

func (t *abstractPerformableActionType) debugArgs() map[string]any {
	args := t.abstractActionType.debugArgs()
	args["conditions"] = t._conditions
	args["causations"] = t._causations
	args["receiverType"] = t._receiverType
	return args
}

func (a *Agent) newAbstractPerformableActionType(self concept, receiverType objectType, args map[int]any,
	out **abstractPerformableActionType) {
	*out = &abstractPerformableActionType{
		_conditions:      map[int]*memReference{},
		_causations:      map[int]*memReference{},
		causationRecords: map[int]*causationRecord{},
	}
	a.newAbstractActionType(self, args, &(*out).abstractActionType)
	if receiverType != nil {
		(*out)._receiverType = receiverType.createReference(self, true)
	}
}

type causationRecord struct {
	occ       int
	condOcc   map[int]int
	condTruth map[int]int
}

func (r *causationRecord) addLog(condTruth map[int]*bool) {
	r.occ++
	for condTypeId, truth := range condTruth {
		if truth != nil {
			r.condOcc[condTypeId]++
			if *truth {
				r.condTruth[condTypeId]++
			}
		}
	}
}

const causationPredictionThreshold = 0.5

func (r *causationRecord) predict(condTruth map[int]*bool) *bool {
	expectedCondTruth, expectedWeight, totalWeight := map[int]bool{}, map[int]int{}, 0
	for condTypeId := range r.condTruth {
		ratio := float64(r.condTruth[condTypeId]) / float64(r.condOcc[condTypeId])
		if ratio > condTruthThreshold {
			expectedCondTruth[condTypeId] = true
			expectedWeight[condTypeId] = r.condOcc[condTypeId]
			totalWeight += r.condOcc[condTypeId]
		} else if 1-ratio > condTruthThreshold {
			expectedCondTruth[condTypeId] = false
			expectedWeight[condTypeId] = r.condOcc[condTypeId]
			totalWeight += r.condOcc[condTypeId]
		}
	}

	score := 0
	for condTypeId, truth := range condTruth {
		expectedTruth, seen := expectedCondTruth[condTypeId]
		if !seen || truth == nil {
			continue
		}

		if expectedTruth == *truth {
			score += expectedWeight[condTypeId]
		} else {
			score -= expectedWeight[condTypeId]
		}
	}

	scoreRatio := float64(score) / float64(totalWeight)
	if scoreRatio > causationPredictionThreshold {
		return ternary(true)
	} else if scoreRatio < causationPredictionThreshold {
		return ternary(false)
	}

	return nil
}

func (t *abstractPerformableActionType) newCausationRecord(causation changeType) *causationRecord {
	if _, seen := t.causationRecords[causation.id()]; !seen {
		t.causationRecords[causation.id()] = &causationRecord{
			occ:       0,
			condTruth: map[int]int{},
			condOcc:   map[int]int{},
		}
	}

	return t.causationRecords[causation.id()]
}

const (
	_ = iota
	actionStateIdle
	actionStateActive
	actionStateDone
)

var actionStateNames = map[int]string{
	actionStateIdle:   "[idle]",
	actionStateActive: "[active]",
	actionStateDone:   "[done]",
}
