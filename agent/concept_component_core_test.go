package agent

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestConceptComponentCore(t *testing.T) {
	agent := newEmptyWorldAgent()
	tcVal := 42
	tc := agent.newTestConcept(tcVal, nil)
	assert.Equal(t, tc.id(), tc.cid)
	assert.Equal(t, tc.self(), tc._self)
	assert.Equal(t, tc.abs(), tc.abstractConcept)
	assert.Nil(t, tc.part(0))
	tc.cycle()
	tcStr := tc.debug("", 2)
	assert.Contains(t, tcStr, "TestConcept")
	assert.Contains(t, tcStr, strconv.Itoa(tcVal))

	assert.Equal(t, tc, agent.newTestConcept(tcVal, nil))
	assert.NotEqual(t, tc, agent.newTestConcept(tcVal+1, nil))
}
