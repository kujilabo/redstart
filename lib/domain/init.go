package domain

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type ContextKey string

var (
	Validator = validator.New()

	ErrInvalidArgument = errors.New("invalid argument")
)
