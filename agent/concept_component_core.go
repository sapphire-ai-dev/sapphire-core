package agent

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type conceptCpntCore interface {
	part(partId int) concept
	setPart(partId int, c concept)
	cycle()
	debug(indent string, depth int) string
	debugArgs() map[string]any
}

type conceptImplCore struct {
	abs *abstractConcept
}

// to be implemented per class
func (c *conceptImplCore) part(_ int) concept {
	return nil
}

// to be implemented per class
func (c *conceptImplCore) setPart(_ int, _ concept) {}

// to be implemented per class
func (c *conceptImplCore) cycle() {}

func (c *abstractConcept) debug(indent string, depth int) string {
	return c.debugImpl(indent, depth, c._self.debugArgs())
}

// to be implemented per class
func (c *abstractConcept) debugArgs() map[string]any {
	return map[string]any{
		"tenses": c.conceptImplLang.tenses,
	}
}

var tab = "  "

func (c *abstractConcept) debugImpl(indent string, depth int, args map[string]any) string {
	result := fmt.Sprintf("[%s] %d", reflect.TypeOf(c._self).String(), c.cid)
	if depth < 1 {
		return result
	}

	sortedArgNames := organizePrintArgs(args)
	body := ""
	for _, organizedArg := range sortedArgNames {
		argName, arg := organizedArg.name, organizedArg.arg
		if memRef, isMemRef := arg.(*memReference); isMemRef {
			if memRef == nil {
				continue
			}
			body += fmt.Sprintf("%s%s: %s\n", indent+tab, argName, memRef.debug(indent+tab, depth-1))
			continue
		}

		if mapRef, isMapRef := arg.(map[int]*memReference); isMapRef {
			if len(mapRef) == 0 {
				//body += fmt.Sprintf("%s%s: []\n", indent+tab, argName)
				continue
			}

			body += fmt.Sprintf("%s%s: [\n", indent+tab, argName)
			for _, mapArg := range mapRef {
				body += fmt.Sprintf("%s%s\n", indent+tab+tab, mapArg.debug(indent+tab+tab, depth-1))
			}
			body += fmt.Sprintf("%s]\n", indent+tab)
			continue
		}

		body += fmt.Sprintf("%s%s: %v\n", indent+tab, argName, arg)
	}

	if len(strings.Split(body, "\n")) > 2 {
		result += ": {\n"
		result += body
		result += fmt.Sprintf("%s}", indent)
	} else {
		result += trimSpaces(body)
	}
	return result
}

func organizePrintArgs(args map[string]any) []*printArg {
	var result []*printArg

	for argName, arg := range args {
		result = append(result, &printArg{
			name: argName,
			arg:  arg,
		})
	}

	// order by type (primitive -> single reference -> map of references)
	// then alphabetical
	sort.SliceStable(result, func(i int, j int) bool {
		_, iIsMapRef := result[i].arg.(map[int]*memReference)
		_, jIsMapRef := result[j].arg.(map[int]*memReference)
		_, iIsMemRef := result[i].arg.(*memReference)
		_, jIsMemRef := result[j].arg.(*memReference)
		if iIsMapRef && !jIsMapRef {
			return false
		}
		if !iIsMapRef && jIsMapRef {
			return true
		}
		if iIsMemRef && !jIsMemRef {
			return false
		}
		if !iIsMemRef && jIsMemRef {
			return true
		}
		return strings.Compare(result[i].name, result[j].name) < 0
	})
	return result
}

type printArg struct {
	name string
	arg  any
}

func trimSpaces(s string) string {
	if len(s) == 0 {
		return s
	}
	trimmed := 0
	for s[0] == ' ' {
		s = s[1:]
		trimmed++
	}
	if trimmed > 0 {
		s = " " + s
	}
	for s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	return s
}

func (a *Agent) newConceptImplCore(abs *abstractConcept) {
	abs.conceptImplCore = &conceptImplCore{abs: abs}
}
