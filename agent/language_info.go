package agent

// used to track what information should be included under what conditions
// if it is not included, this structure also stores the implicit backup to assume
// for example, if user says "I want you to eat an apple", the "want" should include time information
// in its absence, typically assume there is an implicit "now"
// note: implicit id of 0 indicates explicit
type langInfo struct {
	node        *langNode
	partId      int         // the information to include
	implicitIds map[int]int // the implicit information to assume -> occurrence
}

func (i *langInfo) record(implicitId int) {
	if _, seen := i.implicitIds[implicitId]; !seen {
		i.implicitIds[implicitId] = 0
	}
	i.implicitIds[implicitId]++
}

func (n *langNode) newLangInfo(partId int) *langInfo {
	if li, seen := n.infos[partId]; seen {
		return li
	}

	result := &langInfo{
		node:        n,
		partId:      partId,
		implicitIds: map[int]int{},
	}

	n.infos[partId] = result
	return result
}
