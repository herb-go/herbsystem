package herbsystem

import "context"

func NopReveiver(context.Context, System) {}

type Process func(ctx context.Context, system System, next func(context.Context, System))

func ComposeProcess(series ...Process) Process {
	return func(ctx context.Context, system System, next func(context.Context, System)) {
		if len(series) == 0 {
			next(ctx, system)
			return
		}
		series[0](ctx, system, func(ctx context.Context, system System) {
			ComposeProcess(series[1:]...)(ctx, system, next)
		})
	}
}
