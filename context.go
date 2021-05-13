package herbsystem

import "context"

type ContextName string

const ContextNameFinished = ContextName("finished")

func WithFinished(ctx context.Context) context.Context {
	var finished bool
	return context.WithValue(ctx, ContextNameFinished, &finished)
}

func IsFinished(ctx context.Context) bool {
	successed := ctx.Value(ContextNameFinished).(*bool)
	return *successed
}

func Finish(ctx context.Context, s System) {
	v := ctx.Value(ContextNameFinished).(*bool)
	*v = true
}
