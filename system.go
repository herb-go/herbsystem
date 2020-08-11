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
		err = v.Start()
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
		err = v.Stop()
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
	name := service.Name()
	for _, v := range s.services {
		if v.Name() == name {
			return fmt.Errorf("herbsystem: %w (%s)", ErrServiceNameDuplicated, name)
		}
	}
	err = service.Init()
	if err != nil {
		return err
	}
	s.services = append(s.services, service)
	o := service.Actions()
	for _, v := range o {
		s.actions[v.Command] = append(s.actions[v.Command], v)
	}
	return nil
}

func (s *System) ExecActions(ctx context.Context, cmd Command) error {
	err := s.ErrIfNotInStage(StageRunning)
	if err != nil {
		return err
	}
	return execActions(ctx, s.actions[cmd])
}
func (s *System) GetConfigurableService(name string) (Service, error) {
	err := s.ErrIfNotInStage(StageConfiguring)
	if err != nil {
		return nil, err
	}
	for _, v := range s.services {
		if v.Name() == name {
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
