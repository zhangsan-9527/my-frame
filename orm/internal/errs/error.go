package errs

import "errors"

var (
	ErrPointerOnly = errors.New("orm: 只支持指向结构体的一级指针")
)
