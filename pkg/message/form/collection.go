package form

import "errors"

type Factory func() any

var Registry = map[string]Factory{
	(&Contact{}).Name(): func() any { return &Contact{} },
}

func NewForm(name string) (any, error) {
	constructor, ok := Registry[name]
	if !ok {
		return nil, errors.New("unknown form: " + name)
	}

	return constructor(), nil
}
