package herbsystem

import (
	"context"
	"errors"
	"testing"
)

func resultContext() context.Context {
	result := &[]string{}
	return context.WithValue(context.Background(), "result", result)
}

func getResult(ctx context.Context) *[]string {
	v := ctx.Value("result")
	return v.(*[]string)
}

type testService struct {
	NopService
	name    string
	actions []*Action
}

func (s *testService) ServiceActions() []*Action {
	return s.actions
}
func (s *testService) ServiceName() string {
	return s.name
}
func TestStage(t *testing.T) {
	var err error
	s := NewSystem()
	if s.Stage != StageNew {
		t.Fatal(s)
	}
	_, err = s.ExecActions(nil, "")
	if err == nil || errors.Unwrap(err) != ErrInvalidStage {
		t.Fatal(err)
	}
	s.Stage = Stage(-1)
	err = s.InstallService(nil)
	if err == nil || errors.Unwrap(err) != ErrInvalidStage {
		t.Fatal(err)
	}
	err = s.Ready()
	if err == nil || errors.Unwrap(err) != ErrInvalidStage {
		t.Fatal(err)
	}
	err = s.Configuring()
	if err == nil || errors.Unwrap(err) != ErrInvalidStage {
		t.Fatal(err)
	}
	_, err = s.GetConfigurableService("test")
	if err == nil || errors.Unwrap(err) != ErrInvalidStage {
		t.Fatal(err)
	}
	err = s.Start()
	if err == nil || errors.Unwrap(err) != ErrInvalidStage {
		t.Fatal(err)
	}
	err = s.Stop()
	if err == nil || errors.Unwrap(err) != ErrInvalidStage {
		t.Fatal(err)
	}
}

func TestNopService(t *testing.T) {
	var err error
	s := NewSystem()
	err = s.InstallService(NopService{})
	if err != nil {
		panic(err)
	}
	err = s.InstallService(NopService{})
	if err == nil || errors.Unwrap(err) != ErrServiceNameDuplicated {
		panic(err)
	}
	err = s.Ready()
	if err != nil {
		panic(err)
	}
	err = s.Configuring()
	if err != nil {
		panic(err)
	}
	service, err := s.GetConfigurableService("")
	if err != nil {
		panic(err)
	}
	if service == nil {
		t.Fatal(service)
	}
	service, err = s.GetConfigurableService("notexists")
	if err != nil {
		panic(err)
	}
	if service != nil {
		t.Fatal(service)
	}
	err = s.Start()
	if err != nil {
		panic(err)
	}
	_, err = s.ExecActions(nil, "")
	if err != nil {
		panic(err)
	}
	err = s.Stop()
	if err != nil {
		panic(err)
	}
}

func TestAction(t *testing.T) {
	var err error
	action1 := NewAction()
	action1.Command = "test"
	action1.Handler = func(ctx context.Context, next func(context.Context) error) error {
		result := getResult(ctx)
		*result = append(*result, "action1")
		return next(ctx)
	}
	action2 := NewAction()
	action2.Command = "test"
	action2.Handler = func(ctx context.Context, next func(context.Context) error) error {
		result := getResult(ctx)
		*result = append(*result, "action2")
		return next(ctx)
	}
	action3 := NewAction()
	action3.Command = "test2"
	action3.Handler = func(ctx context.Context, next func(context.Context) error) error {
		result := getResult(ctx)
		*result = append(*result, "action3")
		return nil
	}
	action4 := NewAction()
	action4.Command = "test"
	action4.Handler = func(ctx context.Context, next func(context.Context) error) error {
		result := getResult(ctx)
		*result = append(*result, "action4")
		return next(context.WithValue(ctx, "last", true))
	}
	action5 := NewAction()
	action5.Command = "test2"
	action5.Handler = func(ctx context.Context, next func(context.Context) error) error {
		result := getResult(ctx)
		*result = append(*result, "action5")
		return next(ctx)
	}
	action6 := NewAction()
	action6.Command = "test3"
	action6.Handler = func(ctx context.Context, next func(context.Context) error) error {
		result := getResult(ctx)
		*result = append(*result, "action6")
		return errors.New("stop")
	}
	action7 := NewAction()
	action7.Command = "test3"
	action7.Handler = func(ctx context.Context, next func(context.Context) error) error {
		result := getResult(ctx)
		*result = append(*result, "action7")
		return next(ctx)
	}
	servece1 := &testService{
		name: "server1",
		actions: []*Action{
			action1,
			action2,
		},
	}
	servece2 := &testService{
		name: "server2",
		actions: []*Action{
			action3,
			action4,
			action5,
			action6,
			action7,
		},
	}
	s := NewSystem()
	s.InstallService(servece1)
	s.InstallService(servece2)
	err = s.Ready()
	if err != nil {
		panic(err)
	}
	err = s.Configuring()
	if err != nil {
		panic(err)
	}

	err = s.Start()
	if err != nil {
		panic(err)
	}
	ctx, err := s.ExecActions(resultContext(), "test")
	if err != nil {
		panic(err)
	}
	r := getResult(ctx)
	if len(*r) != 3 || (*r)[0] != "action1" || (*r)[1] != "action2" || (*r)[2] != "action4" {
		t.Fatal(*r)
	}
	last := ctx.Value("last").(bool)
	if !last {
		t.Fatal(last)
	}
	ctx, err = s.ExecActions(resultContext(), "test2")
	if err != nil {
		panic(err)
	}
	r = getResult(ctx)
	if len(*r) != 1 || (*r)[0] != "action3" {
		t.Fatal(*r)
	}
	rctx := resultContext()
	ctx, err = s.ExecActions(rctx, "test3")
	if ctx != nil || err == nil {
		t.Fatal(ctx, err)
	}
	r = getResult(rctx)
	if len(*r) != 1 || (*r)[0] != "action6" {
		t.Fatal(*r)
	}
	err = s.Stop()
	if err != nil {
		panic(err)
	}
}
