package herbsystem

import "context"

type Module interface {
	ModuleName() string
	InitProcess(ctx context.Context, system System, next func(context.Context, System))
}

type NopModule struct {
}

func (NopModule) ModuleName() string {
	return ""
}

func (NopModule) InitProcess(ctx context.Context, system System, next func(context.Context, System)) {
	next(ctx, system)
}
