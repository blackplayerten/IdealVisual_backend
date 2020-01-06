package account

import (
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
)

type UniqueConstraintViolationError struct {
	Field string
}

func (u UniqueConstraintViolationError) Error() string {
	return "unique constraint violated"
}

type NotNullConstraintViolationError struct {
	Field string
}

func (n NotNullConstraintViolationError) Error() string {
	return "not null constraint violated"
}
