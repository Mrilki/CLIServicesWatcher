package config

import "errors"

var (
	ErrNotFound = errors.New("config: file not found")
	ErrRead     = errors.New("config: failed to read file")
	ErrValidate = errors.New("config: validation failed")
	ErrParse    = errors.New("config: failed to parse file")
)
