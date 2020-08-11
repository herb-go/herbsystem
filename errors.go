package herbsystem

import "errors"

var ErrServiceNameDuplicated = errors.New("service name duplicated")

var ErrInvalidStage = errors.New("invalid stage")
