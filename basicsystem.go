package herbsystem

import (
	"context"
	"fmt"
)

type BasicSystem struct {
	stage     Stage
	modules   []Module
	processes map[interface{}][]Process
	logger    func(error)
}

func (s *BasicSystem) SystemStage() Stage {
	return s.stage
}

func (s *BasicSystem) SetSystemStage(stage Stage) {
	s.stage = stage
}
func (s *BasicSystem) SystemModules() []Module {
	return s.modules
}
func (s *BasicSystem) SystemContext() context.Context {
	return context.TODO()
}
func (s *BasicSystem) MustRegisterSystemModule(m Module) {
	PanicIfNotInStage(s, StageNew)
	name := m.ModuleName()
	for _, v := range s.modules {
		if v.ModuleName() == name {
			panic(fmt.Errorf("herbsystem: %w (%s)", ErrModuleNameDuplicated, v.ModuleName()))
		}
	}
	s.modules = append(s.modules, m)
}
func (s *BasicSystem) GetSystemActions(cmd interface{}) []Process {
	return s.processes[cmd]
}
func (s *BasicSystem) MountSystemActions(actions ...*Action) {
	for _, v := range actions {
		s.processes[v.Command] = append(s.processes[v.Command], v.Process)
	}
}
func (s *BasicSystem) ResetSystem() {
	s.processes = map[interface{}][]Process{}
}
func (s *BasicSystem) SetSystemLogger(l func(error)) {
	s.logger = l
}
func (s *BasicSystem) LogSystemError(err error) {
	errs := ToErrors(err).All()
	for _, err := range errs {
		s.logger(err)
	}
}
func New() *BasicSystem {
	return &BasicSystem{
		stage:     StageNew,
		processes: map[interface{}][]Process{},
		logger:    DefaultLogger,
	}
}
