package entity

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrInvalidInput   = errors.New("invalid input")
	DataInconsistency = errors.New("data inconsistency")
)
