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

// currently assume nil start is beginning of time, nil end is end of time
// even for a time point object, if it does not provide a clock time, assume it's a time segment that spans infinitely
func (t *agentTime) temporalObjJoin(l, r temporalObject, ordered bool) temporalObject {

	ctx, commonFound := t.agent.commonCtx(l, r)
	if !commonFound {
		return nil
	}

	lp, lIsPoint := l.(*timePointObject)
	rp, rIsPoint := r.(*timePointObject)
	if lIsPoint && rIsPoint && ordered {
		if !isNil(l) && !isNil(r) && l.match(r) {
			return l
		}

		return t.agent.newTimeSegmentObject(lp, rp, map[int]any{partIdConceptContext: ctx})
	}

	if isNil(l) || isNil(r) {
		return nil
	}

	s, e := t.timePointMin(l.start(), r.start()), t.timePointMax(l.end(), r.end())
	return t.agent.newTimeSegmentObject(s, e, map[int]any{partIdConceptContext: ctx})
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
		args[partIdConceptContext] = ctx
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

func (a *Agent) commonCtx(l, r concept) (*contextObject, bool) {
	var lCtx, rCtx *contextObject
	if !isNil(l) {
		lCtx = l.ctx()
	}
	if !isNil(r) {
		rCtx = r.ctx()
	}

	if !matchConcepts(lCtx, rCtx) {
		return nil, false
	}

	return lCtx, true
}

func filterOverlapTemporal[T concept](t *agentTime, concepts map[int]T, temporal temporalObject) map[int]T {
	if temporal == nil {
		return concepts
	}

	var ts, te *int
	if !isNil(temporal.start()) {
		ts = temporal.start().clockTime
	}
	if !isNil(temporal.end()) {
		te = temporal.end().clockTime
	}

	result := map[int]T{}
	for _, c := range concepts {
		var cs, ce *int
		if !isNil(c.time()) && !isNil(c.time().start()) {
			cs = c.time().start().clockTime
		}
		if !isNil(c.time()) && !isNil(c.time().end()) {
			ce = c.time().end().clockTime
		}
		if c.time() == nil || t.overlaps(cs, ce, ts, te) {
			result[c.id()] = c
		}
	}

	return result
}

func (t *agentTime) overlaps(ls, le, rs, re *int) bool {
	leGERs := !(le != nil && rs != nil && *le < *rs)
	lsLERe := !(ls != nil && re != nil && *ls > *re)
	return leGERs && lsLERe
}

func (t *agentTime) timePointMin(l, r *timePointObject) *timePointObject {
	if isNil(l) || isNil(r) || l.clockTime == nil || r.clockTime == nil {
		return nil
	}

	if *l.clockTime <= *r.clockTime {
		return l
	}

	return r
}

func (t *agentTime) timePointMax(l, r *timePointObject) *timePointObject {
	if isNil(l) || isNil(r) || l.clockTime == nil || r.clockTime == nil {
		return nil
	}

	if *l.clockTime >= *r.clockTime {
		return l
	}

	return r
}

func (a *Agent) newAgentTime() {
	a.time = &agentTime{
		agent: a,
		clock: 0,
	}

	a.time.now = a.newTimePointObject(&a.time.clock, nil)
}
