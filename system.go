package herbsystem

import (
	"context"
	"fmt"
)

type System interface {
	SystemStage() Stage
	SetSystemStage(stage Stage)
	SystemModules() []Module
	MustRegisterSystemModule(m Module)
	GetSystemActions(cmd interface{}) []Process
	MountSystemActions(actions ...*Action)
	ResetSystem()
	SystemContext() context.Context
	SetSystemLogger(func(error))
	LogSystemError(error)
}

func PanicIfNotInStage(s System, stage Stage) {
	ss := s.SystemStage()
	if stage != ss {
		name := StageNames[ss]
		if name == "" {
			panic(fmt.Errorf("herbsystem: %w (%d)", ErrInvalidStage, ss))
		}
		panic(fmt.Errorf("herbsystem: %w (%s)", ErrInvalidStage, name))
	}
}

func MustExecActions(ctx context.Context, s System, cmd interface{}) context.Context {
	PanicIfNotInStage(s, StageRunning)
	p := s.GetSystemActions(cmd)
	var result context.Context
	ComposeProcess(p...)(WithFinished(ctx), s, func(newctx context.Context, s System) {
		result = newctx
		Finish(newctx, s)
	})
	return result
}

func MustGetConfigurableModule(s System, name string) Module {
	PanicIfNotInStage(s, StageConfiguring)
	for _, v := range s.SystemModules() {
		if v.ModuleName() == name {
			return v
		}
	}
	return nil
}

func initSystemModules(s System) {
	for _, v := range s.SystemModules() {
		v.InitModule()
	}
}
func MustReady(s System) {
	PanicIfNotInStage(s, StageNew)
	initSystemModules(s)
	s.SetSystemStage(StageReady)
}

func MustConfigure(s System) {
	PanicIfNotInStage(s, StageReady)
	s.ResetSystem()
	modules := s.SystemModules()
	processes := make([]Process, len(modules))
	for k := range modules {
		processes[k] = modules[k].InitProcess
	}
	ComposeProcess(processes...)(WithFinished(s.SystemContext()), s, Finish)
	s.SetSystemStage(StageConfiguring)
}

func MustStart(s System) {
	PanicIfNotInStage(s, StageConfiguring)
	ComposeProcess(s.GetSystemActions(CommandStart)...)(WithFinished(s.SystemContext()), s, Finish)
	s.SetSystemStage(StageRunning)
}

func MustStop(s System) {
	PanicIfNotInStage(s, StageRunning)
	ComposeProcess(s.GetSystemActions(CommandStop)...)(WithFinished(s.SystemContext()), s, Finish)
	initSystemModules(s)
	s.SetSystemStage(StageReady)
}
