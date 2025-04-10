package models

import (
	"github.com/axlle-com/blog/app/logger"
	"reflect"
)

type Field struct {
}

func (f *Field) GetChangedFields(original, updated interface{}) map[string]interface{} {
	changedFields := make(map[string]interface{})

	originalValue := reflect.ValueOf(original)
	updatedValue := reflect.ValueOf(updated)

	if originalValue.Kind() == reflect.Ptr {
		originalValue = originalValue.Elem()
	}
	if updatedValue.Kind() == reflect.Ptr {
		updatedValue = updatedValue.Elem()
	}

	if originalValue.Kind() != reflect.Struct {
		if updatedValue.Kind() != reflect.Struct {
			logger.Error("Both original and updated must be structs")
			return changedFields
		} else {
			return f.fillUpdatedValue(updatedValue)
		}
	}

	if originalValue.Type() != updatedValue.Type() {
		logger.Error("The types of the original and updated structs do not match")
		return changedFields
	}

	if originalValue.IsZero() {
		return f.fillUpdatedValue(updatedValue)
	}

	numFields := originalValue.NumField()
	for i := 0; i < numFields; i++ {
		field := originalValue.Type().Field(i)

		if !field.IsExported() {
			continue
		}

		if tagValue, ok := field.Tag.Lookup("ignore"); ok && tagValue == "true" {
			continue
		}

		originalField := originalValue.Field(i)
		updatedField := updatedValue.Field(i)

		if !originalField.IsValid() || !updatedField.IsValid() {
			continue
		}

		if !originalField.CanInterface() || !updatedField.CanInterface() {
			continue
		}

		if !reflect.DeepEqual(originalField.Interface(), updatedField.Interface()) {
			changedFields[field.Name] = updatedField.Interface()
		}
	}

	return changedFields
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
		}
	}
}

func SetEmptyStringPointersToNil(v interface{}) {
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.String {
			str := field.Elem().String()
			if str == "" {
				field.Set(reflect.Zero(field.Type()))
			}
		}
	}
}

func SetZeroPointersToNil(v interface{}) {
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.Ptr && field.Elem().CanUint() {
			str := field.Elem().Uint()
			if str == 0 {
				field.Set(reflect.Zero(field.Type()))
			}
		}
	}
}

func (f *Field) fillUpdatedValue(updatedValue reflect.Value) map[string]interface{} {
	changedFields := make(map[string]interface{})
	numFields := updatedValue.NumField()
	for i := 0; i < numFields; i++ {
		field := updatedValue.Type().Field(i)
		if !field.IsExported() {
			continue
		}
		if tagValue, ok := field.Tag.Lookup("ignore"); ok && tagValue == "true" {
			continue
		}
		updatedField := updatedValue.Field(i)
		if !updatedField.IsValid() || !updatedField.CanInterface() {
			continue
		}
		changedFields[field.Name] = updatedField.Interface()
	}
	return changedFields
}
