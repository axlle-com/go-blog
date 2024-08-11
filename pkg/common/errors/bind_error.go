package errors

import (
	"fmt"
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
			field := toSnakeCase(fieldErr.Field())
			errors[field] = BindError{
				Field:   field,
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

func ParseBindErrorToMap(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			field := toSnakeCase(fieldErr.Field())
			errors[field] = message(fieldErr)
		}
	} else {
		errors[GeneralFieldName] = fmt.Sprintf("%s", err.Error())
	}

	return errors
}
