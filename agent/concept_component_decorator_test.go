package agent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConceptCtxMatch(t *testing.T) {
	agent := newEmptyWorldAgent()
	tc := agent.newTestConcept(1, nil)

	ccat := agent.newCreateContextActionType()
	cca := agent.newCreateContextAction(ccat, agent.self, 0)
	cot := agent.newContextObjectType(conceptSourceObservation)
	co := agent.newContextObject(cca)
	co.addType(cot)

	assert.Equal(t, tc, agent.newTestConcept(1, nil))
	assert.NotEqual(t, tc, agent.newTestConcept(1, map[int]any{conceptArgContext: co}))
}
