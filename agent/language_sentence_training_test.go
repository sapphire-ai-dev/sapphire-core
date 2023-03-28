package agent

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTrainingSentence(t *testing.T) {
	agent := newEmptyWorldAgent()
	amt1 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "blue"), nil)
	amt2 := agent.newAspectModifierType(agent.aspect.find(world.InfoLabelObservable, "color", "yellow"), nil)
	amt1.addSource(conceptSourceObservation)
	amt2.addSource(conceptSourceObservation)
	sot1 := agent.newSimpleObjectType(conceptSourceObservation, map[int]modifierType{amt1.id(): amt1}, nil)
	sot2 := agent.newSimpleObjectType(conceptSourceObservation, map[int]modifierType{amt2.id(): amt2}, nil)
	so1 := agent.newSimpleObject(123, nil)
	so2 := agent.newSimpleObject(234, nil)
	so1.addType(sot1)
	so2.addType(sot2)
	aat := agent.newAtomicActionType(newTestActionInterface().instantiate(), nil)
	aa := agent.newAtomicAction(aat, so1, nil)
	aa.setReceiver(so2)
	one := agent.symbolic.numerics.number1

	tsnAa := agent.language.newTrainSntcNode(aa.id(), "", false)
	tsnSo1 := agent.language.newTrainSntcNode(so1.id(), "", false)
	tsnSo1C := agent.language.newTrainSntcNode(so1.id(), "bob", false)
	tsnAat := agent.language.newTrainSntcNode(aat.id(), "", false)
	tsnAatC := agent.language.newTrainSntcNode(aat.id(), "ate", false)
	tsnSo2 := agent.language.newTrainSntcNode(so2.id(), "", false)
	tsnN := agent.language.newTrainSntcNode(one.id(), "", false)
	tsnNC := agent.language.newTrainSntcNode(one.id(), "a", false)
	tsnSot2 := agent.language.newTrainSntcNode(sot2.id(), "", false)
	tsnSot2C := agent.language.newTrainSntcNode(sot2.id(), "banana", false)

	tsnAa.setChildren(tsnSo1, tsnAat, tsnSo2)
	tsnSo1.setChildren(tsnSo1C)
	tsnAat.setChildren(tsnAatC)
	tsnSo2.setChildren(tsnN, tsnSot2)
	tsnN.setChildren(tsnNC)
	tsnSot2.setChildren(tsnSot2C)

	concepts := map[int]concept{}
	concepts[one.id()] = one
	concepts[aa.id()] = aa
	concepts[aat.id()] = aat
	concepts[so1.id()] = so1
	concepts[so2.id()] = so2
	concepts[sot2.id()] = sot2

	body := "bob ate a banana"
	ts := agent.language.newTrainSntc(body, tsnAa, concepts)
	ts.rootNode.build(nil)

	sn := agent.language.fit(strings.Split(body, " "), agent.language.newSntcCtx(nil, nil))
	assert.Equal(t, strings.Join(sn.str(), " "), body)
}

func TestParse(t *testing.T) {
	agent := newEmptyWorldAgent()
	sentences, _ := agent.language.trainParser.parse("data/data.json")
	for _, sentence := range sentences {
		sentence.rootNode.build(nil)
	}

	for _, sentence := range sentences {
		ctx := agent.language.newSntcCtx(nil, nil)
		sn := agent.language.fit(strings.Split(sentence.body, " "), ctx)
		assert.Equal(t, strings.Join(sn.str(), " "), sentence.body)
	}
}
