package agent

type langCond interface {
	match(other langCond) bool
	satisfied(root, parent concept, ctx *sntcCtx) *bool

	// generate implicit concepts from the forms being used, i.e. "my apple" -> generate ownershipRelation
	interpret(root, parent concept, truth *bool, ctx *sntcCtx) map[int]concept
}
