package agent

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
	"github.com/sapphire-ai-dev/sapphire-core/world/text"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAgentPerceptionConstructor(t *testing.T) {
	agent := newTextWorldAgent()
	assert.Empty(t, agent.perception.visibleObjects)
	assert.NotNil(t, agent.perception.visibleObjects)
}

func TestAgentPerceptionLook(t *testing.T) {
	agent := newTextWorldAgent()
	rootId, dId := 0, 0
	world.Cmd(text.CmdTypeGetRootDirectoryId, &rootId)
	world.Cmd(text.CmdTypeCreateDirectory, rootId, "", &dId)
	agent.cycle()

	assert.NotZero(t, dId)
	assert.Equal(t, dId, agent.perception.visibleObjects[0].(*simpleObject).worldId)
}
