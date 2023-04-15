package agent

type number struct {
	*abstractConcept
	value int
}

func (c *number) match(other concept) bool {
	o, ok := other.(*number)
	return ok && c.value == o.value
}

func (a *Agent) newNumber(value int, args map[int]any) *number {
	result := &number{
		value: value,
	}

	a.newAbstractConcept(result, args, &result.abstractConcept)
	return result.memorize().(*number)
}
