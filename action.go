package herbsystem

import (
	"context"
)

type Command string

type Handler func(ctx context.Context, next func(context.Context) error) error
type Action struct {
	Command Command
	Handler Handler
}

func NewAction() *Action {
	return &Action{}
}

func execActions(ctx context.Context, o []*Action) (context.Context, error) {
	if len(o) == 0 {
		return ctx, nil
	}
	var lastctx context.Context = ctx
	err := o[0].Handler(ctx, func(nextctx context.Context) error {
		var nexterr error
		lastctx, nexterr = execActions(nextctx, o[1:])
		return nexterr
	})
	if err != nil {
		return nil, err
	}
	return lastctx, nil
}
