package herbsystem

import "context"

type Process func(ctx context.Context, next func(context.Context))

func ComposeProcess(series ...Process) Process {
	return func(ctx context.Context, next func(context.Context)) {
		if len(series) == 0 {
			next(ctx)
			return
		}
		series[0](ctx, func(context.Context) {
			ComposeProcess(series[1:]...)(ctx, next)
		})
	}
}
