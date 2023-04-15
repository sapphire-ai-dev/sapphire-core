package agent

import (
	"reflect"
)

type assembleConceptRecord struct {
	l          *agentLanguage
	assemblers map[reflect.Type]func(args map[int]any) concept
}

func (r *assembleConceptRecord) assemble(class reflect.Type, args map[int]any) concept {
	return r.assemblers[class](args)
}

func (r *assembleConceptRecord) initAssemblers() {
	addAssembler[*atomicAction](r, r.l.agent.newAtomicAction)
	addAssembler[*virtualAction](r, r.l.agent.newVirtualAction)

	addAssembler[*simpleObject](r, r.l.agent.newSimpleObject)

	addAssembler[*auxiliaryRelation](r, r.l.agent.newAuxiliaryRelation)
	addAssembler[*virtualSolutionRelation](r, r.l.agent.newVirtualSolutionRelation)
}

func (l *agentLanguage) newAssembleConceptRecord() {
	l.assembleRecord = &assembleConceptRecord{
		l:          l,
		assemblers: map[reflect.Type]func(args map[int]any) concept{},
	}
	l.assembleRecord.initAssemblers()
}

func addAssembler[T concept](r *assembleConceptRecord, assembler func(args map[int]any) T) {
	r.assemblers[toReflect[T]()] = func(args map[int]any) concept {
		return assembler(args)
	}
}
