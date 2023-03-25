package agent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLangNodeConstructor(t *testing.T) {
	agent := newEmptyWorldAgent()
	class := toReflect[*simpleObject]()
	soLn := agent.language.newLangNode(class)
	assert.Equal(t, agent, soLn.agent)
	assert.Equal(t, class, soLn.class)
	assert.NotNil(t, soLn.conds)
}
