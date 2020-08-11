package herbsystem

import (
	"context"
	"fmt"
)

type System struct {
	Stage    Stage
	services []Service
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
	s.Stage = StageConfiguring
	return nil
}
func (s *System) Start() error {
	err := s.ErrIfNotInStage(StageConfiguring)
	s.Stage = StageStarting
	if err != nil {
		return err
	}
	for _, v := range s.services {
		err = v.StartService()
		if err != nil {
			return err
		}
	}
	s.Stage = StageRunning
	return nil
}
func (s *System) Stop() error {
	err := s.ErrIfNotInStage(StageRunning)
	s.Stage = StageStoping
	if err != nil {
		return err
	}
	for _, v := range s.services {
		err = v.StopService()
		if err != nil {
			return err
		}
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
	for _, v := range s.services {
		if v.ServiceName() == name {
			return v, nil
		}
	}
	return nil, nil
}
func NewSystem() *System {
	return &System{
		Stage:   StageNew,
		actions: map[Command][]*Action{},
	}
}
