package service

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func ConvertStringToType(str string, target any) (any, error) {
	targetType := reflect.TypeOf(target)

	isPtr := false
	if targetType.Kind() == reflect.Ptr {
		isPtr = true
		targetType = targetType.Elem() // Get the element type
	}

	ptrValue := reflect.New(targetType)

	switch targetType.Kind() {
	case reflect.Int:
		intValue, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to int: %v", str, err)
		}
		ptrValue.Elem().SetInt(int64(intValue))

	case reflect.Int8:
		intValue, err := strconv.ParseInt(str, 10, 8)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to int8: %v", str, err)
		}
		ptrValue.Elem().SetInt(int64(int8(intValue)))

	case reflect.Int16:
		intValue, err := strconv.ParseInt(str, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to int16: %v", str, err)
		}
		ptrValue.Elem().SetInt(int64(int16(intValue)))

	case reflect.Int32:
		intValue, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to int32: %v", str, err)
		}
		ptrValue.Elem().SetInt(int64(int32(intValue)))

	case reflect.Int64:
		intValue, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to int64: %v", str, err)
		}
		ptrValue.Elem().SetInt(intValue)

	case reflect.Uint:
		uintValue, err := strconv.ParseUint(str, 10, 0)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to uint: %v", str, err)
		}
		ptrValue.Elem().SetUint(uintValue)

	case reflect.Uint8:
		uintValue, err := strconv.ParseUint(str, 10, 8)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to uint8: %v", str, err)
		}
		ptrValue.Elem().SetUint(uintValue)

	case reflect.Uint16:
		uintValue, err := strconv.ParseUint(str, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to uint16: %v", str, err)
		}
		ptrValue.Elem().SetUint(uintValue)

	case reflect.Uint32:
		uintValue, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to uint32: %v", str, err)
		}
		ptrValue.Elem().SetUint(uintValue)

	case reflect.Uint64:
		uintValue, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to uint64: %v", str, err)
		}
		ptrValue.Elem().SetUint(uintValue)

	case reflect.Float32:
		floatValue, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to float32: %v", str, err)
		}
		ptrValue.Elem().SetFloat(float64(float32(floatValue)))

	case reflect.Float64:
		floatValue, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to float64: %v", str, err)
		}
		ptrValue.Elem().SetFloat(floatValue)

	case reflect.Bool:
		boolValue, err := strconv.ParseBool(str)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to bool: %v", str, err)
		}
		ptrValue.Elem().SetBool(boolValue)

	case reflect.String:
		ptrValue.Elem().SetString(str)

	case reflect.Struct:
		if targetType == reflect.TypeOf(time.Time{}) {
			timeValue, err := time.Parse(time.RFC3339, str)
			if err != nil {
				return nil, fmt.Errorf("cannot convert %s to time.Time: %v", str, err)
			}
			ptrValue.Elem().Set(reflect.ValueOf(timeValue))
		}

	default:
		return nil, fmt.Errorf("unsupported type: %s", targetType.Kind().String())
	}

	if isPtr {
		return ptrValue.Interface(), nil
	}
	return ptrValue.Elem().Interface(), nil
}
