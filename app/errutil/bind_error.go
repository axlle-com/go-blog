package errutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

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
	if err == nil {
		return &Errors{Errors: errs}
	}

	// 1) Validation errors (go-playground/validator)
	var vErrs validator.ValidationErrors
	if errors.As(err, &vErrs) {
		for _, fe := range vErrs {
			field := ToSnakeCase(fe.Field())
			errs[field] = append(errs[field], message(fe))
		}
	}

	// 2) Type mismatch during JSON unmarshalling
	var ute *json.UnmarshalTypeError
	if errors.As(err, &ute) {
		// e.g. "menu_items.id" -> "id"
		field := ToSnakeCase(lastSegment(ute.Field))
		expected := humanType(ute.Type)
		actual := ute.Value
		if field == "" {
			field = "json"
		}
		errs[field] = append(errs[field],
			fmt.Sprintf("expected type %s but received %q", expected, actual))
	}

	// 3) JSON syntax error
	var se *json.SyntaxError
	if errors.As(err, &se) {
		errs["json"] = append(errs["json"],
			fmt.Sprintf("JSON syntax error (offset %d)", se.Offset))
	}

	// 4) If no specific type recognized — keep the raw message
	if len(errs) == 0 {
		errs["_error"] = append(errs["_error"], err.Error())
	}

	return &Errors{
		Errors:  errs,
		Message: "Validation errors",
	}
}

// lastSegment("MenuItemsRequest.menu_items.id") => "id"
func lastSegment(path string) string {
	if path == "" {
		return ""
	}
	i := strings.LastIndex(path, ".")
	if i >= 0 && i+1 < len(path) {
		return path[i+1:]
	}
	return path
}

// Returns a human-readable type name
func humanType(t reflect.Type) string {
	if t == nil {
		return "unknown"
	}
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Map, reflect.Struct:
		return "object"
	case reflect.String:
		return "string"
	default:
		return t.String()
	}
}
