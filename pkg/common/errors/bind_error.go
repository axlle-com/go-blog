package errors

import (
	"github.com/go-playground/validator/v10"
)

const GeneralFieldName = "general"

type BindError struct {
	Field   string
	Message string
	Value   interface{}
}

func ParseBindError(err error) map[string]BindError {
	errors := make(map[string]BindError)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			errors[fieldErr.Field()] = BindError{
				Field:   fieldErr.Field(),
				Message: message(fieldErr),
				Value:   fieldErr.Value(),
			}
		}
	} else {
		errors[GeneralFieldName] = BindError{
			Field:   GeneralFieldName,
			Message: err.Error(),
		}
	}

	return errors
}
