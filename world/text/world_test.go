package text

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTextWorldInit(t *testing.T) {
	Init()
}

func TestTextWorldConstructor(t *testing.T) {
	w := newTextWorld()
	assert.NotNil(t, w.rootDirectory)
	assert.Len(t, w.items, 1)
	assertEmptyNotNil(t, w.actors)
	assertEmptyNotNil(t, w.cycleFuncs)
	assert.Equal(t, w.Name(), "text")
}

func TestTextWorldPlaceholderFuncs(t *testing.T) {
	Init()
	world.Feel(0)
	world.Reset()
}

func TestTextWorldNewActor(t *testing.T) {
	Init()
	actorId1, _ := world.NewActor()
	actorId2, _ := world.NewActor()
	assert.NotEqual(t, actorId1, actorId2)
}

func TestTextWorldRegister(t *testing.T) {
	Init()
	cycleResult := 0
	cycleFunc := func() {
		cycleResult++
	}
	assert.PanicsWithError(t, errActorNotFound.Error(), func() {
		world.Register(0, cycleFunc)
	})

	world.Tick()
	assert.Equal(t, cycleResult, 0)

	actorId, _ := world.NewActor()
	assert.NotPanics(t, func() {
		world.Register(actorId, cycleFunc)
	})

	world.Tick()
	assert.Equal(t, cycleResult, 1)
}

func TestTextWorldLook(t *testing.T) {
	w := newTextWorld()
	assert.Empty(t, w.Look(0))

	actorId, _ := w.NewActor()
	assert.Empty(t, w.Look(actorId))

	w.actors[actorId].currItemId++
	assert.Empty(t, w.Look(actorId))
}
