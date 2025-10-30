package video

import "errors"

var (
	ErrOSNotSupported = errors.New("OS currently not supported")
	ErrInvalidPath    = errors.New("the current path is invalid")
)
