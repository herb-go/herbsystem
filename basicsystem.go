package herbsystem

import "fmt"

type BasicSystem struct {
	stage     Stage
	modules   []Module
	processes map[interface{}][]Process
}

func (s *BasicSystem) Stage() Stage {
	return s.stage
}

func (s *BasicSystem) SetStage(stage Stage) {
	s.stage = stage
}
func (s *BasicSystem) Modules() []Module {
	return s.modules
}
func (s *BasicSystem) MustRegisterModule(m Module) {
	PanicIfNotInStage(s, StageNew)
	for _, v := range s.modules {
		if v.ModuleName() == m.ModuleName() {
			panic(fmt.Errorf("herbsystem: %w (%s)", ErrModuleNameDuplicated, v.ModuleName()))
		}
	}
	s.modules = append(s.modules, m)
}
func (s *BasicSystem) GetActions(cmd interface{}) []Process {
	return s.processes[cmd]
}
func (s *BasicSystem) MountActions(actions ...*Action) {
	for _, v := range actions {
		s.processes[v.Command] = append(s.processes[v.Command], v.Process)
	}
}
func (s *BasicSystem) Reset() {
	s.processes = map[interface{}][]Process{}
}

func New() *BasicSystem {
	return &BasicSystem{
		stage:     StageNew,
		processes: map[interface{}][]Process{},
	}
}
