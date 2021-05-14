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

func Wrap(h ...func()) Process {
	return func(ctx context.Context, system System, next func(context.Context, System)) {
		for k := range h {
			h[k]()
		}
		next(ctx, system)
	}
}

func WrapOrPanic(h ...func() error) Process {
	return func(ctx context.Context, system System, next func(context.Context, System)) {
		for k := range h {
			if err := h[k](); err != nil {
				panic(err)
			}

		}
		next(ctx, system)
	}
}

func WrapOrLog(h ...func() error) Process {
	return func(ctx context.Context, system System, next func(context.Context, System)) {
		for k := range h {
			if err := h[k](); err != nil {
				system.LogSystemError(err)
			}

		}
		next(ctx, system)
	}
}
