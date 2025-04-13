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

type Errors struct {
	Errors  map[string]string
	Message string
}

type MapErrors map[string]string

func ParseBindError(err error) map[string]BindError {
	errors := make(map[string]BindError)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			field := ToSnakeCase(fieldErr.Field())
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

func ParseBindErrorToMap(err error) *Errors {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			field := ToSnakeCase(fieldErr.Field())
			errors[field] = message(fieldErr)
		}
	} else {
		errors[GeneralFieldName] = err.Error()
	}

	return &Errors{
		Errors:  errors,
		Message: "Ошибки валидации",
	}
}
