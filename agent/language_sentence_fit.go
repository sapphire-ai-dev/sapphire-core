package agent

type sntcFit struct {
	start         int
	end           int
	c             concept
	sntc          sntcPart
	lang          langPart
	parent        *sntcFit
	children      []*sntcFit
	mismatchCount int
}

func newSntcFit(start, end int, sntc sntcPart, c concept, lang langPart, mismatchCount int) *sntcFit {
	return &sntcFit{
		start:         start,
		end:           end,
		c:             c,
		sntc:          sntc,
		lang:          lang,
		children:      []*sntcFit{},
		mismatchCount: mismatchCount,
	}
}
