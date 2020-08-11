package herbsystem

type Service interface {
	Init() error
	Name() string
	Start() error
	Stop() error
	Actions() []*Action
}
