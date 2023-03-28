package agent

type agentSymbolicRecord struct {
	agent    *Agent
	numerics *numericRecord
}

func (a *Agent) newAgentSymbolicRecord() {
	a.symbolic = &agentSymbolicRecord{
		agent: a,
	}

	a.symbolic.newNumerics()
}

type numericRecord struct {
	number0 *number
	number1 *number
}

func (s *agentSymbolicRecord) newNumerics() {
	s.numerics = &numericRecord{
		number0: s.agent.newNumber(0, nil),
		number1: s.agent.newNumber(1, nil),
	}
}
