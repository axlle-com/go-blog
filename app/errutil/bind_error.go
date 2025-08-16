package errutil

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

const GeneralFieldName = "general"

type BindError struct {
	Field   string
	Message string
	Value   interface{}
}

type Errors struct {
	Errors  map[string][]string `json:"errors"`
	Message string              `json:"message"`
}

type MapErrors map[string]string

func NewBindError(err error) map[string]BindError {
	errs := make(map[string]BindError)
	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) {
		for _, fieldErr := range validationErrors {
			field := ToSnakeCase(fieldErr.Field())
			errs[field] = BindError{
				Field:   field,
				Message: message(fieldErr),
				Value:   fieldErr.Value(),
			}
		}
	}

	return errs
}

func NewErrors(err error) *Errors {
	errs := make(map[string][]string)
	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) {
		for _, fieldErr := range validationErrors {
			field := ToSnakeCase(fieldErr.Field())
			if _, ok := errs[field]; !ok {
				errs[field] = make([]string, 0)
			}
			errs[field] = append(errs[field], message(fieldErr))
		}
	}

	return &Errors{
		Errors:  errs,
		Message: "Ошибки валидации",
	}
}
