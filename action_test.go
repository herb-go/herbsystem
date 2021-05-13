package herbsystem

import "testing"

func TestCreateAction(t *testing.T) {
	var a *Action
	a = CreateStartAction(nil)
	if a.Command != CommandStart {
		t.Fatal()
	}
	a = CreateStopAction(nil)
	if a.Command != CommandStop {
		t.Fatal()
	}
}
