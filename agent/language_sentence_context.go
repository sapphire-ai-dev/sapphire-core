package agent

import (
	"reflect"
)

type sntcCtx struct {
	// general information
	convCtx *convCtx
	src     object
	dst     object

	// state variables for fitting
	sentence []string
	matches  map[int]map[langPart]map[*sntcFit]bool
	fitDone  map[int]map[langPart]bool // start position -> part

	// state variables interpretation
	// interpretedConds: for each language part, the language conditions its concept must have
	//   satisfied, used to narrow down pronoun match, populated when interpreting forms and used
	//   to interpret wordLangPart / conceptLangPart
	interpretedConds map[langPart]map[langCond]*bool
	newConcepts      map[int]concept
}

func (c *sntcCtx) fitStatus(p langPart, pos int) bool {
	done := false
	if m, dSeen := c.fitDone[pos]; dSeen {
		done = m[p]
	}

	return done
}

func (c *sntcCtx) setFitStatus(p langPart, pos int, done bool) {
	if _, dSeen := c.fitDone[pos]; !dSeen {
		c.fitDone[pos] = map[langPart]bool{}
	}

	c.fitDone[pos][p] = done
}

func (c *sntcCtx) addMatch(start int, p langPart, match *sntcFit) {
	if _, seen := c.matches[start]; !seen {
		c.matches[start] = map[langPart]map[*sntcFit]bool{}
	}

	if _, seen := c.matches[start][p]; !seen {
		c.matches[start][p] = map[*sntcFit]bool{}
	}

	c.matches[start][p][match] = true
}

func (c *sntcCtx) getMatch(start int, p langPart) []*sntcFit {
	if _, seen := c.matches[start]; !seen {
		return []*sntcFit{}
	}

	if _, seen := c.matches[start][p]; !seen {
		return []*sntcFit{}
	}

	var result []*sntcFit
	for match := range c.matches[start][p] {
		result = append(result, match)
	}

	return result
}

func (c *sntcCtx) filterConcepts(class reflect.Type) map[int]concept {
	candidates, result := map[int]concept{}, map[int]concept{}
	if c.src != nil {
		candidates[c.src.id()] = c.src
	}
	if c.dst != nil {
		candidates[c.dst.id()] = c.dst
	}

	for _, cc := range c.convCtx.mentioned {
		candidates[cc.id()] = cc
	}

	for _, cc := range candidates {
		if reflect.TypeOf(cc) == class {
			result[cc.id()] = cc
		}
	}

	return candidates
}

func (l *agentLanguage) newSntcCtx(src, dst object) *sntcCtx {
	return &sntcCtx{
		convCtx:          l.newConvCtx(),
		src:              src,
		dst:              dst,
		fitDone:          map[int]map[langPart]bool{},
		matches:          map[int]map[langPart]map[*sntcFit]bool{},
		interpretedConds: map[langPart]map[langCond]*bool{},
		newConcepts:      map[int]concept{},
	}
}
