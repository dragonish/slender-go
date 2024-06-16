package model

import (
	"errors"
)

// DB_USERNAME defines database users table key.
const DB_USERNAME = "slender"

var (
	// ErrExists indicates that the data already exists.
	ErrExists = errors.New("data already exists")

	// ErrNotExist indicates that the data does not exist.
	ErrNotExist = errors.New("data does not exist")

	// ErrDoNothing indicates that no operation.
	ErrDoNothing = errors.New("no operation")

	// ErrQueryParamMissing indicates that the required query conditions are missing
	ErrQueryParamMissing = errors.New("the required query conditions are missing")
)
