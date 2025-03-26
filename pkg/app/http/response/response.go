package response

import (
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	"net/http"
)

type Body map[string]any
type message map[string]any

func Message(message string) map[string]any {
	return map[string]any{
		"message": message,
	}
}

func OK(data any, message string, pagination contracts.Paginator) map[string]any {
	return Successful(http.StatusOK, data, message, pagination)
}

func Created(data any, message string) map[string]any {
	return Successful(http.StatusCreated, data, message, nil)
}

func Successful(code int, data any, message string, pagination contracts.Paginator) map[string]any {
	result := Body{
		"successful": true,
		"code":       code,
	}
	body := Body{
		"result":  result,
		"data":    data,
		"errors":  nil,
		"message": message,
	}

	if pagination != nil {
		body["pagination"] = pagination
	}

	return body
}

func Fail(code int, message string, errors map[string]string) map[string]any {
	result := Body{
		"successful": false,
		"code":       code,
	}
	body := Body{
		"result":     result,
		"data":       nil,
		"pagination": nil,
		"errors":     errors,
		"message":    message,
	}

	return body
}
