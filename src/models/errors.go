package models

import (
	"github.com/pkg/errors"
)

var (
	ErrNotFound            = errors.New("item is not found")
	ErrConflict            = errors.New("item already exists")
	ErrBadRequest          = errors.New("bad request")
	ErrInternalServerError = errors.New("internal server error")
)
