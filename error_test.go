package herbsystem

import (
	"errors"
	"testing"
)

func TestErrors(t *testing.T) {
	errs := NewErrors()
	if errs.ToError() != nil {
		t.Fatal()
	}
	errs.Add(errors.New("a"))
	errs.Add(errors.New("b"))
	errs.Add(nil)
	if errs.ToError().Error() != "a | b" {
		t.Fatal()
	}
	var testerr = errors.New("test")
	var test2err = errors.New("test2")
	err := MergeError(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = MergeError(testerr, nil)
	if err != testerr {
		t.Fatal(err)
	}
	err = MergeError(testerr, test2err)
	if err.Error() != "test | test2" {
		t.Fatal(err)
	}
	errs = NewErrors()
	errs.Add(testerr)
	err = MergeError(errs.ToError(), test2err)
	if err.Error() != "test | test2" {
		t.Fatal(err)
	}
	errs = NewErrors()
	errs.Add(test2err)
	err = MergeError(testerr, errs.ToError())
	if err.Error() != "test | test2" {
		t.Fatal(err)
	}
}
