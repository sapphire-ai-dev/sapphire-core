package agent

import "strings"

type conceptCpntLang interface {
	explicitName(tenseId int) string
	setExplicitName(tenseId int, name string)
	toSentence(ctx *sntcCtx) string
}

type conceptImplLang struct {
	abs    *abstractConcept
	tenses map[int]string
}

func (l *conceptImplLang) explicitName(tenseId int) string {
	return l.tenses[tenseId]
}

func (l *conceptImplLang) setExplicitName(tenseId int, name string) {
	l.tenses[tenseId] = name
}

func (l *conceptImplLang) toSentence(ctx *sntcCtx) string {
	return strings.Join(l.abs.agent.language.toSentence(l.abs.self(), ctx).str(), " ")
}

func (a *Agent) newConceptImplLang(abs *abstractConcept) {
	result := &conceptImplLang{
		abs:    abs,
		tenses: map[int]string{},
	}

	abs.conceptImplLang = result
}
