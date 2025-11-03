package service

import (
	"encoding/json"
	"errors"

	"gorm.io/datatypes"
)

var ErrEmpty = errors.New("empty json")

// UnmarshalJSON разбирает datatypes.JSON или *datatypes.JSON в произвольный тип.
func UnmarshalJSON[T any](j *datatypes.JSON, out *T) error {
	if j == nil || len(*j) == 0 {
		return ErrEmpty
	}
	return json.Unmarshal(*j, out)
}
