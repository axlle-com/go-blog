package errors

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

func message(fieldErr validator.FieldError) (message string) {
	fieldName := fieldErr.Field()
	tag := fieldErr.Tag()
	switch tag {
	case "required":
		message = fmt.Sprintf("%s is required", fieldName)
	case "email":
		message = fmt.Sprintf("%s must be a valid email", fieldName)
	case "min":
		message = fmt.Sprintf("%s must be at least %s characters long", fieldName, fieldErr.Param())
	case "max":
		message = fmt.Sprintf("%s must be at most %s characters long", fieldName, fieldErr.Param())
	default:
		message = fmt.Sprintf("%s is invalid", fieldName)
	}
	return
}
