package agent

import (
    "github.com/sapphire-ai-dev/sapphire-core/world"
)

func (l *agentLanguage) initConceptParsers() {
    l.conceptParsers["AspectModifierType"] = l.dataParserAspectModifierType
    l.conceptParsers["SimpleObjectType"] = l.dataParserSimpleObjectType
    l.conceptParsers["SimpleObject"] = l.dataParserSimpleObject
    l.conceptParsers["AtomicActionType"] = l.dataParserAtomicActionType
    l.conceptParsers["AtomicAction"] = l.dataParserAtomicAction
    l.conceptParsers["Number"] = l.dataParserNumber
    l.conceptParsers["IdentityRelationType"] = l.dataParserIdentityRelationType
    l.conceptParsers["IdentityRelation"] = l.dataParserIdentityRelation
    l.conceptParsers["AuxiliaryRelationType"] = l.dataParserAuxiliaryRelationType
    l.conceptParsers["AuxiliaryRelation"] = l.dataParserAuxiliaryRelation
}

func (l *agentLanguage) dataParserAspectModifierType(_ *trainSntcData, data map[string]any) concept {
    labels, ok := mapListVal[string](data, "labels")
    if !ok {
        return nil
    }

    labels = append([]string{world.InfoLabelObservable}, labels...)
    result := l.agent.newAspectModifierType(l.agent.aspect.find(labels...), nil)
    result.addSource(conceptSourceObservation)
    return result
}

func (l *agentLanguage) dataParserSimpleObjectType(d *trainSntcData, data map[string]any) concept {
    modifierTypes, ok := mapListConcept[modifierType](d, data, "modifierTypes")
    if !ok {
        return nil
    }

    return l.agent.newSimpleObjectType(conceptSourceObservation, modifierTypes, nil)
}

func (l *agentLanguage) dataParserAtomicActionType(d *trainSntcData, data map[string]any) concept {
    actionInterfaceId, ok := mapInt(data, "interfaceId")
    if !ok {
        return nil
    }

    return l.agent.newAtomicActionType(d.newActionInterface(actionInterfaceId), nil)
}

func (l *agentLanguage) dataParserSimpleObject(d *trainSntcData, data map[string]any) concept {
    worldId, worldIdOk := mapInt(data, "worldId")
    objectTypes, objectTypesOk := mapListConcept[objectType](d, data, "types")
    if !worldIdOk || !objectTypesOk {
        return nil
    }

    result := l.agent.newSimpleObject(worldId, nil)
    for _, t := range objectTypes {
        result.addType(t)
    }

    return result
}

func (l *agentLanguage) dataParserAtomicAction(d *trainSntcData, data map[string]any) concept {
    aat, aatOk := mapConcept[*atomicActionType](d, data, "type")
    if !aatOk {
        return nil
    }

    performer, performerOk := mapConcept[object](d, data, "performer")
    if !performerOk {
        performer = l.agent.self
    }

    result := l.agent.newAtomicAction(aat, performer, nil)
    receiver, receiverOk := mapConcept[object](d, data, "receiver")
    if receiverOk {
        result.setReceiver(receiver)
    }

    return result
}

func (l *agentLanguage) dataParserNumber(_ *trainSntcData, data map[string]any) concept {
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

func (l *agentLanguage) dataParserIdentityRelationType(_ *trainSntcData, _ map[string]any) concept {
    return l.agent.newIdentityRelationType(nil)
}

func (l *agentLanguage) dataParserIdentityRelation(d *trainSntcData, data map[string]any) concept {
    irt, irtOk := mapConcept[*identityRelationType](d, data, "type")
    lTarget, lTargetOk := mapConcept[object](d, data, "lTarget")
    rTarget, rTargetOk := mapConcept[objectType](d, data, "rTarget")
    if !irtOk || !lTargetOk || !rTargetOk {
        return nil
    }

    result := l.agent.newIdentityRelation(irt, lTarget, rTarget, nil)
    return result
}

func (l *agentLanguage) dataParserAuxiliaryRelationType(_ *trainSntcData, data map[string]any) concept {
    auxiliaryTypeName, ok := mapVal[string](data, "type")
    if !ok {
        return nil
    }

    auxiliaryTypeId, seen := auxiliaryTypeIdNames[auxiliaryTypeName]
    if !seen {
        return nil
    }

    return l.agent.newAuxiliaryRelationType(auxiliaryTypeId, nil)
}

func (l *agentLanguage) dataParserAuxiliaryRelation(d *trainSntcData, data map[string]any) concept {
    art, artOk := mapConcept[*auxiliaryRelationType](d, data, "type")
    lTarget, lTargetOk := mapConcept[object](d, data, "lTarget")
    rTarget, rTargetOk := mapConcept[action](d, data, "rTarget")
    if !artOk || !lTargetOk || !rTargetOk {
        return nil
    }

    result := l.agent.newAuxiliaryRelation(art, lTarget, rTarget, nil)
    return result
}
