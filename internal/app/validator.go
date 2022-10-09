package app

import "github.com/go-playground/validator/v10"

type CustomValidator struct {
	v *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.v.Struct(i)
}
