package agent

type langCond interface {
	match(other langCond) bool
	satisfied(root, parent concept, ctx *sntcCtx) *bool

	// generate implicit concepts from the forms being used, i.e. "my apple" -> generate ownershipRelation
	// also updates root if necessary
	interpret(root, parent concept, truth *bool, ctx *sntcCtx) (concept, map[int]concept)
}
