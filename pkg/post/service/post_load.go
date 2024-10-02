package service

import (
	"reflect"
)

func LoadFromRequest(dest, request any) any {
	requestVal := reflect.ValueOf(request).Elem()
	destVal := reflect.ValueOf(dest).Elem()

	for i := 0; i < requestVal.NumField(); i++ {
		srcField := requestVal.Field(i)
		destField := destVal.FieldByName(requestVal.Type().Field(i).Name)
		if requestVal.Type().Field(i).IsExported() && destField.IsValid() && destField.CanSet() && srcField.Type() == destField.Type() {
			destField.Set(srcField)
		}
	}
	return dest
}
