package engine

import (
	"errors"
)

var errInterrupted = errors.New("workflow engine: interrupted, may be ended")
