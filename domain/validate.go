package domain

import (
	"gopkg.in/go-playground/validator.v9"
)

var Validator *validator.Validate

func init() {
	Validator = validator.New()
}
