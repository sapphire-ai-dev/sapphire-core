package agent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTimePointObjectConstructor(t *testing.T) {
	agent := newEmptyWorldAgent()
	tpo := agent.newTimePointObject(&agent.time.clock, nil)
	assert.Equal(t, *tpo.clockTime, agent.time.clock)
	assert.Equal(t, tpo, agent.newTimePointObject(&agent.time.clock, nil))
	assert.NotEqual(t, tpo, agent.newTimePointObject(nil, nil))
	agent.time.clock++
	assert.NotEqual(t, *tpo.clockTime, agent.time.clock)
	assert.NotEqual(t, tpo, agent.newTimePointObject(&agent.time.clock, nil))
	assert.Equal(t, tpo, tpo.start())
	assert.Equal(t, tpo, tpo.end())
}

func TestTimePointObjectStartEnd(t *testing.T) {
	agent := newEmptyWorldAgent()
	tpo0 := agent.newTimePointObject(&agent.time.clock, nil)
	assert.Equal(t, tpo0, tpo0.start())
	assert.Equal(t, tpo0, tpo0.end())
	assert.Equal(t, tpo0, tpo0.join(tpo0))

	agent.time.clock++
	tpo1 := agent.newTimePointObject(&agent.time.clock, nil)
	tso := tpo0.join(tpo1)
	assert.Equal(t, tpo0, tso.start())
	assert.Equal(t, tpo1, tso.end())
	assert.Equal(t, tso, tpo1.join(tpo0))

	assert.Nil(t, tpo0.join(nil))
}

func assertComparativeRelationSeen(t *testing.T, rels map[int]relation, crt *comparativeRelationType,
	l, r concept) *comparativeRelation {
	var result *comparativeRelation
	for _, rel := range rels {
		if cr, ok := rel.(*comparativeRelation); ok {
			result = cr
			assert.Equal(t, crt, cr._type())
			assert.Equal(t, l, cr.lTarget())
			assert.Equal(t, r, cr.rTarget())
		}
	}
	assert.NotNil(t, result)
	return result
}

func assertIntervalRelationSeen(t *testing.T, rels map[int]relation, irt *intervalRelationType,
	l, r concept) *intervalRelation {
	var result *intervalRelation
	for _, rel := range rels {
		if cr, ok := rel.(*intervalRelation); ok {
			result = cr
			assert.Equal(t, irt, cr._type())
			assert.Equal(t, l, cr.lTarget())
			assert.Equal(t, r, cr.rTarget())
		}
	}
	assert.NotNil(t, result)
	return result
}

func TestTimePointObjectCompare(t *testing.T) {
	agent := newEmptyWorldAgent()
	tpo0 := agent.newTimePointObject(&agent.time.clock, nil)
	agent.time.clock++
	tpo1 := agent.newTimePointObject(&agent.time.clock, nil)
	assertComparativeRelationSeen(t, tpo0.compare(tpo1), agent.logic.comparatives.lt, tpo0, tpo1)
	assertComparativeRelationSeen(t, tpo1.compare(tpo0), agent.logic.comparatives.lt, tpo0, tpo1)
	assertComparativeRelationSeen(t, tpo0.compare(tpo0), agent.logic.comparatives.eq, tpo0, tpo0)
	assert.Empty(t, tpo0.compare(nil))
	assert.Empty(t, tpo0.compare(agent.newTimePointObject(nil, nil)))
}

func newTestContext(agent *Agent, contextId int) *contextObject {
	ccat := agent.newCreateContextActionType()
	cca := agent.newCreateContextAction(ccat, agent.self, contextId)
	cot := agent.newContextObjectType(conceptSourceObservation)
	co := agent.newContextObject(cca)
	co.addType(cot)
	return co
}

func TestTimePointObjectWithContext(t *testing.T) {
	agent := newEmptyWorldAgent()
	co0 := newTestContext(agent, 0)
	tpo0 := agent.newTimePointObject(&agent.time.clock, map[int]any{partIdConceptContext: co0})
	assert.Equal(t, co0, tpo0.ctx())
	assert.Equal(t, co0, tpo0.part(partIdConceptContext))
	assert.Equal(t, tpo0, agent.newTimePointObject(&agent.time.clock, map[int]any{partIdConceptContext: co0}))
	assert.NotEqual(t, tpo0, agent.newTimePointObject(&agent.time.clock, nil))

	agent.time.clock++
	assert.Nil(t, tpo0.join(agent.newTimePointObject(&agent.time.clock, nil)))
	tpo1 := agent.newTimePointObject(&agent.time.clock, map[int]any{partIdConceptContext: co0})
	tso := tpo0.join(tpo1)
	assert.Equal(t, co0, tso.ctx())

	cr := assertComparativeRelationSeen(t, tpo0.compare(tpo1), agent.logic.comparatives.lt, tpo0, tpo1)
	assert.Equal(t, co0, cr.ctx())

	assert.Equal(t, tso, agent.newTimeSegmentObject(tpo0, tpo1, nil))
	co1 := newTestContext(agent, 1)
	assert.NotNil(t, agent.newTimeSegmentObject(nil, nil, nil))
	assert.Nil(t, agent.newTimeSegmentObject(tpo0, tpo1, map[int]any{partIdConceptContext: co1}))
}

func TestTemporalObjectJoin(t *testing.T) {
	agent := newEmptyWorldAgent()
	v0, v1, v2, v3 := 0, 1, 2, 3
	c0, c1, c2, c3 := &v0, &v1, &v2, &v3
	tpo0 := agent.newTimePointObject(c0, nil)
	tpo1 := agent.newTimePointObject(c1, nil)
	tpo2 := agent.newTimePointObject(c2, nil)
	tpo3 := agent.newTimePointObject(c3, nil)
	tso01 := agent.newTimeSegmentObject(tpo0, tpo1, nil)
	tso02 := agent.newTimeSegmentObject(tpo0, tpo2, nil)
	tso03 := agent.newTimeSegmentObject(tpo0, tpo3, nil)
	tso12 := agent.newTimeSegmentObject(tpo1, tpo2, nil)
	tso13 := agent.newTimeSegmentObject(tpo1, tpo3, nil)
	tso23 := agent.newTimeSegmentObject(tpo2, tpo3, nil)
	assert.Equal(t, tpo0.join(tpo0), tpo0)
	assert.Equal(t, tpo1.join(tso01), tso01)
	assert.Equal(t, tpo0.join(tso12), tso02)
	assert.Equal(t, tpo3.join(tso12), tso13)
	assert.Equal(t, tso03.join(tso12), tso03)
	assert.Equal(t, tso03.join(tso23), tso03)
}

func TestTemporalObjectCompare(t *testing.T) {
	agent := newEmptyWorldAgent()
	v0, v1, v2, v3 := 0, 1, 2, 3
	c0, c1, c2, c3 := &v0, &v1, &v2, &v3
	tpo0 := agent.newTimePointObject(c0, nil)
	tpo1 := agent.newTimePointObject(c1, nil)
	tpo2 := agent.newTimePointObject(c2, nil)
	tpo3 := agent.newTimePointObject(c3, nil)
	tso01 := agent.newTimeSegmentObject(tpo0, tpo1, nil)
	tso02 := agent.newTimeSegmentObject(tpo0, tpo2, nil)
	tso03 := agent.newTimeSegmentObject(tpo0, tpo3, nil)
	tso12 := agent.newTimeSegmentObject(tpo1, tpo2, nil)
	tso13 := agent.newTimeSegmentObject(tpo1, tpo3, nil)
	tso23 := agent.newTimeSegmentObject(tpo2, tpo3, nil)
	assertIntervalRelationSeen(t, tso01.compare(tso01), agent.logic.intervals.eq, tso01, tso01)
	assertIntervalRelationSeen(t, tso01.compare(tso12), agent.logic.intervals.lt, tso01, tso12)
	assertIntervalRelationSeen(t, tso23.compare(tso12), agent.logic.intervals.gt, tso23, tso12)
	assertIntervalRelationSeen(t, tso03.compare(tso12), agent.logic.intervals.ct, tso03, tso12)
	assertIntervalRelationSeen(t, tso02.compare(tso03), agent.logic.intervals.ct, tso03, tso02)
	assertIntervalRelationSeen(t, tso13.compare(tso03), agent.logic.intervals.ct, tso03, tso13)
}
