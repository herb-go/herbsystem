package herbsystem

import (
	"context"
	"fmt"
)

type System struct {
	Stage    Stage
	services []Service
	locked   map[string]bool
	actions  map[Command][]*Action
}

func (s *System) ErrIfNotInStage(stage Stage) error {
	if stage != s.Stage {
		name := StageNames[s.Stage]
		if name == "" {
			return fmt.Errorf("herbsystem: %w (%d)", ErrInvalidStage, s.Stage)
		}
		return fmt.Errorf("herbsystem: %w (%s)", ErrInvalidStage, name)
	}
	return nil
}
func (s *System) Ready() error {
	err := s.ErrIfNotInStage(StageNew)
	if err != nil {
		return err
	}
	s.Stage = StageReady
	return nil
}
func (s *System) Configuring() error {
	err := s.ErrIfNotInStage(StageReady)
	if err != nil {
		return err
	}
	errs := NewErrors()
	for _, v := range s.services {
		errs.Add(v.ConfigurService())
	}
	err = errs.ToError()
	if err != nil {
		return err
	}
	s.Stage = StageConfiguring
	s.locked = map[string]bool{}
	return nil
}
func (s *System) Start() error {
	err := s.ErrIfNotInStage(StageConfiguring)
	if err != nil {
		return err
	}
	errs := NewErrors()
	for _, v := range s.services {
		errs.Add(v.StartService())
	}
	err = errs.ToError()
	if err != nil {
		return err
	}
	s.Stage = StageRunning
	return nil
}
func (s *System) Stop() error {
	err := s.ErrIfNotInStage(StageRunning)
	if err != nil {
		return err
	}
	s.Stage = StageStoping
	errs := NewErrors()
	for _, v := range s.services {
		errs.Add(v.StopService())
	}
	err = errs.ToError()
	if err != nil {
		return err
	}

	s.Stage = StageReady
	return nil
}
func (s *System) InstallService(service Service) error {
	err := s.ErrIfNotInStage(StageNew)
	if err != nil {
		return err
	}
	name := service.ServiceName()
	for _, v := range s.services {
		if v.ServiceName() == name {
			return fmt.Errorf("herbsystem: %w (%s)", ErrServiceNameDuplicated, name)
		}
	}
	err = service.InitService()
	if err != nil {
		return err
	}
	s.services = append(s.services, service)
	o := service.ServiceActions()
	for _, v := range o {
		s.actions[v.Command] = append(s.actions[v.Command], v)
	}
	return nil
}

func (s *System) ExecActions(ctx context.Context, cmd Command) (context.Context, error) {
	err := s.ErrIfNotInStage(StageRunning)
	if err != nil {
		return nil, err
	}
	return execActions(ctx, s.actions[cmd])
}
func (s *System) GetConfigurableService(name string) (Service, error) {
	err := s.ErrIfNotInStage(StageConfiguring)
	if err != nil {
		return nil, err
	}
	if s.locked[name] {
		return nil, nil
	}
	for _, v := range s.services {
		if v.ServiceName() == name {
			return v, nil
		}
	}
	return nil, nil
}
func (s *System) LockConfigurableService(name string) error {
	_, err := s.GetConfigurableService(name)
	if err != nil {
		return err
	}
	s.locked[name] = true
	return nil
}
func NewSystem() *System {
	return &System{
		Stage:   StageNew,
		actions: map[Command][]*Action{},
	}
}
