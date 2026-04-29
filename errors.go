package annon

import "errors"

var (
	ErrInvalidJSON = errors.New("invalid json input")
	ErrInvalidYAML = errors.New("invalid yaml input")
)
