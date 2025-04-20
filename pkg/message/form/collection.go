package form

import (
	"errors"
	"reflect"
)

var Registry = map[string]reflect.Type{
	(&Contact{}).Name(): reflect.TypeOf(Contact{}),
}

type Name struct {
	FormName string `json:"form_name" binding:"required"`
}

func (n *Name) NewForm() (any, error) {
	typ, ok := Registry[n.FormName]
	if !ok {
		return nil, errors.New("unknown form: " + n.FormName)
	}

	return reflect.New(typ).Interface(), nil
}
