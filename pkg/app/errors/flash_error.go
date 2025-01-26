package errors

import (
	"encoding/json"
)

func FlashErrorString(field BindError) string {
	jsonError, _ := json.Marshal(field)
	return string(jsonError)
}

func ParseFlashes(flashes []interface{}) map[string]*BindError {
	errorMessages := make(map[string]*BindError)
	for _, flash := range flashes {
		var fe BindError
		err := json.Unmarshal([]byte(flash.(string)), &fe)
		if err == nil {
			errorMessages[ToSnakeCase(fe.Field)] = &fe
		}
	}
	return errorMessages
}
