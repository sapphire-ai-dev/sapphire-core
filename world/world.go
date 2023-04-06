package world

var currentWorld World

/*
World

	# skeletal design for training environments

	# methods:
	    # Reset: resets the world to default factory state
	        # sets clock to 0
	        # eliminates all units
	    # Tick: advances world clock by 1, executes all cycle functions sequentially in registration order
	    # NewActor: creates a new actor
	        # return: unit id of actor
	        # return: list of atomic action interfaces provided by the world
	        # return: an action response to communicate outcome of invoking an action interface
	    # Register: registers a cycle function
	        # all registered cycle functions will be executed per tick
	    # Look: an actor looks, receiving a collection of images
	        # id: id of the actor that looks
	        # return: list of images the actor sees
	    # Cmd: used to enable implementation-specific commands per world
*/
type World interface {
	Name() string
	Reset()
	Tick()
	NewActor(args ...any) (int, []*ActionInterface)
	Register(actorId int, cycle func())
	Look(actorId int) []*Image
	Feel(actorId int) []*Touch
	Listen(id int) []*LangMessage
	Speak(speakerId, listenerId *int, content string)
	Cmd(args ...any)
}

type AbstractWorld struct {
	clock              int
	cycleFuncs         map[int]func()
	langReceiveBuffers map[int]*langReceiveBuffer
}

func (w *AbstractWorld) Tick() {
	for _, f := range w.cycleFuncs {
		f()
	}
}

func (w *AbstractWorld) Reset() {
	w.clock = 0
	w.cycleFuncs = map[int]func(){}
	w.langReceiveBuffers = map[int]*langReceiveBuffer{}
	w.newLangReceiveBuffer(nil) // everyone hears from this buffer
}

func (w *AbstractWorld) Register(id int, cycle func()) {
	w.cycleFuncs[id] = cycle
	w.newLangReceiveBuffer(&id)
}

func NewAbstractWorld() *AbstractWorld {
	result := &AbstractWorld{}
	result.Reset()
	return result
}
