package herbsystem

import "context"

type ContextKey string

const ContextKeyFinished = ContextKey("finished")

func WithFinished(ctx context.Context) context.Context {
	var finished bool
	return context.WithValue(ctx, ContextKeyFinished, &finished)
}

func IsFinished(ctx context.Context) bool {
	successed := ctx.Value(ContextKeyFinished).(*bool)
	return *successed
}

func Finish(ctx context.Context, s System) {
	v := ctx.Value(ContextKeyFinished).(*bool)
	*v = true
}
