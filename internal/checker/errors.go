package checker

import "errors"

var (
	ErrUnknownType = errors.New("checker: unknown check type")
	ErrTimeout     = errors.New("checker: timeout")
	ErrNetwork     = errors.New("checker: network error")
)
