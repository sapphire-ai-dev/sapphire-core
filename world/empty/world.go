package empty

import "github.com/sapphire-ai-dev/sapphire-core/world"

/*
emptyWorld

	# an empty implementation of world.World
	# used in tests where no environment interaction is needed
	# essentially there to prevent nil-pointer exceptions

	# methods:
	    # all implementations of world.World
*/
type emptyWorld struct{}

func (w *emptyWorld) Name() string {
	return "empty"
}

func (w *emptyWorld) Reset() {}

func (w *emptyWorld) Tick() {}

func (w *emptyWorld) NewActor(_ ...any) (int, []*world.ActionInterface) {
	return 0, nil
}

func (w *emptyWorld) Register(_ int, _ func()) {}

func (w *emptyWorld) Look(_ int) []*world.Image {
	return nil
}

func (w *emptyWorld) Feel(_ int) []*world.Touch {
	return nil
}

func (w *emptyWorld) Cmd(_ ...any) {}

func Init() {
	world.SetWorld(&emptyWorld{})
}
