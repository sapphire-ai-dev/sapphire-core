package agent

type selfObject struct {
	*abstractObject
	worldId int
}

func (o *selfObject) match(other concept) bool {
	n, ok := other.(*selfObject)
	return ok && o.worldId == n.worldId
}

func (o *selfObject) debugArgs() map[string]any {
	args := o.abstractObject.debugArgs()
	args["worldId"] = o.worldId
	return args
}

func (a *Agent) newSelfObject(worldId int) *selfObject {
	result := &selfObject{
		worldId: worldId,
	}
	a.newAbstractObject(result, &result.abstractObject)
	return result.memorize().(*selfObject)
}
