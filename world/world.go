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
	Cmd(args ...any)
}
