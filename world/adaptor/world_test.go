package adaptor

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdaptorWorldConstructor(t *testing.T) {
	w := newAdaptorWorld()
	assertEmptyNotNil(t, w.actors)
	assertEmptyNotNil(t, w.children)
	assertEmptyNotNil(t, w.cycleFuncs)
}

func TestAdaptorWorldNewActor(t *testing.T) {
	w := newAdaptorWorld()
	actorId1, _ := w.NewActor()
	actorId2, _ := w.NewActor()
	assert.NotEqual(t, actorId1, actorId2)
	w.NewActor(map[int][]any{1: {1, 2, 3}})
}

func TestTextWorldRegister(t *testing.T) {
	w := newAdaptorWorld()
	cycleResult := 0
	cycleFunc := func() {
		cycleResult++
	}
	assert.PanicsWithError(t, errActorNotFound.Error(), func() {
		w.Register(0, cycleFunc)
	})

	w.Tick()
	assert.Equal(t, cycleResult, 0)

	actorId, _ := w.NewActor()
	assert.NotPanics(t, func() {
		w.Register(actorId, cycleFunc)
	})

	w.Tick()
	assert.Equal(t, cycleResult, 1)
}

func TestAdaptorWorldInit(t *testing.T) {
	tempSingleton = nil
	assert.PanicsWithError(t, errWorldNotFound.Error(), func() {
		Proxy()
	})
	assert.PanicsWithError(t, errWorldNotFound.Error(), func() {
		InitComplete()
	})

	InitStart()
	assert.PanicsWithError(t, errWorldNotFound.Error(), func() {
		Proxy()
	})

	tw := &testWorld{}
	world.SetWorld(tw)
	testWorldId := Proxy()
	assert.Equal(t, testWorldId, Proxy())
	assert.Equal(t, tempSingleton.children[testWorldId], tw)
	assert.Equal(t, tw, world.GetWorld())

	InitComplete()
	assert.Equal(t, tempSingleton, world.GetWorld())
}

func TestAdaptorWorldLifecycle(t *testing.T) {
	InitStart()
	tw := &testWorld{}
	world.SetWorld(tw)
	testWorldId := Proxy()
	InitComplete()

	assert.Contains(t, tempSingleton.Name(), "adaptor")
	assert.Contains(t, tempSingleton.Name(), testWorldName)

	assert.Zero(t, tw.newActorCalled, 0)
	newActorArgs := []any{1234, "abcd"}
	adaptorActorId, _ := world.NewActor(map[int][]any{testWorldId: newActorArgs})
	assert.Equal(t, tw.newActorCalled, 1)
	assert.Equal(t, tw.newActorReturnId, tempSingleton.actors[adaptorActorId].links[testWorldId].childActorId)
	for _, newActorArg := range newActorArgs {
		assert.Contains(t, tw.newActorArgs, newActorArg)
	}

	assert.Zero(t, tw.lookCalled)
	assert.Zero(t, tw.lookActorId)
	imgs := world.Look(adaptorActorId)
	assert.Equal(t, tw.lookCalled, 1)
	assert.Equal(t, tw.lookActorId, tempSingleton.actors[adaptorActorId].links[testWorldId].childActorId)
	assert.Equal(t, imgs[0].Id, tw.lookReturnId)

	assert.Zero(t, tw.feelCalled)
	assert.Zero(t, tw.feelActorId)
	tchs := world.Feel(adaptorActorId)
	assert.Equal(t, tw.feelCalled, 1)
	assert.Equal(t, tw.feelActorId, tempSingleton.actors[adaptorActorId].links[testWorldId].childActorId)
	assert.Equal(t, tchs[0].Id, tw.feelReturnId)

	assert.Empty(t, world.Look(1234))
	assert.Empty(t, world.Feel(1234))
}

func TestAdaptorWorldCmd(t *testing.T) {
	InitStart()
	tw := &testWorld{}
	world.SetWorld(tw)
	testWorldId := Proxy()
	InitComplete()

	assert.PanicsWithError(t, world.ErrInvalidArgs.Error(), func() {
		world.Cmd()
	})

	assert.PanicsWithError(t, world.ErrInvalidArgs.Error(), func() {
		world.Cmd("1")
	})

	assert.PanicsWithError(t, world.ErrInvalidArgs.Error(), func() {
		world.Cmd(-1)
	})

	world.Cmd(CmdTypeLocal)

	assert.PanicsWithError(t, world.ErrInvalidArgs.Error(), func() {
		world.Cmd(CmdTypeChild)
	})

	assert.PanicsWithError(t, world.ErrInvalidArgs.Error(), func() {
		world.Cmd(CmdTypeChild, "1")
	})

	assert.PanicsWithError(t, errWorldNotFound.Error(), func() {
		world.Cmd(CmdTypeChild, testWorldId+1)
	})

	assert.Zero(t, tw.cmdCalled)
	cmdArgs := []any{1234, "abcd"}
	world.Cmd(CmdTypeChild, testWorldId, cmdArgs[0], cmdArgs[1])
	assert.Equal(t, tw.cmdCalled, 1)
	for _, cmdArg := range cmdArgs {
		assert.Contains(t, tw.cmdArgs, cmdArg)
	}
}
