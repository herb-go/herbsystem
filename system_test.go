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

type testModule struct {
	NopModule
	name    string
	actions []*Action
}

func (s *testModule) InstallProcess(ctx context.Context, system System, next func(context.Context, System)) {
	system.MountSystemActions(s.actions...)
	next(ctx, system)
}
func (s *testModule) ModuleName() string {
	return s.name
}
func TestStage(t *testing.T) {
	var err error
	s := New()
	if s.SystemStage() != StageNew {
		t.Fatal(s)
	}
	err = Catch(func() {
		MustExecActions(context.TODO(), s, "")
	})
	if err == nil || errors.Unwrap(err) != ErrInvalidStage {
		t.Fatal(err)
	}
	s.SetSystemStage(Stage(-1))
	err = Catch(func() { s.MustRegisterSystemModule(nil) })
	if err == nil || errors.Unwrap(err) != ErrInvalidStage {
		t.Fatal(err)
	}
	err = Catch(func() { MustReady(s) })
	if err == nil || errors.Unwrap(err) != ErrInvalidStage {
		t.Fatal(err)
	}
	err = Catch(func() { MustConfigure(s) })
	if err == nil || errors.Unwrap(err) != ErrInvalidStage {
		t.Fatal(err)
	}
	err = Catch(func() { MustGetConfigurableModule(s, "test") })
	if err == nil || errors.Unwrap(err) != ErrInvalidStage {
		t.Fatal(err)
	}
	err = Catch(func() { MustStart(s) })
	if err == nil || errors.Unwrap(err) != ErrInvalidStage {
		t.Fatal(err)
	}
	err = Catch(func() { MustStop(s) })
	if err == nil || errors.Unwrap(err) != ErrInvalidStage {
		t.Fatal(err)
	}
}

func TestNopModule(t *testing.T) {
	var err error
	s := New()
	s.MustRegisterSystemModule(NopModule{})
	err = Catch(func() { s.MustRegisterSystemModule(NopModule{}) })
	if err == nil || errors.Unwrap(err) != ErrModuleNameDuplicated {
		panic(err)
	}
	s = New()
	s.MustRegisterSystemModule(NopModule{})
	MustReady(s)
	MustConfigure(s)
	m := MustGetConfigurableModule(s, "")
	if m == nil {
		t.Fatal(m)
	}
	m = MustGetConfigurableModule(s, "notexists")
	if m != nil {
		t.Fatal(m)
	}
	MustStart(s)
	MustExecActions(context.TODO(), s, "")
	MustStop(s)
}

func TestAction(t *testing.T) {
	var err error
	action1 := NewAction()
	action1.Command = "test"
	action1.Process = func(ctx context.Context, system System, next func(context.Context, System)) {
		result := getResult(ctx)
		*result = append(*result, "action1")
		next(ctx, system)
	}
	action2 := NewAction()
	action2.Command = "test"
	action2.Process = func(ctx context.Context, system System, next func(context.Context, System)) {
		result := getResult(ctx)
		*result = append(*result, "action2")
		next(ctx, system)
	}
	action3 := NewAction()
	action3.Command = "test2"
	action3.Process = func(ctx context.Context, system System, next func(context.Context, System)) {
		result := getResult(ctx)
		*result = append(*result, "action3")
	}
	action4 := NewAction()
	action4.Command = "test"
	action4.Process = func(ctx context.Context, system System, next func(context.Context, System)) {
		result := getResult(ctx)
		*result = append(*result, "action4")
		next(context.WithValue(ctx, "last", true), system)
	}
	action5 := NewAction()
	action5.Command = "test2"
	action5.Process = func(ctx context.Context, system System, next func(context.Context, System)) {
		result := getResult(ctx)
		*result = append(*result, "action5")
		next(ctx, system)
	}
	action6 := NewAction()
	action6.Command = "test3"
	action6.Process = func(ctx context.Context, system System, next func(context.Context, System)) {
		result := getResult(ctx)
		*result = append(*result, "action6")
		panic(errors.New("stop"))
	}
	action7 := NewAction()
	action7.Command = "test3"
	action7.Process = func(ctx context.Context, system System, next func(context.Context, System)) {
		result := getResult(ctx)
		*result = append(*result, "action7")
		next(ctx, system)
	}
	module1 := &testModule{
		name: "server1",
		actions: []*Action{
			action1,
			action2,
		},
	}
	module2 := &testModule{
		name: "server2",
		actions: []*Action{
			action3,
			action4,
			action5,
			action6,
			action7,
		},
	}
	s := New()
	s.MustRegisterSystemModule(module1)
	s.MustRegisterSystemModule(module2)
	MustReady(s)
	MustConfigure(s)
	MustStart(s)
	ctx := MustExecActions(resultContext(), s, "test")
	if !IsFinished(ctx) {
		t.Fatal()
	}
	r := getResult(ctx)
	if len(*r) != 3 || (*r)[0] != "action1" || (*r)[1] != "action2" || (*r)[2] != "action4" {
		t.Fatal(*r)
	}
	last := ctx.Value("last").(bool)
	if !last {
		t.Fatal(last)
	}
	ctx = MustExecActions(resultContext(), s, "test2")

	if ctx != nil {
		t.Fatal()
	}
	rctx := resultContext()
	err = Catch(func() {
		MustExecActions(rctx, s, "test3")
	})
	if err == nil {
		t.Fatal(ctx, err)
	}
	r = getResult(rctx)
	if len(*r) != 1 || (*r)[0] != "action6" {
		t.Fatal(*r)
	}
	MustStop(s)
}
