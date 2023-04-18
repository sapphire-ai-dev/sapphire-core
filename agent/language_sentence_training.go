package agent

import (
	"encoding/json"
	"github.com/sapphire-ai-dev/sapphire-core/world"
	"io"
	"os"
	"reflect"
)

type trainSntc struct {
	body     string
	language *agentLanguage
	rootNode *trainSntcNode
	concepts map[int]concept
	speaker  object
	listener object
}

func (l *agentLanguage) newTrainSntc(body string, rootNode *trainSntcNode, concepts map[int]concept,
	speaker, listener object) *trainSntc {
	result := &trainSntc{
		body:     body,
		language: l,
		rootNode: rootNode,
		concepts: concepts,
		speaker:  speaker,
		listener: listener,
	}

	rootNode.setSentence(result)
	return result
}

type trainSntcNode struct {
	sentence  *trainSntc
	parent    *trainSntcNode
	conceptId int
	word      string
	isPronoun bool
	children  []*trainSntcNode
	infos     map[int]*trainSntcNodeInfo // partId -> info
}

func (l *agentLanguage) collectParts(root concept) (map[int]int, map[int]int) {
	forward, backward := map[int]int{}, map[int]int{} // part ID -> concept ID, concept ID -> part ID
	for partId := range l.agent.record.classes[reflect.TypeOf(root)].parts {
		if !isNil(root.part(partId)) {
			recConceptId := root.part(partId).id()
			forward[partId] = recConceptId
			backward[recConceptId] = partId
		}
	}

	return forward, backward
}

// internal as in not a leaf, in the semantic tree
func (t *trainSntcNode) buildInternalForm(root concept, ln *langNode, ctx *sntcCtx) ([]string, *langForm) {
	_, backward := t.sentence.language.collectParts(root)
	var sentence []string

	newForm := ln.newLangForm()
	for _, child := range t.children {
		phrase, childPart := child.buildForm(ctx)
		if lf, ok := childPart.(*langForm); ok {
			childPart = lf.newRecursiveLangPart(lf.node.class, backward[child.conceptId])
			t.sentence.language.registerWordPart(phrase[0], childPart)
		}
		newForm.parts = append(newForm.parts, childPart)
		sentence = append(sentence, phrase...)
	}

	newForm = ln.selectForm(newForm)
	t.sentence.language.registerWordPart(sentence[0], newForm)
	return sentence, newForm
}

func (t *trainSntcNode) buildLeafForm(ln *langNode) ([]string, *langForm) {
	newForm := ln.newLangForm()
	phrase, childPart := t.buildLeafPart(newForm, ln)
	newForm.parts = append(newForm.parts, childPart)
	newForm = ln.selectForm(newForm)
	t.sentence.language.registerWordPart(phrase[0], newForm)
	return phrase, newForm
}

func (t *trainSntcNode) buildLeafPart(workingForm *langForm, ln *langNode) ([]string, langPart) {
	if t.isPronoun {
		wlp := workingForm.newWordLangPart(ln.class, t.word)
		t.sentence.language.registerWordPart(t.word, wlp)
		return []string{t.word}, wlp
	} else {
		tenseId := -1
		c := t.sentence.concepts[t.conceptId]
		for i, w := range c.abs().tenses {
			if w == t.word {
				tenseId = i
			}
		}

		if tenseId == -1 {
			tenseId = len(c.abs().tenses)
			c.abs().tenses[tenseId] = t.word
		}

		clp := workingForm.newConceptLangPart(ln.class, tenseId)
		t.sentence.language.registerWordPart(t.word, clp)
		t.sentence.language.registerWordConcept(t.word, c, tenseId)
		return []string{t.word}, clp
	}
}

func (t *trainSntcNode) buildForm(ctx *sntcCtx) ([]string, langPart) {
	root := t.sentence.concepts[t.conceptId]
	ln := t.sentence.language.findLangNode(reflect.TypeOf(root))
	var sentence []string
	var newForm *langForm

	if len(t.children) != 0 {
		sentence, newForm = t.buildInternalForm(root, ln, ctx)
	} else {
		sentence, newForm = t.buildLeafForm(ln)
	}

	var parentConcept concept
	if t.parent != nil {
		parentConcept = t.sentence.concepts[t.parent.conceptId]
	}

	condTruth := ln.conditionTruth(root, parentConcept, ctx, map[langCond]bool{})
	infoExplicit := t.buildInfo(ln, root, ctx)
	ln.addLog(condTruth, newForm, formValueTrainSntc, infoExplicit)
	return sentence, newForm
}

func (t *trainSntcNode) buildInfo(ln *langNode, root concept, ctx *sntcCtx) map[int]bool {
	infoExplicit := map[int]bool{}
	for partId, i := range t.infos {
		li := ln.newLangInfo(partId)
		infoExplicit[partId] = i.explicit
		if i.explicit {
			li.record(implicitExplicit)
		} else {
			backupId, _ := ln.agent.language.fieldBackupId(root, i.implicitBackup, ctx)
			li.record(backupId)
		}
	}

	return infoExplicit
}

func (t *trainSntcNode) build() {
	ctx := t.sentence.language.newSntcCtx(t.sentence.speaker, t.sentence.listener)
	t.buildForm(ctx)
}

func (l *agentLanguage) newTrainSntcNode(conceptId int, word string, isPronoun bool) *trainSntcNode {
	return &trainSntcNode{
		conceptId: conceptId,
		word:      word,
		isPronoun: isPronoun,
		infos:     map[int]*trainSntcNodeInfo{},
	}
}

func (t *trainSntcNode) setSentence(s *trainSntc) {
	t.sentence = s
	for _, child := range t.children {
		child.setSentence(s)
	}
}

func (t *trainSntcNode) setChildren(children ...*trainSntcNode) {
	t.children = children
	for _, child := range children {
		child.parent = t
	}
}

type trainSntcNodeInfo struct {
	partId         int
	explicit       bool
	implicitBackup concept
}

func (t *trainSntcNode) newInfo(partId int, explicit bool, implicitBackup concept) {
	t.infos[partId] = &trainSntcNodeInfo{
		partId:         partId,
		explicit:       explicit,
		implicitBackup: implicitBackup,
	}
}

type trainSntcParser struct {
	l *agentLanguage
}

func (p *trainSntcParser) parse(file string) (*trainSntcData, []*trainSntc, map[int]concept) {
	var data *trainSntcData
	jsonFile, err := os.Open(file)
	printErr(err)

	byteValue, err := io.ReadAll(jsonFile)
	printErr(err)
	printErr(json.Unmarshal(byteValue, &data))

	if data == nil {
		return nil, nil, nil
	}

	printErr(jsonFile.Close())
	data.l = p.l
	sentences, concepts := data.parse()
	return data, sentences, concepts
}

func (l *agentLanguage) newTrainSntcParser() {
	l.trainParser = &trainSntcParser{l: l}
}

type trainSntcData struct {
	l                    *agentLanguage
	Concepts             []map[string]any     `json:"concepts"`
	Sentences            []*TrainSntcRootData `json:"sentences"`
	namedConcepts        map[string]concept
	actionInterfaces     map[int]*world.ActionInterface
	testActionInterfaces map[int]*TestActionInterface
}

func (d *trainSntcData) parse() ([]*trainSntc, map[int]concept) {
	d.namedConcepts = map[string]concept{}
	for _, entry := range d.Concepts {
		d.parseSingleConcept(entry)
	}

	concepts := map[int]concept{}
	for _, c := range d.namedConcepts {
		concepts[c.id()] = c
	}

	var sentences []*trainSntc
	for _, s := range d.Sentences {
		sentences = append(sentences, d.parseSingleSntc(concepts, s))
	}

	return sentences, concepts
}

func (d *trainSntcData) parseConceptArgs(data map[string]any) map[int]any {
	result := map[int]any{}
	if ctx, ctxOk := mapConcept[*contextObject](d, data, "ctx"); ctxOk {
		result[partIdConceptContext] = ctx
	}
	if temporal, temporalOk := mapConcept[temporalObject](d, data, "time"); temporalOk {
		result[partIdConceptTime] = temporal
	} else if temporalStr, temporalStrOk := mapVal[string](data, "time"); temporalStrOk {
		if temporalStr == "now" {
			result[partIdConceptTime] = d.l.agent.time.now
		}
	}

	return result
}

func (d *trainSntcData) parseSingleConcept(data map[string]any) {
	name, nameOk := mapVal[string](data, "name")
	class, classOk := mapVal[string](data, "class")
	if !nameOk || !classOk {
		return
	}

	args := d.parseConceptArgs(data)
	if parser, parserSeen := d.l.conceptParsers[class]; parserSeen {
		if c := parser(d, data, args); c != nil {
			d.namedConcepts[name] = c
		}
	}
}

func (d *trainSntcData) parseSingleSntc(concepts map[int]concept, r *TrainSntcRootData) *trainSntc {
	root := d.parseSingleNode(r.Root)
	var speaker, listener object
	if raw, seen := d.namedConcepts[r.Speaker]; seen {
		if t, ok := raw.(object); ok {
			speaker = t
		}
	}
	if raw, seen := d.namedConcepts[r.Listener]; seen {
		if t, ok := raw.(object); ok {
			listener = t
		}
	}

	sntc := d.l.newTrainSntc(r.Body, root, concepts, speaker, listener)
	return sntc
}

func (d *trainSntcData) parseSingleNode(n *TrainSntcNodeData) *trainSntcNode {
	rootC := d.namedConcepts[n.Concept]
	node := d.l.newTrainSntcNode(rootC.id(), n.Word, n.IsPronoun)
	_, backward := d.l.collectParts(rootC)
	var children []*trainSntcNode

	for _, childData := range n.Children {
		childNode := d.parseSingleNode(childData)
		children = append(children, childNode)
		childPartId, seen := backward[childNode.conceptId]
		if !seen {
			panic("child concept unrecognized by part record")
		}

		node.newInfo(childPartId, true, nil)
	}

	for _, implicitFieldName := range n.Implicit {
		implicitPartId, seen := d.l.fieldPartId(reflect.TypeOf(rootC), implicitFieldName)
		if !seen {
			panic("implicit field name unrecognized by parser")
		}

		implicitPart := rootC.part(implicitPartId)
		if isNil(implicitPart) {
			panic("implicit part does not exist in root")
		}

		node.newInfo(implicitPartId, false, implicitPart)
	}

	node.setChildren(children...)
	return node
}

func (d *trainSntcData) newActionInterface(i int) *world.ActionInterface {
	if d.actionInterfaces == nil {
		d.testActionInterfaces = map[int]*TestActionInterface{}
		d.actionInterfaces = map[int]*world.ActionInterface{}
	}

	if _, seen := d.actionInterfaces[i]; !seen {
		d.testActionInterfaces[i] = newTestActionInterface()
		d.actionInterfaces[i] = d.testActionInterfaces[i].instantiate()
	}

	return d.actionInterfaces[i]
}

func mapInt(m map[string]any, key string) (int, bool) {
	if raw, seen := m[key]; seen {
		if t, ok := raw.(float64); ok {
			return int(t), true
		}
	}

	return 0, false
}

func mapVal[T any](m map[string]any, key string) (T, bool) {
	if raw, seen := m[key]; seen {
		if t, ok := raw.(T); ok {
			return t, true
		}
	}

	var t T
	return t, false
}

func mapListVal[T any](m map[string]any, key string) ([]T, bool) {
	var result []T
	raw, seen := m[key]
	if !seen {
		return result, false
	}

	rawList, ok := raw.([]any)
	if !ok {
		return result, false
	}

	for _, elem := range rawList {
		if t, isT := elem.(T); isT {
			result = append(result, t)
		} else {
			return []T{}, false
		}
	}

	return result, true
}

func mapConcept[T concept](d *trainSntcData, m map[string]any, key string) (T, bool) {
	var result T
	name, nameOk := mapVal[string](m, key)
	if !nameOk {
		return result, false
	}

	if raw, seen := d.namedConcepts[name]; seen {
		if t, ok := raw.(T); ok {
			return t, true
		}
	}

	return result, false
}

func mapListConcept[T concept](d *trainSntcData, m map[string]any, key string) (map[int]T, bool) {
	result := map[int]T{}
	nameList, nameListOk := mapListVal[string](m, key)
	if !nameListOk {
		return result, false
	}

	for _, name := range nameList {
		elem, seen := d.namedConcepts[name]
		if !seen {
			return map[int]T{}, false
		}

		t, ok := elem.(T)
		if !ok {
			return map[int]T{}, false
		}

		result[t.id()] = t
	}

	return result, true
}

type TrainSntcRootData struct {
	Body     string             `json:"body"`
	Speaker  string             `json:"speaker"`
	Listener string             `json:"listener"`
	Root     *TrainSntcNodeData `json:"root"`
}

type TrainSntcNodeData struct {
	Word      string               `json:"word"`
	Concept   string               `json:"concept"`
	IsPronoun bool                 `json:"isPronoun"`
	Children  []*TrainSntcNodeData `json:"children"`
	Implicit  []string             `json:"implicit"`
}
