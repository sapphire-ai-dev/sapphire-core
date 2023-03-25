package agent

type sntcPart interface {
	str() []string
	parent() langPart
}

type wordSntcPart struct {
	s string
	p langPart
}

func (w *wordSntcPart) str() []string {
	return []string{w.s}
}

func (w *wordSntcPart) parent() langPart {
	return w.p
}
