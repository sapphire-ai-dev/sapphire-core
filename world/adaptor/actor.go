package adaptor

import "github.com/sapphire-ai-dev/sapphire-core/world"

type actor struct {
	w     *adaptorWorld
	id    int
	links map[int]*link // child world id -> link
}

func (a *actor) collectActionInterfaces(argsMap map[int][]any) []*world.ActionInterface {
	var result []*world.ActionInterface
	for childWorldId, childWorld := range a.w.children {
		childActorId, childActions := childWorld.NewActor(argsMap[childWorldId]...)
		result = append(result, childActions...)
		a.links[childWorldId] = a.newLink(childWorldId, childActorId)
	}

	return result
}

func (a *actor) look() []*world.Image {
	var result []*world.Image
	for childWorldId, childWorld := range a.w.children {
		result = append(result, childWorld.Look(a.links[childWorldId].childActorId)...)
	}

	return result
}

func (a *actor) feel() []*world.Touch {
	var result []*world.Touch
	for childWorldId, childWorld := range a.w.children {
		result = append(result, childWorld.Feel(a.links[childWorldId].childActorId)...)
	}

	return result
}

func (w *adaptorWorld) newActor() int {
	id := world.NewUnitId()
	w.actors[id] = &actor{
		w:     w,
		id:    id,
		links: map[int]*link{},
	}

	return id
}

type link struct {
	actor        *actor
	childWorldId int
	childActorId int
}

func (a *actor) newLink(childWorldId, childActorId int) *link {
	return &link{
		actor:        a,
		childWorldId: childWorldId,
		childActorId: childActorId,
	}
}
