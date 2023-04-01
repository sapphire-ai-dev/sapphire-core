package agent

type agentTime struct {
	agent *Agent
	clock int
	now   temporalObject
}

func (t *agentTime) cycle() {
	t.clock++
	t.now = t.agent.newTimePointObject(&t.clock, nil)
}

// currently only supports time points with non-nil clock times
func (t *agentTime) temporalObjJoin(l, r temporalObject) temporalObject {
	ok, lsc, lec, rsc, rec, ctx := t.temporalObjBreakdownHelper(l, r)
	if !ok {
		return nil
	}

	if lsc == rsc && lsc == lec && rsc == rec {
		r.replace(l)
		return l
	}

	var s, e *timePointObject
	if lsc > rsc {
		s = r.start()
	} else {
		s = l.start()
	}

	if lec > rec {
		e = l.end()
	} else {
		e = r.end()
	}

	args := map[int]any{}
	if ctx != nil {
		args[conceptArgContext] = ctx
	}

	return t.agent.newTimeSegmentObject(s, e, args)
}

// currently only supports time points with non-nil clock times
func (t *agentTime) temporalObjCompare(l, r temporalObject) map[int]relation {
	result := map[int]relation{}

	ok, lsc, lec, rsc, rec, ctx := t.temporalObjBreakdownHelper(l, r)
	if !ok {
		return result
	}

	args := map[int]any{}
	if ctx != nil {
		args[conceptArgContext] = ctx
	}

	if lsc == lec && rsc == rec { // both are time points
		if lsc > rec {
			cr := t.agent.newComparativeRelation(t.agent.logic.comparatives.gt, l, r, args)
			result[cr.id()] = cr
		} else if lec < rsc {
			cr := t.agent.newComparativeRelation(t.agent.logic.comparatives.lt, l, r, args)
			result[cr.id()] = cr
		} else if lsc == rsc {
			cr := t.agent.newComparativeRelation(t.agent.logic.comparatives.eq, l, r, args)
			result[cr.id()] = cr
		}
	} else { // at least one is a time segment
		if lsc == rsc && lec == rec {
			cr := t.agent.newIntervalRelation(t.agent.logic.intervals.eq, l, r, args)
			result[cr.id()] = cr
		} else if lsc <= rsc && lec >= rec {
			cr := t.agent.newIntervalRelation(t.agent.logic.intervals.ct, l, r, args)
			result[cr.id()] = cr
		} else if lsc >= rsc && lec <= rec {
			cr := t.agent.newIntervalRelation(t.agent.logic.intervals.ct, r, l, args)
			result[cr.id()] = cr
		} else if lsc >= rec {
			cr := t.agent.newIntervalRelation(t.agent.logic.intervals.gt, l, r, args)
			result[cr.id()] = cr
		} else if lec <= rsc {
			cr := t.agent.newIntervalRelation(t.agent.logic.intervals.lt, l, r, args)
			result[cr.id()] = cr
		}
	}

	return result
}

func (t *agentTime) temporalObjBreakdownHelper(l, r temporalObject) (bool, int, int, int, int, *contextObject) {
	if isNil(l) || isNil(r) || !matchConcepts(l.ctx(), r.ctx()) {
		return false, 0, 0, 0, 0, nil
	}

	ls, le, rs, re := l.start(), l.end(), r.start(), r.end()
	if ls.clockTime == nil || le.clockTime == nil || rs.clockTime == nil || re.clockTime == nil {
		return false, 0, 0, 0, 0, nil
	}

	return true, *ls.clockTime, *le.clockTime, *rs.clockTime, *re.clockTime, l.ctx()
}

func (a *Agent) newAgentTime() {
	a.time = &agentTime{
		agent: a,
		clock: 0,
	}

	a.time.now = a.newTimePointObject(&a.time.clock, nil)
}
