package m

import (
	"errors"
)

var (
	ErrPID     = errors.New("package.pid is invalid")
	ErrMID     = errors.New("package.mid is invalid")
	ErrAccount = errors.New("account is invalid")
	ErrVCode   = errors.New("vcode is invalid")
)
