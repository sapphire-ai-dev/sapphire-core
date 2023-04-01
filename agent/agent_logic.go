package agent

type agentLogicRecord struct {
	agent        *Agent
	comparatives *comparativeRecord
	intervals    *intervalRecord
}

func (a *Agent) newAgentLogicRecord() {
	a.logic = &agentLogicRecord{
		agent: a,
	}

	a.logic.newComparatives()
	a.logic.newIntervals()
}

type comparativeRecord struct {
	gt *comparativeRelationType // greater than
	ge *comparativeRelationType // greater than or equal to
	eq *comparativeRelationType // equal to
	le *comparativeRelationType // less than or equal to
	lt *comparativeRelationType // less than
}

func (s *agentLogicRecord) newComparatives() {
	s.comparatives = &comparativeRecord{}
	s.comparatives.gt = s.agent.newComparativeRelationType(comparativeIdGT, nil)
	s.comparatives.ge = s.agent.newComparativeRelationType(comparativeIdGE, nil)
	s.comparatives.eq = s.agent.newComparativeRelationType(comparativeIdEQ, nil)
	s.comparatives.le = s.agent.newComparativeRelationType(comparativeIdLE, nil)
	s.comparatives.lt = s.agent.newComparativeRelationType(comparativeIdLT, nil)
}

type intervalRecord struct {
	gt *intervalRelationType // greater than
	eq *intervalRelationType // equal to
	lt *intervalRelationType // less than
	ct *intervalRelationType // contains
}

func (s *agentLogicRecord) newIntervals() {
	s.intervals = &intervalRecord{}
	s.intervals.gt = s.agent.newIntervalRelationType(intervalIdGT, nil)
	s.intervals.eq = s.agent.newIntervalRelationType(intervalIdEQ, nil)
	s.intervals.lt = s.agent.newIntervalRelationType(intervalIdLT, nil)
	s.intervals.ct = s.agent.newIntervalRelationType(intervalIdCT, nil)
}
