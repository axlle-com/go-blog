package models

import (
	"reflect"
)

type Field struct {
}

func (f *Field) SetEmptyPointersToNil(v interface{}) {
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() != reflect.Ptr || field.IsNil() {
			continue
		}

		elem := field.Elem()
		switch elem.Kind() {
		case reflect.String:
			if elem.String() == "" {
				field.Set(reflect.Zero(field.Type()))
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if elem.Uint() == 0 {
				field.Set(reflect.Zero(field.Type()))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if elem.Int() == 0 {
				field.Set(reflect.Zero(field.Type()))
			}
		case reflect.Float32, reflect.Float64:
			if elem.Float() == 0 {
				field.Set(reflect.Zero(field.Type()))
			}
		case reflect.Struct:
			if elem.IsZero() {
				field.Set(reflect.Zero(field.Type()))
			}
		case reflect.Bool:
			if !elem.Bool() {
				field.Set(reflect.Zero(field.Type()))
			}
		default:
			panic("unhandled default case")
		}
	}
}
