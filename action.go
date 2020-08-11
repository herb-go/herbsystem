package herbsystem

import (
	"context"
)

type Command string

type Action struct {
	Command Command
	Handler func(ctx context.Context, next func(context.Context) error) error
}

func NewAction() *Action {
	return &Action{}
}

func nextAction(o []*Action) func(context.Context) error {
	return func(ctx context.Context) error {
		return execActions(ctx, o)
	}
}

func execActions(ctx context.Context, o []*Action) error {
	if len(o) == 0 {
		return nil
	}
	return o[0].Handler(ctx, nextAction(o[1:]))
}
