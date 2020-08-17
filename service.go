package herbsystem

type Service interface {
	InitService() error
	ServiceName() string
	StartService() error
	StopService() error
	ConfigurService() error
	ServiceActions() []*Action
}

type NopService struct {
}

func (s NopService) InitService() error {
	return nil
}
func (s NopService) ServiceName() string {
	return ""
}
func (s NopService) StartService() error {
	return nil
}
func (s NopService) StopService() error {
	return nil
}
func (s NopService) ServiceActions() []*Action {
	return nil
}
func (s NopService) ConfigurService() error {
	return nil
}
