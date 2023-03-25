package agent

type sntcNode struct {
	ln    *langNode
	lf    *langForm
	parts []sntcPart
}

func (n *sntcNode) str() []string {
	var result []string
	for _, p := range n.parts {
		result = append(result, p.str()...)
	}

	return result
}

func (n *sntcNode) parent() langPart {
	return n.lf
}

func (n *langNode) newSntcNode(lf *langForm) *sntcNode {
	return &sntcNode{
		ln:    n,
		lf:    lf,
		parts: []sntcPart{},
	}
}
