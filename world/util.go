package world

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidArgs = errors.New("invalid args")
	ErrInvalidCmd  = errors.New("invalid cmd")
)

func GetArg[T any](position int, required bool, backup T, args []any) T {
	if position >= 0 && position < len(args) {
		if tArg, ok := args[position].(T); ok {
			return tArg
		}
	}

	if required {
		panic(ErrInvalidArgs)
	}

	return backup
}

func SetOutArg[T any](position int, required bool, value T, args []any) {
	if position >= 0 && position < len(args) {
		if out, ok := args[position].(*T); ok {
			*out = value
			return
		}
	}

	if required {
		panic(ErrInvalidArgs)
	}
}

func PrintErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
