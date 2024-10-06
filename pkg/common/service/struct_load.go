package service

import (
	"github.com/axlle-com/blog/pkg/common/logger"
	"reflect"
	"strconv"
	"time"
)

// LoadFromRequest TODO return error
func LoadFromRequest(dest, request any) any {
	requestVal := reflect.ValueOf(request).Elem() // Исходная структура (строки/базовые типы/указателя)
	destVal := reflect.ValueOf(dest).Elem()       // Целевая структура (с указателями и без)

	for i := 0; i < requestVal.NumField(); i++ {
		srcField := requestVal.Field(i)                                   // Поле в исходной структуре
		destField := destVal.FieldByName(requestVal.Type().Field(i).Name) // Поле в целевой структуре

		// Проверяем, что поле экспортируемое и его можно установить
		if !(requestVal.Type().Field(i).IsExported() && destField.IsValid() && destField.CanSet()) {
			continue
		}

		// Если типы совпадают, просто копируем
		if srcField.Type() == destField.Type() {
			destField.Set(srcField)
			continue
		}

		if srcField.IsZero() {
			destField.SetZero()
			continue
		}

		srcStr := srcField.String() // Значение строки из request
		switch destField.Kind() {
		case reflect.Ptr:
			if srcField.Kind() == reflect.String {
				strPtr(&srcField, &destField)
			} else {
				typePtr(&srcField, &destField)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if val, err := strconv.ParseInt(srcStr, 10, destField.Type().Bits()); err == nil {
				destField.SetInt(val)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if val, err := strconv.ParseUint(srcStr, 10, destField.Type().Bits()); err == nil {
				destField.SetUint(val)
			}
		case reflect.Float32, reflect.Float64:
			if val, err := strconv.ParseFloat(srcStr, destField.Type().Bits()); err == nil {
				destField.SetFloat(val)
			}
		case reflect.Bool:
			if val, err := strconv.ParseBool(srcStr); err == nil {
				destField.SetBool(val)
			}
		case reflect.Struct:
			if destField.Type() == reflect.TypeOf(time.Time{}) {
				if val, err := time.Parse("2006-01-02 15:04:05", srcStr); err == nil {
					destField.Set(reflect.ValueOf(val))
				}
			}
		}
	}
	return dest
}

func strPtr(src *reflect.Value, dest *reflect.Value) {
	if src.Kind() != reflect.String || dest.Kind() != reflect.Ptr {
		return
	}
	srcStr := src.String()
	destElemType := dest.Type().Elem() // Тип значения, на которое указывает указатель
	switch destElemType.Kind() {
	case reflect.String:
		newVal := reflect.New(destElemType).Elem()
		newVal.SetString(srcStr)
		dest.Set(newVal.Addr())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val, err := strconv.ParseInt(srcStr, 10, destElemType.Bits()); err == nil {
			newVal := reflect.New(destElemType).Elem()
			newVal.SetInt(val)
			dest.Set(newVal.Addr())
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val, err := strconv.ParseUint(srcStr, 10, destElemType.Bits()); err == nil {
			newVal := reflect.New(destElemType).Elem()
			newVal.SetUint(val)
			dest.Set(newVal.Addr())
		}
	case reflect.Float32, reflect.Float64:
		if val, err := strconv.ParseFloat(srcStr, destElemType.Bits()); err == nil {
			newVal := reflect.New(destElemType).Elem()
			newVal.SetFloat(val)
			dest.Set(newVal.Addr())
		}
	case reflect.Bool:
		if val, err := strconv.ParseBool(srcStr); err == nil {
			newVal := reflect.New(destElemType).Elem()
			newVal.SetBool(val)
			dest.Set(newVal.Addr())
		}
	case reflect.Struct:
		if destElemType == reflect.TypeOf(time.Time{}) {
			if val, err := time.Parse("2006-01-02 15:04:05", srcStr); err == nil {
				logger.Print(val)
				newVal := reflect.New(destElemType).Elem()
				newVal.Set(reflect.ValueOf(val))
				dest.Set(newVal.Addr())
			}
			logger.Print("+")
		}
	}
}

func typePtr(src *reflect.Value, dest *reflect.Value) {
	destElemType := dest.Type().Elem() // Тип значения, на которое указывает указатель
	if dest.Kind() != reflect.Ptr || src.Type() != destElemType {
		return
	}

	newValue := reflect.New(dest.Type().Elem())
	newValue.Elem().Set(*src)
	dest.Set(newValue)
}
