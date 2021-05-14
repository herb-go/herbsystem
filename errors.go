package herbsystem

import (
	"errors"
	"log"
)

var ErrModuleNameDuplicated = errors.New("module name duplicated")

var ErrInvalidStage = errors.New("invalid stage")

func Catch(f func()) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()
	f()
	return nil
}

var DefaultLogger = func(err error) {
	log.Println(err.Error())
}
