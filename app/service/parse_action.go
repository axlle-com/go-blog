package service

import (
	"fmt"
)

func ParseAction(key string, obj map[string]any) (string, error) {
	rawAction, ok := obj[key]
	if !ok {
		return "", fmt.Errorf("%s key missing", key)
	}

	action, ok := rawAction.(string)
	if !ok {
		return "", fmt.Errorf("%s is not a string, got %T", key, rawAction)
	}

	return action, nil
}
