package agent

import "github.com/sapphire-ai-dev/sapphire-core/world"

type agentAspect struct {
	agent *Agent
	root  *aspectNode
}

func (a *agentAspect) find(args ...string) *aspectNode {
	return a.root.find(args...)
}

func (a *agentAspect) ternary(name string) *aspectNode {
	return a.root.find(aspectTernaryDebugName, name)
}

func (a *agentAspect) dimension(name string) *aspectNode {
	return a.root.find(aspectDimensionDebugName, name)
}

func (a *agentAspect) interval(name string) *aspectNode {
	return a.root.find(aspectIntervalDebugName, name)
}

func (a *agentAspect) qualitative(name string) *aspectNode {
	return a.root.find(aspectQualitativeDebugName, name)
}

func (a *agentAspect) lowestCommonAncestor(l, r *aspectNode) *aspectNode {
	var lAncestors []*aspectNode
	var rAncestors []*aspectNode
	lA, rA := l, r
	for lA != nil {
		lAncestors = append(lAncestors, lA)
		lA = lA.parent
	}
	for rA != nil {
		rAncestors = append(rAncestors, rA)
		rA = rA.parent
	}

	lP, rP := len(lAncestors)-1, len(rAncestors)-1
	for lAncestors[lP] == rAncestors[rP] {
		lP, rP = lP-1, rP-1
	}

	return lAncestors[lP+1]
}

func (a *Agent) newAgentAspect() {
	a.aspect = &agentAspect{
		agent: a,
	}

	a.aspect.initAspectTree()
}

func (a *agentAspect) initAspectTree() {
	a.root = &aspectNode{
		parent:   nil,
		children: []*aspectNode{},
		name:     aspectRootDebugName,
	}

	for category, nodes := range aspectConstants {
		categoryRoot := a.root.addChild(category)
		for _, node := range nodes {
			categoryRoot.addChild(node)
		}
	}

	a.root.addChild(aspectObjectInfoDebugName)
}

type aspectNode struct {
	parent   *aspectNode
	children []*aspectNode
	name     string
}

func (n *aspectNode) addChild(name string) *aspectNode {
	newNode := &aspectNode{
		parent:   n,
		children: []*aspectNode{},
		name:     name,
	}

	for _, child := range n.children {
		if child.name == name {
			return child
		}
	}

	n.children = append(n.children, newNode)
	return newNode
}

func (n *aspectNode) find(args ...string) *aspectNode {
	if len(args) == 0 {
		return n
	}

	for _, child := range n.children {
		if child.name == args[0] {
			return child.find(args[1:]...)
		}
	}

	return n.addChild(args[0]).find(args[1:]...)
}

func (n *aspectNode) toString() string {
	if n.parent == nil {
		return n.name
	}

	return n.parent.toString() + " " + n.name
}

const (
	aspectRootDebugName        = "[aspectNode]"
	aspectObjectInfoDebugName  = "[objectInfo]"
	aspectTernaryDebugName     = "[ternary]"
	aspectDimensionDebugName   = "[dimension]"
	aspectIntervalDebugName    = "[interval]"
	aspectQualitativeDebugName = "[qualitative]"
)

const (
	aspectDimensionTime = "[time]"
)

const (
	aspectIntervalStartStart = "[start-start]"
	aspectIntervalStartEnd   = "[start-end]"
	aspectIntervalEndStart   = "[end-start]"
	aspectIntervalEndEnd     = "[end-end]"
)

const (
	aspectQualitativeEqual     = "[equal]"
	aspectQualitativeDifferent = "[different]"
	aspectQualitativeOpposite  = "[opposite]"
)

var aspectConstants = map[string][]string{
	aspectObjectInfoDebugName: {},
	aspectTernaryDebugName: {
		world.TernaryPos,
		world.TernaryZro,
		world.TernaryNeg,
	},
	aspectDimensionDebugName: {
		aspectDimensionTime,
	},
	aspectIntervalDebugName: {
		aspectIntervalStartStart,
		aspectIntervalStartEnd,
		aspectIntervalEndStart,
		aspectIntervalEndEnd,
	},
	aspectQualitativeDebugName: {
		aspectQualitativeEqual,
		aspectQualitativeDifferent,
		aspectQualitativeOpposite,
	},
}
