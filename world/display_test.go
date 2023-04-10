package world

import (
	"strconv"
	"testing"
	"time"
)

type testWorld struct {
	*AbstractWorld
}

func (w *testWorld) Name() string {
	return "test"
}

func (w *testWorld) NewActor(args ...any) (int, []*ActionInterface) {
	//TODO implement me
	panic("implement me")
}

func (w *testWorld) Look(actorId int) []*Image {
	//TODO implement me
	panic("implement me")
}

func (w *testWorld) Feel(actorId int) []*Touch {
	//TODO implement me
	panic("implement me")
}

func (w *testWorld) Cmd(args ...any) {
	//TODO implement me
	panic("implement me")
}

func TestNewDisplayClient(t *testing.T) {
	dc := NewDisplayClient(&testWorld{NewAbstractWorld()})

	for i := 0; i < 1; i++ {
		time.Sleep(time.Second)
		dc.SendState([]byte(strconv.Itoa(i)))
	}
}
