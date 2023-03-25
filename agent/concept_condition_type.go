package agent

type conditionType interface {
	concept
	verify(args ...any) *bool
}
