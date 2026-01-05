package app

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	_validator *validator.Validate
	once       sync.Once
)

func InitValidator() *validator.Validate {
	once.Do(func() {
		_validator = validator.New()
	})

	return _validator
}
