package agent

import (
	"fmt"
	"github.com/sapphire-ai-dev/sapphire-core/world"
)

type agentPerception struct {
	agent          *Agent
	visibleObjects []object
}

func (p *agentPerception) cycle() {
	p.look()
	p.listen()
}

func (p *agentPerception) look() {
	p.visibleObjects = []object{}
	for _, image := range world.Look(p.agent.self.worldId) {
		p.processImage(image)
	}
}

func (p *agentPerception) processImage(img *world.Image) {
	modifTypes, instantiateArgs := p.parseImage(img)
	obj := p.identifyObjectInst(img)
	objType := p.identifyObjectType(modifTypes)
	knownTypes := obj.types()
	if _, seen := knownTypes[objType.id()]; !seen {
		obj.addType(objType)
	}
	for modifTypeId, modifType := range modifTypes {
		modif := modifType.instantiate(obj, conceptSourceObservation, nil, instantiateArgs[modifTypeId]...)
		p.agent.mind.add(modif)
	}

	p.visibleObjects = append(p.visibleObjects, obj)
}

func (p *agentPerception) parseImage(img *world.Image) (map[int]modifierType, map[int][]any) {
	modifTypes := map[int]modifierType{}
	instantiateArgs := map[int][]any{}
	for _, permanentInfo := range img.Permanent {
		permanentLabels := append([]string{aspectObjectInfoDebugName}, permanentInfo.Labels...)
		modifType := p.agent.newAspectModifierType(p.agent.aspect.find(permanentLabels...), nil)
		modifTypes[modifType.id()] = modifType
		instantiateArgs[modifType.id()] = []any{permanentInfo.Value}
	}

	for _, transientInfo := range img.Transient {
		transientLabels := append([]string{aspectObjectInfoDebugName}, transientInfo.Labels...)
		modifType := p.agent.newAspectModifierType(p.agent.aspect.find(transientLabels...), nil)
		modifTypes[modifType.id()] = modifType
		instantiateArgs[modifType.id()] = []any{transientInfo.Value}
	}

	return modifTypes, instantiateArgs
}

func (p *agentPerception) identifyObjectInst(img *world.Image) object {
	return p.agent.newSimpleObject(map[int]any{partIdObjectWorldId: img.Id})
}

func (p *agentPerception) identifyObjectType(modifTypes map[int]modifierType) objectType {
	return p.agent.newSimpleObjectType(conceptSourceObservation, modifTypes, nil)
}

func (p *agentPerception) listen() {
	sentences := world.Listen(p.agent.self.worldId)
	for _, sentence := range sentences {
		sp := p.agent.language.listen(sentence)
		if sp != nil {
			fmt.Println("agent heard", sp.str())
		} else {
			fmt.Println("agent heard and did not understand:", sentence.Body)
		}
	}
}

func (a *Agent) newAgentPerception() {
	a.perception = &agentPerception{
		agent:          a,
		visibleObjects: []object{},
	}
}
