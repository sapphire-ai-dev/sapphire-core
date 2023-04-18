package agent

import (
	"github.com/sapphire-ai-dev/sapphire-core/world"
	"reflect"
)

type backupConceptRecord struct{}

const (
	implicitExplicit = iota // use this first one for explicit mentions
	implicitObjectSelf
	implicitTemporalNow
	implicitContext0 // temporary todo remove
	implicitActionStateChange
)

func (l *agentLanguage) fieldBackupId(root, backupInst concept, ctx *sntcCtx) (int, bool) {
	if backupInst == l.agent.time.now {
		return implicitTemporalNow, true
	}

	if reflect.TypeOf(backupInst) == toReflect[*actionStateChange]() {
		return l.fieldBackupIdActionStateChange(root, backupInst, ctx)
	}

	if reflect.TypeOf(backupInst).Implements(toReflect[object]()) {
		if backupInst == l.agent.self {
			return implicitObjectSelf, true
		}
	}

	if reflect.TypeOf(backupInst) == toReflect[*contextObject]() {
		// todo remove workaround
		return implicitContext0, true
	}

	return 0, false
}

func (l *agentLanguage) fieldBackupIdActionStateChange(root, backupInst concept, ctx *sntcCtx) (int, bool) {
	return implicitActionStateChange, true
}

func (l *agentLanguage) generateImplicitConcept(implicitId int) concept {
	if implicitId == implicitObjectSelf {
		return l.agent.self
	}

	if implicitId == implicitTemporalNow {
		return l.agent.time.now
	}

	if implicitId == implicitContext0 { // todo remove workaround
		ccat := l.agent.newCreateContextActionType()
		cca := l.agent.newCreateContextAction(ccat, l.agent.self, 0)
		cot := l.agent.newContextObjectType(conceptSourceObservation)
		co := l.agent.newContextObject(cca)
		co.addType(cot)
		return co
	}

	if implicitId == implicitActionStateChange { // todo remove workaround
		aat := l.agent.newAtomicActionType(&world.ActionInterface{}, nil)
		aa := l.agent.newAtomicAction(map[int]any{partIdActionT: aat})
		asct := l.agent.newActionStateChangeType(aat, nil)
		asct.addValue(10.0)
		asc := l.agent.newActionStateChange(asct, aa, nil)
		return asc
	}

	return nil
}

func (l *agentLanguage) newBackupConceptRecord() {
	l.backupRecord = &backupConceptRecord{}
}
