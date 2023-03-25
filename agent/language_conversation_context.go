package agent

type convCtx struct {
	agent     *Agent
	mentioned map[int]concept
}

func (l *agentLanguage) newConvCtx() *convCtx {
	return &convCtx{
		agent:     l.agent,
		mentioned: map[int]concept{},
	}
}
