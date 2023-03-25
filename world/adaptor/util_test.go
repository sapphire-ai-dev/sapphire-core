package adaptor

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func assertEmptyNotNil(t *testing.T, a any) {
	assert.Empty(t, a)
	assert.NotNil(t, a)
}

var testWorldName = "test"

type testWorld struct {
	resetCalled      int
	newActorCalled   int
	newActorArgs     []any
	newActorReturnId int
	lookCalled       int
	lookActorId      int
	lookReturnId     int
	feelCalled       int
	feelActorId      int
	feelReturnId     int
	cmdCalled        int
	cmdArgs          []any
}

func (w *testWorld) Name() string {
	return testWorldName
}

func (w *testWorld) Reset() {
	w.resetCalled++
}

func (w *testWorld) Tick() {}

func (w *testWorld) NewActor(args ...any) (int, []*world.ActionInterface) {
	w.newActorCalled++
	w.newActorArgs = args
	w.newActorReturnId = rand.Intn(1 << 20)
	return w.newActorReturnId, nil
}

func (w *testWorld) Register(_ int, _ func()) {}

func (w *testWorld) Look(actorId int) []*world.Image {
	w.lookCalled++
	w.lookActorId = actorId
	w.lookReturnId = rand.Intn(1 << 20)
	return []*world.Image{{Id: w.lookReturnId}}
}

func (w *testWorld) Feel(actorId int) []*world.Touch {
	w.feelCalled++
	w.feelActorId = actorId
	w.feelReturnId = rand.Intn(1 << 20)
	return []*world.Touch{{Id: w.feelReturnId}}
}

func (w *testWorld) Cmd(args ...any) {
	w.cmdCalled++
	w.cmdArgs = args
}
