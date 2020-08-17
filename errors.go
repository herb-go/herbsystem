package herbsystem

import (
	"errors"
	"strings"
)

var ErrServiceNameDuplicated = errors.New("service name duplicated")

var ErrInvalidStage = errors.New("invalid stage")

type Errors struct {
	Errors []error
}

func (e *Errors) ToError() error {
	if len(e.Errors) == 0 {
		return nil
	}
	return &ErrorsError{
		Errors: e.Errors,
	}
}

func (e *Errors) Add(err error) {
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
}

type ErrorsError struct {
	Errors []error
}

func (e *ErrorsError) Error() string {
	errormsgs := make([]string, len(e.Errors))
	for k := range e.Errors {
		errormsgs[k] = e.Errors[k].Error()
	}
	return strings.Join(errormsgs, " | ")
}

func NewErrors() *Errors {
	return &Errors{
		Errors: []error{},
	}
}

func MergeError(err error, newerr error) error {
	var errs *ErrorsError
	var ok bool
	if err == nil {
		return newerr
	}
	if newerr == nil {
		return err
	}
	errs, ok = err.(*ErrorsError)
	if ok {
		errs.Errors = append(errs.Errors, newerr)
		return errs
	}
	errs, ok = newerr.(*ErrorsError)
	if ok {
		var newerrs = append([]error{}, err)
		errs.Errors = append(newerrs, errs.Errors...)
		return errs
	}
	return &ErrorsError{
		Errors: []error{err, newerr},
	}
}
