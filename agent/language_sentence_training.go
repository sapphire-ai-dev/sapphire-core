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
}

func (l *agentLanguage) newTrainSntc(body string, rootNode *trainSntcNode, concepts map[int]concept) *trainSntc {
	result := &trainSntc{
		body:     body,
		language: l,
		rootNode: rootNode,
		concepts: concepts,
	}

	rootNode.setSentence(result)
	return result
}

type trainSntcNode struct {
	sentence  *trainSntc
	conceptId int
	word      string
	isPronoun bool
	children  []*trainSntcNode
}

func (l *agentLanguage) collectParts(root concept) (map[int]int, map[int]int) {
	forward, backward := map[int]int{}, map[int]int{} // part ID -> concept ID, concept ID -> part ID
	for partId := range l.agent.record.classes[reflect.TypeOf(root)].parts {
		recConceptId := root.part(partId).id()
		forward[partId] = recConceptId
		backward[recConceptId] = partId
	}

	return forward, backward
}

func (t *trainSntcNode) build(workingForm *langForm) ([]string, langPart) {
	root := t.sentence.concepts[t.conceptId]
	vn := t.sentence.language.findLangNode(reflect.TypeOf(root))

	if len(t.children) != 0 {
		_, backward := t.sentence.language.collectParts(root)
		var sentence []string

		newForm := vn.newLangForm()
		for _, child := range t.children {
			phrase, childPart := child.build(newForm)
			if lf, ok := childPart.(*langForm); ok {
				childPart = lf.newRecursiveLangPart(lf.node.class, backward[child.conceptId])
				t.sentence.language.registerWordPart(phrase[0], childPart)
			}
			newForm.parts = append(newForm.parts, childPart)
			sentence = append(sentence, phrase...)
		}

		vn.selectForm(newForm)
		t.sentence.language.registerWordPart(sentence[0], newForm)
		return sentence, newForm
	}

	if t.isPronoun {
		wvp := workingForm.newWordLangPart(vn.class, t.word)
		t.sentence.language.registerWordPart(t.word, wvp)
		return []string{t.word}, wvp
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

		cvp := workingForm.newConceptLangPart(vn.class, tenseId)
		t.sentence.language.registerWordPart(t.word, cvp)
		t.sentence.language.registerWordConcept(t.word, c, tenseId)
		return []string{t.word}, cvp
	}
}

func (l *agentLanguage) newTrainSntcNode(conceptId int, word string, isPronoun bool) *trainSntcNode {
	return &trainSntcNode{
		conceptId: conceptId,
		word:      word,
		isPronoun: isPronoun,
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
}

type trainSntcParser struct {
	l *agentLanguage
}

func (p *trainSntcParser) parse(file string) ([]*trainSntc, map[int]concept) {
	var data *trainSntcData
	jsonFile, err := os.Open(file)
	printErr(err)

	byteValue, err := io.ReadAll(jsonFile)
	printErr(err)
	printErr(json.Unmarshal(byteValue, &data))

	if data == nil {
		return nil, nil
	}

	printErr(jsonFile.Close())
	data.l = p.l
	return data.parse()
}

func (l *agentLanguage) newTrainSntcParser() {
	l.trainParser = &trainSntcParser{l: l}
}

type trainSntcData struct {
	l                *agentLanguage
	Concepts         []map[string]any     `json:"concepts"`
	Sentences        []*TrainSntcRootData `json:"sentences"`
	namedConcepts    map[string]concept
	actionInterfaces map[int]*world.ActionInterface
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

func (d *trainSntcData) parseSingleConcept(data map[string]any) {
	name, nameOk := mapVal[string](data, "name")
	class, classOk := mapVal[string](data, "class")
	if !nameOk || !classOk {
		return
	}

	if parser, parserSeen := d.l.conceptParsers[class]; parserSeen {
		if c := parser(d, data); c != nil {
			d.namedConcepts[name] = c
		}
	}
}

// todo add body, speaker and listener
func (d *trainSntcData) parseSingleSntc(concepts map[int]concept, r *TrainSntcRootData) *trainSntc {
	root := d.parseSingleNode(r.Root)
	sntc := d.l.newTrainSntc(r.Body, root, concepts)
	return sntc
}

func (d *trainSntcData) parseSingleNode(n *TrainSntcNodeData) *trainSntcNode {
	node := d.l.newTrainSntcNode(d.namedConcepts[n.Concept].id(), n.Word, n.IsPronoun)
	var children []*trainSntcNode
	for _, childData := range n.Children {
		children = append(children, d.parseSingleNode(childData))
	}

	node.setChildren(children...)
	return node
}

func (d *trainSntcData) newActionInterface(i int) *world.ActionInterface {
	if d.actionInterfaces == nil {
		d.actionInterfaces = map[int]*world.ActionInterface{}
	}

	if _, seen := d.actionInterfaces[i]; !seen {
		d.actionInterfaces[i] = newTestActionInterface().instantiate()
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
}
