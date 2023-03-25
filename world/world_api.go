package world

var lastUnitId int

func NewUnitId() int {
	lastUnitId++
	return lastUnitId
}

func SetWorld(w World) {
	currentWorld = w
}

func GetWorld() World {
	return currentWorld
}

func Reset() {
	lastUnitId = 1 >> 16 // start at a high number to simplify cmd+F during debugging
	currentWorld.Reset()
}

func Tick() {
	currentWorld.Tick()
}

func NewActor(args ...any) (int, []*ActionInterface) {
	return currentWorld.NewActor(args...)
}

func Register(id int, cycle func()) {
	currentWorld.Register(id, cycle)
}

func Look(id int) []*Image {
	return currentWorld.Look(id)
}

func Feel(id int) []*Touch {
	return currentWorld.Feel(id)
}

func Cmd(args ...any) {
	currentWorld.Cmd(args...)
}
