package text

import (
	"errors"
	"github.com/sapphire-ai-dev/sapphire-core/world"
)

type textWorld struct {
	*world.AbstractWorld
	rootDirectory *directory
	items         map[int]item
	actors        map[int]*actorPos
}

type actorPos struct {
	currItemId int
	cursorItem int
	cursorLine int
	cursorChar int
}

func (w *textWorld) newActorPos() *actorPos {
	return &actorPos{
		currItemId: w.rootDirectory.id(),
		cursorItem: 0,
		cursorLine: 0,
		cursorChar: 0,
	}
}

func (w *textWorld) Name() string {
	return "text"
}

func (w *textWorld) Reset() {
	w.AbstractWorld.Reset()
	w.items = map[int]item{}
	w.actors = map[int]*actorPos{}
	w.rootDirectory = &directory{
		content: []item{},
	}

	w.newAbstractItem(w.rootDirectory, nil, "", &w.rootDirectory.abstractItem)
}

func (w *textWorld) NewActor(_ ...any) (int, []*world.ActionInterface) {
	id := world.NewUnitId()
	w.actors[id] = w.newActorPos()
	return id, w.newActionInterfaces(id)
}

var (
	errActorNotFound = errors.New("actor not found")
)

func (w *textWorld) Register(id int, cycle func()) {
	if _, seen := w.actors[id]; !seen {
		panic(errActorNotFound)
	}

	w.AbstractWorld.Register(id, cycle)
}

func (w *textWorld) Look(id int) []*world.Image {
	actor, actorSeen := w.actors[id]
	if !actorSeen {
		return []*world.Image{}
	}
	currItem, itemSeen := w.items[actor.currItemId]
	if !itemSeen {
		return []*world.Image{}
	}

	result := currItem.fileImgs(actor.cursorLine, actor.cursorChar)
	result = append(result, currItem.dirImgs(actor.cursorItem)...)
	return result
}

func (w *textWorld) Feel(_ int) []*world.Touch {
	return []*world.Touch{}
}

func newTextWorld() *textWorld {
	result := &textWorld{
		AbstractWorld: world.NewAbstractWorld(),
	}
	result.Reset()
	return result
}

var display *world.DisplayClient

func Init() {
	w := newTextWorld()
	world.SetWorld(w)
	display = world.NewDisplayClient(w.Name())
}
