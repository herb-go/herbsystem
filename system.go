package herbsystem

import (
	"context"
	"fmt"
)

type System interface {
	Stage() Stage
	SetStage(stage Stage)
	Modules() []Module
	MustRegisterModule(m Module)
	GetActions(cmd interface{}) []Process
	MountActions(actions ...*Action)
	Reset()
}

func PanicIfNotInStage(s System, stage Stage) {
	ss := s.Stage()
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
	p := s.GetActions(cmd)
	var result context.Context
	ComposeProcess(p...)(WithFinished(ctx), s, func(newctx context.Context, s System) {
		result = newctx
		Finish(newctx, s)
	})
	return result
}

func MustGetConfigurableModule(s System, name string) Module {
	PanicIfNotInStage(s, StageConfiguring)
	for _, v := range s.Modules() {
		if v.ModuleName() == name {
			return v
		}
	}
	return nil
}

func MustReady(s System) {
	PanicIfNotInStage(s, StageNew)
	s.SetStage(StageReady)
}

func MustConfigure(s System) {
	PanicIfNotInStage(s, StageReady)
	s.Reset()
	modules := s.Modules()
	processes := make([]Process, len(modules))
	for k := range modules {
		processes[k] = modules[k].InitProcess
	}
	ComposeProcess(processes...)(WithFinished(context.TODO()), s, Finish)
	s.SetStage(StageConfiguring)
}

func MustStart(s System) {
	PanicIfNotInStage(s, StageConfiguring)
	ComposeProcess(s.GetActions(CommandStart)...)(WithFinished(context.TODO()), s, Finish)
	s.SetStage(StageRunning)
}

func MustStop(s System) {
	PanicIfNotInStage(s, StageRunning)
	ComposeProcess(s.GetActions(CommandStop)...)(WithFinished(context.TODO()), s, Finish)
	s.SetStage(StageReady)
}
