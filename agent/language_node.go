package agent

import "reflect"

type langNode struct {
	agent *Agent
	class reflect.Type
	forms []*langForm
	log   []*sntcLog
	conds map[langCond]bool
}

func (n *langNode) selectForm(f *langForm) *langForm {
	for _, form := range n.forms {
		if form.match(f) {
			return form
		}
	}

	n.forms = append(n.forms, f)
	return f
}

// create a sentence part with a concept given some sentence context
// non-zero part Id indicate that root is that part of some parent
func (n *langNode) instantiate(root, parent concept, ctx *sntcCtx, extraConds map[langCond]bool) sntcPart {
	sortedForms := n.sortForms(root, parent, ctx, extraConds)
	if len(sortedForms) == 0 {
		return nil
	}

	return sortedForms[0].instantiate(root, ctx)
}

func (n *langNode) sortForms(root, parent concept, ctx *sntcCtx, extraConds map[langCond]bool) []*langForm {
	var scores []float64
	condTruth := n.conditionTruth(root, parent, ctx, extraConds)
	formLog := n.formLog()

	for _, form := range n.forms {
		scores = append(scores, n.formValue(formLog[form], condTruth))
	}

	return sortSlice[*langForm](n.forms, scores)
}

func (n *langNode) conditionTruth(root, parent concept, ctx *sntcCtx, extraConds map[langCond]bool) map[langCond]*bool {
	conds := n.agent.language.genConds(root, ctx)
	for cond := range extraConds {
		conds[cond] = true
	}

	result := map[langCond]*bool{}
	for cond := range conds {
		n.conds[cond] = true
		result[cond] = cond.satisfied(root, parent, ctx)
	}

	return result
}

func (n *langNode) formLog() map[*langForm][]*sntcLog {
	result := map[*langForm][]*sntcLog{}
	for _, entry := range n.log {
		if _, seen := result[entry.form]; !seen {
			result[entry.form] = []*sntcLog{}
		}
		result[entry.form] = append(result[entry.form], entry)
	}
	return result
}

func (n *langNode) formValue(subLog []*sntcLog, condTruth map[langCond]*bool) float64 {
	condValues := map[langCond]float64{}
	for cond := range condTruth {
		condValues[cond] = 0.0
	}

	for _, entry := range subLog {
		for cond, truth := range entry.condTruth {
			if _, seen := condTruth[cond]; seen && ternaryEqual(truth, condTruth[cond]) {
				condValues[cond] += entry.value
			}
		}
	}

	result := 0.0
	for cond := range condTruth {
		result += condValues[cond]
	}

	return result
}

func (n *langNode) addLog(conds map[langCond]*bool, form *langForm, value float64) {
	f := n.selectForm(form)
	n.log = append(n.log, &sntcLog{
		condTruth: conds,
		form:      f,
		value:     value,
	})
}

func (n *langNode) filterConcepts(parentC concept, concepts map[int]concept,
	conds map[langCond]*bool, ctx *sntcCtx) map[int]concept {
	result := map[int]concept{}
	for _, c := range concepts {
		satisfied := true
		for cond, truth := range conds {
			if cond.satisfied(c, parentC, ctx) != truth {
				satisfied = false
			}
		}

		if satisfied {
			result[c.id()] = c
		}
	}

	return result
}

func (l *agentLanguage) newLangNode(class reflect.Type) *langNode {
	return &langNode{
		agent: l.agent,
		class: class,
		forms: []*langForm{},
		log:   []*sntcLog{},
		conds: map[langCond]bool{},
	}
}

type sntcLog struct {
	condTruth map[langCond]*bool
	form      *langForm
	value     float64
}
