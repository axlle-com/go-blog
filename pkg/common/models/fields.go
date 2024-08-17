package models

import (
	"errors"
	"github.com/axlle-com/blog/pkg/common/logger"
	"reflect"
)

func GetChangedFields(original, updated interface{}) map[string]interface{} {
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
			logger.Error(errors.New("Both original and updated must be structs"))
			return changedFields
		} else {
			return fillUpdatedValue(updatedValue)
		}
	}

	if originalValue.Type() != updatedValue.Type() {
		logger.Error(errors.New("The types of the original and updated structs do not match"))
		return changedFields
	}

	if originalValue.IsZero() {
		return fillUpdatedValue(updatedValue)
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

func fillUpdatedValue(updatedValue reflect.Value) map[string]interface{} {
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
