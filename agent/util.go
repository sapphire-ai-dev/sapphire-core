package agent

import (
	"fmt"
	"github.com/sapphire-ai-dev/sapphire-core/world"
	"reflect"
	"sort"
)

type TestActionInterface struct {
	ReadyResult bool
	StepCount   int
}

func (t *TestActionInterface) instantiate() *world.ActionInterface {
	return &world.ActionInterface{
		Name: reflect.TypeOf(t).Name(),
		Ready: func() bool {
			return t.ReadyResult
		},
		Step: func() {
			t.StepCount++
		},
	}
}

func newTestActionInterface() *TestActionInterface {
	return &TestActionInterface{
		ReadyResult: false,
		StepCount:   0,
	}
}

func matchParams(a, b map[string]any) bool {
	for paramName := range a {
		if a[paramName] != b[paramName] {
			return false
		}
	}

	for paramName := range b {
		if a[paramName] != b[paramName] {
			return false
		}
	}

	return true
}

func toReflect[T any]() reflect.Type {
	return reflect.TypeOf(func(T) {}).In(0)
}

func ternary(b any) *bool {
	if b != nil {
		if bb, ok := b.(bool); ok {
			return &bb
		}
	}

	return nil
}

func sortSlice[T any](lst []T, scores []float64) []T {
	var lstCopy []T
	for _, t := range lst {
		lstCopy = append(lstCopy, t)
	}
	sort.SliceStable(lstCopy, func(i, j int) bool {
		return scores[i] > scores[j]
	})
	return lstCopy
}

func interpretPart[T concept](concepts map[int]concept, partId int) T {
	var result T
	if c, seen := concepts[partId]; seen {
		if r, ok := c.(T); ok {
			result = r
		}
	}

	return result
}

func ternaryEqual(a, b *bool) bool {
	if a == nil || b == nil {
		return false
	}

	return *a == *b
}

func printErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
