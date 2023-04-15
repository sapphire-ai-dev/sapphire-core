package agent

import (
	"math"
)

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
	value(inst performableAction) float64
	predictValue(args map[int]any) map[int]float64 // hypothesized action id -> hypothesized value
	predictValueHelper(args map[int]any, visitedTypes map[int]performableActionType, result map[int]float64)
	conditions() map[int]conditionType
	instantiate(args map[int]any) map[int]performableAction
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

func (a *abstractPerformableAction) instShareParts() (map[int]concept, map[int]int) {
	parts, sync := map[int]concept{partIdConceptSelf: a._self}, map[int]int{a._self.id(): partIdConceptSelf}
	if a._performer != nil {
		parts[partIdActionPerformer] = a.performer()
		sync[a._performer.c.id()] = partIdActionPerformer
	}
	if a._receiver != nil {
		parts[partIdActionReceiver] = a.receiver()
		sync[a._receiver.c.id()] = partIdActionReceiver
	}
	return parts, sync
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

func (a *Agent) newAbstractPerformableAction(self concept, args map[int]any, out **abstractPerformableAction) {
	*out = &abstractPerformableAction{_state: actionStateIdle, _snapshots: map[int]*snapshot{}}
	a.newAbstractAction(self, args, &(*out).abstractAction)
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

	if overrideMind == nil {
		overrideMind = mindConcepts[concept](a.agent.mind)
	}

	for _, c := range overrideMind {
		a._snapshots[timing].mind[c.id()] = c.createReference(a._self, false)
	}

	instParts, _ := a.instShareParts()
	for _, condType := range a.t.c.(performableActionType).conditions() {
		condType.typeLockSync(a.t.c, instParts)
		a._snapshots[timing].condTruth[condType.id()] = condType.verify(map[int]any{partIdConceptTime: a.agent.time.now})
		condType.typeUnlockSync()
	}
}

func (a *abstractPerformableAction) getSnapshot(timing int) *snapshot {
	return a._snapshots[timing]
}

type abstractPerformableActionType struct {
	*abstractActionType
	_receiverType    *memReference
	_conditions      map[int]*memReference
	_causations      map[int]*memReference
	causationRecords map[int]*causationRecord
}

func (t *abstractPerformableActionType) predictValue(args map[int]any) map[int]float64 {
	visitedTypes, result := map[int]performableActionType{}, map[int]float64{}
	t.predictValueHelper(args, visitedTypes, result)
	return result
}

func (t *abstractPerformableActionType) predictValueHelper(args map[int]any,
	visitedTypes map[int]performableActionType, result map[int]float64) {
	if _, seen := visitedTypes[t.cid]; seen {
		return
	}

	for _, inst := range t._self.(performableActionType).instantiate(args) {
		result[inst.id()] = t.value(inst)
	}
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
	_, actionInstSync := inst.instShareParts()
	log := instPrevSnapshot.condTruth

	for _, modif := range filterConcepts[modifier](t.agent, instPrevSnapshot.mind) {
		condType := modif._type()
		if _, seen := t._conditions[condType.id()]; !seen {
			t._conditions[condType.id()] = condType.createReference(t._self, false)
		}
		log[condType.id()] = ternary(true)
		_, modifSync := modif.instShareParts()
		t.typeUpdateSync(condType, actionInstSync, modifSync)
	}

	for _, rel := range filterConcepts[relation](t.agent, instPrevSnapshot.mind) {
		condType := rel._type()
		if _, seen := t._conditions[condType.id()]; !seen {
			t._conditions[condType.id()] = condType.createReference(t._self, false)
		}
		log[condType.id()] = ternary(true)
		_, relSync := rel.instShareParts()
		t.typeUpdateSync(condType, actionInstSync, relSync)
	}

	for _, causation := range filterConcepts[change](t.agent, instPostSnapshot.mind) {
		causationType := causation._type().(changeType)
		if _, seen := t._causations[causationType.id()]; !seen {
			t._causations[causationType.id()] = causationType.createReference(t._self, false)
		}
		t.newCausationRecord(causationType).addLog(log)
	}
}

func (t *abstractPerformableActionType) value(inst performableAction) float64 {
	actionInstParts, _ := inst.instShareParts()
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
	syncMap := map[int]concept{partIdConceptSelf: t._self}
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
		result[conditionId] = condition.verify(map[int]any{partIdConceptTime: t.agent.time.now})
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
