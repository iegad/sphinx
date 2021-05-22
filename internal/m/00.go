package m

import (
	"errors"
)

var (
	ErrPID = errors.New("package.pid is invalid")
	ErrMID = errors.New("package.mid is invalid")
)
