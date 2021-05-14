package herbsystem

import "context"

type Module interface {
	ModuleName() string
	InitModule()
	InstallProcess(ctx context.Context, system System, next func(context.Context, System))
}

type NopModule struct {
}

func (NopModule) ModuleName() string {
	return ""
}
func (NopModule) InitModule() {

}
func (NopModule) InstallProcess(ctx context.Context, system System, next func(context.Context, System)) {
	next(ctx, system)
}
