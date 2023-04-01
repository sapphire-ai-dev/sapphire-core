package agent

type temporalObject interface {
	object
	start() *timePointObject
	end() *timePointObject
	join(other temporalObject) temporalObject
	compare(other temporalObject) map[int]relation
}

func nillableIntCopy(a *int) *int {
	if a == nil {
		return nil
	}

	b := *a
	return &b
}

func nillableIntEqual(a, b *int) bool {
	if a == nil || b == nil {
		return false
	}

	return *a == *b
}

type timePointObject struct {
	*abstractObject
	clockTime *int
}

func (o *timePointObject) match(other concept) bool {
	n, ok := other.(*timePointObject)
	return ok && o.abstractObject.match(n.abstractObject) && nillableIntEqual(o.clockTime, n.clockTime)
}

func (o *timePointObject) start() *timePointObject {
	return o
}

func (o *timePointObject) end() *timePointObject {
	return o
}

func (o *timePointObject) join(other temporalObject) temporalObject {
	return o.agent.time.temporalObjJoin(o, other)
}

func (o *timePointObject) compare(other temporalObject) map[int]relation {
	return o.agent.time.temporalObjCompare(o, other)
}

func (a *Agent) newTimePointObject(clockTime *int, args map[int]any) *timePointObject {
	result := &timePointObject{clockTime: nillableIntCopy(clockTime)}
	a.newAbstractObject(result, args, &result.abstractObject)
	return result.memorize().(*timePointObject)
}

type timeSegmentObject struct {
	*abstractObject
	_start *memReference
	_end   *memReference
}

func (o *timeSegmentObject) match(other concept) bool {
	n, ok := other.(*timeSegmentObject)
	return ok && o.abstractObject.match(n.abstractObject) && o._start.c.match(n._start.c) && o._end.c.match(n._end.c)
}

func (o *timeSegmentObject) start() *timePointObject {
	return parseRef[*timePointObject](o.agent, o._start)
}

func (o *timeSegmentObject) end() *timePointObject {
	return parseRef[*timePointObject](o.agent, o._end)
}

func (o *timeSegmentObject) join(other temporalObject) temporalObject {
	return o.agent.time.temporalObjJoin(o, other)
}

func (o *timeSegmentObject) compare(other temporalObject) map[int]relation {
	return o.agent.time.temporalObjCompare(o, other)
}

func (a *Agent) newTimeSegmentObject(start, end *timePointObject, args map[int]any) *timeSegmentObject {
	if start == nil || end == nil || !matchConcepts(start.ctx(), end.ctx()) {
		return nil
	}

	collision := false
	if args, collision = injectConceptArg(args, conceptArgContext, start.ctx()); collision {
		return nil
	}
	result := &timeSegmentObject{}
	a.newAbstractObject(result, args, &result.abstractObject)
	result._start = start.createReference(result, true)
	result._end = end.createReference(result, true)
	return result.memorize().(*timeSegmentObject)
}
