package agent

import "strings"

type agentSymbolicRecord struct {
	agent     *Agent
	numerics  *numericRecord
	breakdown *symbolicBreakdownCpnt
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

type symbolicBreakdownCpnt struct {
	strategies map[int]func(s string) []string
}

func (c *symbolicBreakdownCpnt) breakdown(s string) map[int][]string {
	result := map[int][]string{}
	for strategyId, strategy := range c.strategies {
		result[strategyId] = strategy(s)
	}

	return result
}

func (s *agentSymbolicRecord) newBreakdown() {
	s.breakdown = &symbolicBreakdownCpnt{
		strategies: map[int]func(s string) []string{},
	}

	s.breakdown.strategies[breakdownStrategyChar] = s.breakdown.breakdownStrategyChar
	s.breakdown.strategies[breakdownStrategyWord] = s.breakdown.breakdownStrategyWord
}

func (c *symbolicBreakdownCpnt) breakdownStrategyChar(s string) []string {
	var result []string
	for i := range s {
		result = append(result, s[i:i+1])
	}
	return result
}

func (c *symbolicBreakdownCpnt) breakdownStrategyWord(s string) []string {
	return strings.Split(s, " ")
}

const (
	breakdownStrategyChar = iota
	breakdownStrategyWord
)
