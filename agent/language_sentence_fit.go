package agent

import (
	"fmt"
	"reflect"
)

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

func (f *sntcFit) debug() string {
	return f.debugHelper("")
}

func (f *sntcFit) debugHelper(indent string) string {
	var c any
	if isNil(f.c) {
		c = "nil"
	} else {
		c = fmt.Sprint(reflect.TypeOf(f.c), f.c.id())
	}

	result := fmt.Sprint(indent, c, " ", f.lang.debug())
	for _, child := range f.children {
		result += child.debugHelper(indent + tab)
	}

	return result
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
