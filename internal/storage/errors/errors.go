package storage_errors

import (
	"errors"
)

var (
	ErrUniqueConstraintViolation = errors.New("unique violation constraint")
)
