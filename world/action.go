package world

/*
ActionInterface

	# the atomic action interface type
	# each action an actor may take in a world is provided as an instance of this type
	# instantiated per world per actor per action
	# essentially, each instance of this type is an active ability for a game character

	# fields:
	    # Name: the name of the action interface, used for debugging only
	    # Ready: determine whether it is currently legal to perform this action
	    # Step: perform the action
*/
type ActionInterface struct {
	Name  string
	Ready func() bool
	Step  func()
}
