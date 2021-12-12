package herbsystem

import "strings"

type Errors struct {
	Err  error
	Prev *Errors
}

func (e *Errors) NilOrError() error {
	if e == nil {
		return nil
	}
	return e
}
func (e *Errors) Error() string {
	errs := e.All()
	msgs := []string{}
	for _, v := range errs {
		msgs = append(msgs, v.Error())
	}
	return strings.Join(msgs, "\n")
}

func (e *Errors) Append(err error) *Errors {
	return &Errors{
		Err:  err,
		Prev: e,
	}
}

func (e *Errors) All() []error {
	if e == nil {
		return nil
	}
	return append(e.Prev.All(), e.Err)
}

func ToErrors(err error) *Errors {
	if err == nil {
		return nil
	}
	e, ok := err.(*Errors)
	if ok {
		return e
	}
	return &Errors{Err: err}
}

func NewErrors() *Errors {
	return nil
}
