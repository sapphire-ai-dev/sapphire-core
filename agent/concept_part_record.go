package agent

import (
	"reflect"
)

// the idea of partId is used to enumerate part concepts of concepts
// i.e. consider the action A = "romeo meets juliet"
//   if we have performer = 1, receiver = 2
//   then we can have A.part(1) = romeo, A.part(2) = juliet
// this works fine, but only because the list ["performer", "receiver"] is finite and can be
//   exhaustively hardcoded, sometimes this is not the case
// i.e. consider a lemon L that is YELLOW and ROUND
//   other than "color" and "shape", there are infinitely more possible parts such as "taste",
//   "price", "whether bob likes it", etc. the list of partIds is impossible to exhaustively
//   predetermine, thus it must be dynamic

// partRecord is the class to record these dynamic partIds, there will be a permanent singleton
//	 instance stored in the agent

// part Ids are 32-bit signed integers, for arbitrary part Id P:
//   - P < 2^16: P is a hardcoded part Id (known at development time, i.e. performer of action)
//   - P > 2^16:
type partRecord struct {
	*abstractConcept
	agent               *Agent
	classes             map[reflect.Type]*conceptClass
	partTypeIds         map[int]*memReference // getting a part id -> search for instance with this type
	imaginaryGenerators map[reflect.Type]func(map[int]any) concept
}

func (r *partRecord) match(_ concept) bool {
	return false
}

func (r *partRecord) partImpls(class reflect.Type, partId int) map[reflect.Type]bool {
	cRecord, classSeen := r.classes[class]
	if !classSeen {
		return map[reflect.Type]bool{}
	}

	pRecord, partSeen := cRecord.parts[partId]
	if !partSeen {
		return map[reflect.Type]bool{}
	}

	return r.classImpls(pRecord.class)
}

func (r *partRecord) classImpls(class reflect.Type) map[reflect.Type]bool {
	result := map[reflect.Type]bool{}
	r.classImplsHelper(class, result)
	return result
}

func (r *partRecord) classImplsHelper(class reflect.Type, result map[reflect.Type]bool) {
	if result[class] {
		return
	}

	result[class] = true
	for childClass := range r.classes[class].children {
		r.classImplsHelper(childClass.class, result)
	}
}

func (a *Agent) newPartRecord() *partRecord {
	result := &partRecord{
		agent:       a,
		classes:     map[reflect.Type]*conceptClass{},
		partTypeIds: map[int]*memReference{},
	}

	result.initClasses()
	result.initImagineReflects()
	a.newAbstractConcept(result, nil, &result.abstractConcept)
	return result.memorize().(*partRecord)
}

type conceptClass struct {
	record   *partRecord
	class    reflect.Type
	parent   *conceptClass          // self implement parent
	children map[*conceptClass]bool // child implement self
	parts    map[int]*conceptClass  // partId -> part class
}

func (c *conceptClass) addChild(childClass *conceptClass) {
	if childClass.parent != nil && childClass.parent.children[childClass] {
		delete(childClass.parent.children, childClass)
	}

	c.children[childClass] = true
	childClass.parent = c
}

func (c *conceptClass) addPart(partId int, partClass *conceptClass) {
	c.parts[partId] = partClass
}

func (r *partRecord) newConceptClass(class reflect.Type) *conceptClass {
	_, seen := r.classes[class]
	if !seen {
		r.classes[class] = &conceptClass{
			record:   r,
			class:    class,
			children: map[*conceptClass]bool{},
			parts:    map[int]*conceptClass{},
		}
	}

	return r.classes[class]
}

const (
	partIdStart = iota
	partIdActionT
	partIdActionPerformer
	partIdActionReceiver
	partIdActionSimpleChild
	partIdActionSequentialFirst
	partIdActionSequentialNext
	partIdModifierT
	partIdModifierTarget
	partIdObjectT
	partIdObjectGroupSize
	partIdRelationT
	partIdRelationLTarget
	partIdRelationRTarget
)

func (r *partRecord) initClasses() {
	r.initClassesSingle(toReflect[concept](), nil, map[int]reflect.Type{})

	r.initClassesSingle(toReflect[actionType](), toReflect[concept](), map[int]reflect.Type{})
	r.initClassesSingle(toReflect[action](), toReflect[concept](), map[int]reflect.Type{
		partIdActionT:         toReflect[actionType](),
		partIdActionPerformer: toReflect[object](),
		partIdActionReceiver:  toReflect[object](),
	})
	r.initClassesSingle(toReflect[*atomicActionType](), toReflect[actionType](), map[int]reflect.Type{})
	r.initClassesSingle(toReflect[*atomicAction](), toReflect[action](), map[int]reflect.Type{
		partIdActionT: toReflect[*atomicActionType](),
	})
	r.initClassesSingle(toReflect[*simpleActionType](), toReflect[actionType](), map[int]reflect.Type{})
	r.initClassesSingle(toReflect[*simpleAction](), toReflect[action](), map[int]reflect.Type{
		partIdActionT:           toReflect[*simpleActionType](),
		partIdActionSimpleChild: toReflect[action](),
	})
	r.initClassesSingle(toReflect[*sequentialActionType](), toReflect[actionType](), map[int]reflect.Type{})
	r.initClassesSingle(toReflect[*sequentialAction](), toReflect[action](), map[int]reflect.Type{
		partIdActionT:               toReflect[*sequentialActionType](),
		partIdActionSequentialFirst: toReflect[action](),
		partIdActionSequentialNext:  toReflect[action](),
	})

	r.initClassesSingle(toReflect[modifierType](), toReflect[concept](), map[int]reflect.Type{})
	r.initClassesSingle(toReflect[modifier](), toReflect[concept](), map[int]reflect.Type{
		partIdModifierT:      toReflect[modifierType](),
		partIdModifierTarget: toReflect[concept](),
	})
	r.initClassesSingle(toReflect[*aspectModifierType](), toReflect[modifierType](), map[int]reflect.Type{})
	r.initClassesSingle(toReflect[*aspectModifier](), toReflect[modifierType](), map[int]reflect.Type{
		partIdModifierT: toReflect[*aspectModifierType](),
	})

	r.initClassesSingle(toReflect[objectType](), toReflect[concept](), map[int]reflect.Type{})
	r.initClassesSingle(toReflect[object](), toReflect[concept](), map[int]reflect.Type{
		partIdObjectT:         toReflect[objectType](),
		partIdObjectGroupSize: toReflect[*number](),
	})
	r.initClassesSingle(toReflect[*simpleObjectType](), toReflect[objectType](), map[int]reflect.Type{})
	r.initClassesSingle(toReflect[*simpleObject](), toReflect[object](), map[int]reflect.Type{})

	r.initClassesSingle(toReflect[*selfObject](), toReflect[object](), map[int]reflect.Type{})

	r.initClassesSingle(toReflect[relationType](), toReflect[concept](), map[int]reflect.Type{})
	r.initClassesSingle(toReflect[relation](), toReflect[concept](), map[int]reflect.Type{
		partIdRelationT:       toReflect[relationType](),
		partIdRelationLTarget: toReflect[concept](),
		partIdRelationRTarget: toReflect[concept](),
	})

	r.initClassesSingle(toReflect[*identityRelationType](), toReflect[relationType](), map[int]reflect.Type{})
	r.initClassesSingle(toReflect[*identityRelation](), toReflect[relation](), map[int]reflect.Type{
		partIdRelationT:       toReflect[*identityRelationType](),
		partIdRelationLTarget: toReflect[object](),
		partIdRelationRTarget: toReflect[objectType](),
	})
	r.initClassesSingle(toReflect[*auxiliaryRelationType](), toReflect[relationType](), map[int]reflect.Type{})
	r.initClassesSingle(toReflect[*auxiliaryRelation](), toReflect[relation](), map[int]reflect.Type{
		partIdRelationT:       toReflect[*auxiliaryRelationType](),
		partIdRelationLTarget: toReflect[object](),
		partIdRelationRTarget: toReflect[action](),
	})

	r.initClassesSingle(toReflect[*number](), toReflect[concept](), map[int]reflect.Type{})
}

func (r *partRecord) initClassesSingle(class, parent reflect.Type, parts map[int]reflect.Type) {
	c := r.newConceptClass(class)
	if parent != nil {
		p := r.newConceptClass(parent)
		p.addChild(c)

		// inherit parts from parents, can be overwritten later by specified parts
		for partId, part := range p.parts {
			c.parts[partId] = part
		}
	}

	for partId, part := range parts {
		c.parts[partId] = r.newConceptClass(part)
	}
}

func (r *partRecord) initImagineReflects() {
	r.imaginaryGenerators = map[reflect.Type]func(map[int]any) concept{}

	r.imaginaryGenerators[toReflect[object]()] = r.agent.newImaginaryObject
}

func (r *partRecord) generateImaginary(class reflect.Type, args map[int]any) concept {
	generator := r.imaginaryGenerators[class]
	classNode := r.classes[class]
	for generator == nil {
		classNode = classNode.parent
		generator = r.imaginaryGenerators[classNode.class]
	}
	return generator(args)
}

const (
	conceptArgContext = iota
	conceptArgTime

	conceptArgRelationAuxiliaryWantChange
)

func conceptArg[T any](m map[int]any, key int) (T, bool) {
	if m != nil {
		if raw, seen := m[key]; seen {
			if t, ok := raw.(T); ok {
				return t, true
			}
		}
	}

	var t T
	return t, false
}
