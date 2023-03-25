package world

const InfoLabelObservable = "observable"

const TernaryPos = "[pos]"
const TernaryZro = "[zro]"
const TernaryNeg = "[neg]"

/*
Info

	# used to feed general information into aspectNodes in the agent package

	# fields:
	    # Labels: used to categorize qualitative information
	        # example: ["direction", "left", "positive"]
	    # Value: used to store quantitative information
	        # example: THREE steps
*/
type Info struct {
	Labels []string
	Value  any
}

/*
Image

	# used to collect information corresponding to a single unit

	# fields:
	    # Id: the unit Id
	    # Permanent: collection of permanent information, i.e. appearance
	    # Transient: collection of transient information, i.e. location
*/
type Image struct {
	Id        int
	Name      string
	Permanent []*Info
	Transient []*Info
}

/*
Touch

	# used to send action response information back to actor

	# fields:
	    # Id: the unit Id
	    # Info: to store contact information
*/
type Touch struct {
	Id   int
	Name string
	Info *Info
}
