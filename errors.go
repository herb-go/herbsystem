package herbsystem

import (
	"errors"
)

var ErrModuleNameDuplicated = errors.New("module name duplicated")

var ErrInvalidStage = errors.New("invalid stage")
