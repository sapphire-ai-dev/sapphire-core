package agent

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
)

func (l *agentLanguage) initConceptParsers() {
	l.conceptParsers["AspectModifierType"] = l.parserAspectModifierType
	l.conceptParsers["SimpleObjectType"] = l.parserSimpleObjectType
	l.conceptParsers["SimpleObject"] = l.parserSimpleObject
	l.conceptParsers["SelfObject"] = l.parserSelfObject
	l.conceptParsers["AtomicActionType"] = l.parserAtomicActionType
	l.conceptParsers["AtomicAction"] = l.parserAtomicAction
	l.conceptParsers["Number"] = l.parserNumber
	l.conceptParsers["IdentityRelationType"] = l.parserIdentityRelationType
	l.conceptParsers["IdentityRelation"] = l.parserIdentityRelation
	l.conceptParsers["AuxiliaryRelationType"] = l.parserAuxiliaryRelationType
	l.conceptParsers["AuxiliaryRelation"] = l.parserAuxiliaryRelation
	l.conceptParsers["ActionStateChangeType"] = l.parserActionStateChangeType
	l.conceptParsers["ActionStateChange"] = l.parserActionStateChange
	l.conceptParsers["ContextObjectType"] = l.parserContextObjectType
	l.conceptParsers["ContextObject"] = l.parserContextObject
	l.conceptParsers["CreateContextActionType"] = l.parserCreateContextActionType
	l.conceptParsers["CreateContextAction"] = l.parserCreateContextAction
}

func (l *agentLanguage) parserAspectModifierType(_ *trainSntcData, data map[string]any, args map[int]any) concept {
	labels, ok := mapListVal[string](data, "labels")
	if !ok {
		return nil
	}

	labels = append([]string{world.InfoLabelObservable}, labels...)
	result := l.agent.newAspectModifierType(l.agent.aspect.find(labels...), args)
	result.addSource(conceptSourceObservation)
	return result
}

func (l *agentLanguage) parserSimpleObjectType(d *trainSntcData, data map[string]any, args map[int]any) concept {
	modifierTypes, ok := mapListConcept[modifierType](d, data, "modifierTypes")
	if !ok {
		return nil
	}

	return l.agent.newSimpleObjectType(conceptSourceObservation, modifierTypes, args)
}

func (l *agentLanguage) parserAtomicActionType(d *trainSntcData, data map[string]any, args map[int]any) concept {
	actionInterfaceId, ok := mapInt(data, "interfaceId")
	if !ok {
		return nil
	}

	result := l.agent.newAtomicActionType(d.newActionInterface(actionInterfaceId), args)
	l.agent.mind.add(result)
	return result
}

func (l *agentLanguage) parserCreateContextActionType(_ *trainSntcData, _ map[string]any, _ map[int]any) concept {
	return l.agent.newCreateContextActionType()
}

func (l *agentLanguage) parserCreateContextAction(d *trainSntcData, data map[string]any, _ map[int]any) concept {
	ccat, ccatOk := mapConcept[*createContextActionType](d, data, "type")
	contextId, contextIdOk := mapInt(data, "contextId")
	if !ccatOk || !contextIdOk {
		return nil
	}

	performer, performerOk := mapConcept[object](d, data, "performer")
	if !performerOk {
		performer = l.agent.self
	}

	return l.agent.newCreateContextAction(ccat, performer, contextId)
}

func (l *agentLanguage) parserContextObjectType(_ *trainSntcData, _ map[string]any, _ map[int]any) concept {
	return l.agent.newContextObjectType(conceptSourceObservation)
}

func (l *agentLanguage) parserContextObject(d *trainSntcData, data map[string]any, _ map[int]any) concept {
	cot, cotOk := mapConcept[*contextObjectType](d, data, "type")
	ca, caOk := mapConcept[*createContextAction](d, data, "creation")
	if !cotOk || !caOk {
		return nil
	}

	result := l.agent.newContextObject(ca)
	result.addType(cot)
	return result
}

const (
	dataParserSelfObjectAttachSelf = "self"
)

func (l *agentLanguage) parserSelfObject(d *trainSntcData, data map[string]any, args map[int]any) concept {
	attach, attachOK := mapVal[string](data, "attach")
	worldId, worldIdOk := mapInt(data, "worldId")
	objectTypes, objectTypesOk := mapListConcept[objectType](d, data, "types")
	var result *selfObject
	if attachOK && attach == dataParserSelfObjectAttachSelf {
		result = l.agent.self
	} else if worldIdOk {
		result = l.agent.newSelfObject(worldId, args)
	} else {
		return nil
	}

	if objectTypesOk {
		for _, t := range objectTypes {
			result.addType(t)
		}
	}

	return result
}

func (l *agentLanguage) parserSimpleObject(d *trainSntcData, data map[string]any, args map[int]any) concept {
	worldId, worldIdOk := mapInt(data, "worldId")
	objectTypes, objectTypesOk := mapListConcept[objectType](d, data, "types")
	if !worldIdOk || !objectTypesOk {
		return nil
	}

	result := l.agent.newSimpleObject(worldId, args)
	for _, t := range objectTypes {
		result.addType(t)
	}

	return result
}

func (l *agentLanguage) parserAtomicAction(d *trainSntcData, data map[string]any, args map[int]any) concept {
	aat, aatOk := mapConcept[*atomicActionType](d, data, "type")
	if !aatOk {
		return nil
	}

	performer, performerOk := mapConcept[object](d, data, "performer")
	if !performerOk {
		performer = l.agent.self
	}

	result := l.agent.newAtomicAction(aat, performer, args)
	receiver, receiverOk := mapConcept[object](d, data, "receiver")
	if receiverOk {
		result.setReceiver(receiver)
	}

	return result
}

func (l *agentLanguage) parserNumber(_ *trainSntcData, data map[string]any, _ map[int]any) concept {
	value, ok := mapInt(data, "value")
	if !ok {
		return nil
	}

	if value == 0 {
		return l.agent.symbolic.numerics.number0
	} else if value == 1 {
		return l.agent.symbolic.numerics.number1
	}

	return nil
}

func (l *agentLanguage) parserIdentityRelationType(_ *trainSntcData, _ map[string]any, args map[int]any) concept {
	return l.agent.newIdentityRelationType(args)
}

func (l *agentLanguage) parserIdentityRelation(d *trainSntcData, data map[string]any, args map[int]any) concept {
	irt, irtOk := mapConcept[*identityRelationType](d, data, "type")
	lTarget, lTargetOk := mapConcept[object](d, data, "lTarget")
	rTarget, rTargetOk := mapConcept[objectType](d, data, "rTarget")
	if !irtOk || !lTargetOk || !rTargetOk {
		return nil
	}

	result := l.agent.newIdentityRelation(irt, lTarget, rTarget, args)
	return result
}

func (l *agentLanguage) parserAuxiliaryRelationType(d *trainSntcData, data map[string]any, args map[int]any) concept {
	auxiliaryTypeName, atNameOk := mapVal[string](data, "type")
	lType, lTypeOk := mapConcept[objectType](d, data, "lType")
	rType, rTypeOk := mapConcept[performableActionType](d, data, "rType")
	negative, negativeOk := mapVal[bool](data, "negative")
	if !atNameOk || !lTypeOk || !rTypeOk || !negativeOk {
		return nil
	}

	auxiliaryTypeId, atIdSeen := auxiliaryTypeIdNames[auxiliaryTypeName]
	if !atIdSeen {
		return nil
	}

	return l.agent.newAuxiliaryRelationType(auxiliaryTypeId, negative, lType, rType, args)
}

func (l *agentLanguage) parserAuxiliaryRelation(d *trainSntcData, data map[string]any, args map[int]any) concept {
	art, artOk := mapConcept[*auxiliaryRelationType](d, data, "type")
	lTarget, lTargetOk := mapConcept[object](d, data, "lTarget")
	rTarget, rTargetOk := mapConcept[performableAction](d, data, "rTarget")
	if !artOk || !lTargetOk || !rTargetOk {
		return nil
	}

	wantChange, wantChangeOk := mapConcept[*actionStateChange](d, data, "wantChange")
	if wantChangeOk {
		args[conceptArgRelationAuxiliaryWantChange] = wantChange
	}

	result := l.agent.newAuxiliaryRelation(art, lTarget, rTarget, args)
	result.interpret()
	return result
}

func (l *agentLanguage) parserActionStateChangeType(d *trainSntcData, data map[string]any, args map[int]any) concept {
	target, targetOk := mapConcept[performableActionType](d, data, "target")
	if !targetOk {
		return nil
	}

	return l.agent.newActionStateChangeType(target, args)
}

func (l *agentLanguage) parserActionStateChange(d *trainSntcData, data map[string]any, args map[int]any) concept {
	asct, asctOk := mapConcept[*actionStateChangeType](d, data, "type")
	target, targetOk := mapConcept[performableAction](d, data, "target")
	if !asctOk || !targetOk {
		return nil
	}

	value, valueOk := mapVal[float64](data, "value")
	if valueOk {
		asct.addValue(value)
	}

	return l.agent.newActionStateChange(asct, target, args)
}
