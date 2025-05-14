package service

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/google/uuid"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// LoadStruct TODO return error
func LoadStruct(dest, request any) any {
	if reflect.TypeOf(request).Kind() == reflect.Interface ||
		reflect.TypeOf(dest).Kind() == reflect.Interface ||
		reflect.TypeOf(request).Kind() != reflect.Ptr ||
		reflect.TypeOf(dest).Kind() != reflect.Ptr {
		return dest
	}

	// Если request nil, нужно установить dest тоже в nil
	if reflect.ValueOf(request).IsNil() {
		destVal := reflect.ValueOf(dest).Elem()
		// Проверяем, что dest — указатель и он может быть нулевым
		if destVal.CanSet() {
			destVal.Set(reflect.Zero(destVal.Type()))
		}
		return dest
	}

	requestVal := reflect.ValueOf(request).Elem() // Исходная структура (строки/базовые типы/указателя)
	destVal := reflect.ValueOf(dest).Elem()       // Целевая структура (с указателями и без)

	if reflect.TypeOf(requestVal).Kind() == reflect.Interface || reflect.TypeOf(destVal).Kind() == reflect.Interface {
		return dest
	}

	requestType := requestVal.Type()
	destType := destVal.Type()

	for i := 0; i < requestVal.NumField(); i++ {
		srcField := requestVal.Field(i)                                   // Поле в исходной структуре
		destField := destVal.FieldByName(requestVal.Type().Field(i).Name) // Поле в целевой структуре
		srcStructField := requestType.Field(i)
		destStructField, ok := destType.FieldByName(srcStructField.Name)
		var destTag reflect.StructTag
		if ok {
			destTag = destStructField.Tag
		}

		if destField.Kind() == reflect.Interface || srcField.Kind() == reflect.Interface {
			continue
		}

		if !destField.IsValid() || !destField.CanSet() || !requestVal.Type().Field(i).IsExported() {
			continue
		}

		if srcField.IsZero() {
			destField.SetZero()
			continue
		}

		// Если типы совпадают, просто копируем
		if srcField.Type() == destField.Type() {
			destField.Set(srcField)
			continue
		}

		// Если поля строки
		if srcField.Kind() == reflect.String {
			strToType(&srcField, &destField, destTag)
			continue
		}

		switch destField.Kind() {
		case reflect.Ptr:
			if destField.Type().Elem().Kind() == reflect.Struct && destField.Type().Elem() != reflect.TypeOf(time.Time{}) {
				if destField.IsNil() {
					// Создаем новую структуру, на которую будет указывать указатель
					destField.Set(reflect.New(destField.Type().Elem()))
				}
				if srcField.Kind() == reflect.Ptr {
					// Разыменовываем указатель и вызываем LoadStruct
					LoadStruct(destField.Interface(), srcField.Interface())
				}
				if srcField.Kind() == reflect.Struct {
					LoadStruct(destField.Interface(), srcField.Addr().Interface())
				}
				continue
			}
			typeToPtr(&srcField, &destField)
		case reflect.Struct:
			if srcField.Kind() == reflect.Struct {
				LoadStruct(destField.Addr().Interface(), srcField.Addr().Interface())
			}
			if srcField.Kind() == reflect.Ptr {
				LoadStruct(destField.Addr().Interface(), srcField.Interface())
			}
		case reflect.Slice:
			if srcField.Kind() != reflect.Slice {
				continue
			}

			newDst := reflect.MakeSlice(destField.Type(), srcField.Len(), srcField.Len())
			srcFieldSliceType := srcField.Type().Elem()
			destFieldSliceType := destField.Type().Elem()

			if srcFieldSliceType.Kind() == reflect.Struct {
				for i := 0; i < srcField.Len(); i++ {
					// Получаем элемент исходного слайса
					srcElem := srcField.Index(i)

					if destFieldSliceType.Kind() == reflect.Ptr {
						// Создаем указатель на новую структуру для целевого слайса
						destElemPtr := reflect.New(destFieldSliceType.Elem()) // Создаем указатель на структуру целевого типа

						// Рекурсивно копируем поля структуры из исходного элемента в элемент, на который указывает указатель
						LoadStruct(destElemPtr.Interface(), srcElem.Addr().Interface())

						// Устанавливаем указатель в целевой слайс
						newDst.Index(i).Set(destElemPtr)
					} else {
						// Создаем новый элемент для целевого слайса
						destElem := reflect.New(destFieldSliceType).Elem() // Создаем пустую структуру целевого типа

						// Рекурсивно копируем поля структуры из исходного элемента в целевой элемент
						LoadStruct(destElem.Addr().Interface(), srcElem.Addr().Interface())

						// Устанавливаем элемент в целевой слайс
						newDst.Index(i).Set(destElem)
					}
				}

				// Устанавливаем новый слайс в целевое поле
				destField.Set(newDst)
				continue
			}

			if srcFieldSliceType.Kind() == reflect.Ptr {
				if destFieldSliceType.Kind() == reflect.Interface {
					continue
				}
				for i := 0; i < srcField.Len(); i++ {
					// Получаем элемент исходного слайса
					srcElem := srcField.Index(i)

					if srcElem.IsNil() {
						newDst.Index(i).Set(reflect.Zero(destFieldSliceType))
						continue
					}

					// Разыменовываем указатель, чтобы получить структуру
					srcElemStruct := srcElem.Elem()
					if destFieldSliceType.Kind() == reflect.Ptr {
						// Создаем новый указатель на структуру для целевого слайса
						destElemPtr := reflect.New(destFieldSliceType.Elem()) // Создаем указатель на структуру

						// Рекурсивно копируем поля из разыменованного указателя исходного слайса в разыменованный указатель целевого слайса
						LoadStruct(destElemPtr.Interface(), srcElemStruct.Addr().Interface())

						// Устанавливаем указатель на структуру в целевой слайс
						newDst.Index(i).Set(destElemPtr)
					} else {
						// Создаем новый элемент (структуру) для целевого слайса
						destElem := reflect.New(destFieldSliceType).Elem()

						// Рекурсивно копируем поля из разыменованного указателя в структуру
						LoadStruct(destElem.Addr().Interface(), srcElemStruct.Addr().Interface())

						// Устанавливаем элемент в целевой слайс
						newDst.Index(i).Set(destElem)
					}
				}

				// Устанавливаем новый слайс в целевое поле
				destField.Set(newDst)
				continue
			}

			if srcFieldSliceType.Kind() == destFieldSliceType.Kind() {
				reflect.Copy(newDst, srcField)
				destField.Set(newDst)
				continue
			}

			if srcFieldSliceType.Kind() == reflect.String {
				for i := 0; i < srcField.Len(); i++ {
					srcFieldSliceTypeElem := srcField.Index(i)
					destFieldSliceTypeElem := newDst.Index(i)
					strToType(&srcFieldSliceTypeElem, &destFieldSliceTypeElem, destTag)
				}
				destField.Set(newDst)
				continue
			}

			if destFieldSliceType.Kind() == reflect.Ptr {
				for i := 0; i < srcField.Len(); i++ {
					srcFieldSliceTypeElem := srcField.Index(i)
					destFieldSliceTypeElem := newDst.Index(i)
					typeToPtr(&srcFieldSliceTypeElem, &destFieldSliceTypeElem)
				}
				destField.Set(newDst)
				continue
			}
		}
	}
	return dest
}

func strToType(src *reflect.Value, dest *reflect.Value, destTag reflect.StructTag) {
	if src.Kind() != reflect.String {
		return
	}
	srcStr := strings.TrimSpace(src.String()) // Значение строки из request
	switch dest.Kind() {
	case reflect.Ptr:
		strToPtr(src, dest, destTag)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val, err := strconv.ParseInt(srcStr, 10, dest.Type().Bits()); err == nil {
			dest.SetInt(val)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val, err := strconv.ParseUint(srcStr, 10, dest.Type().Bits()); err == nil {
			dest.SetUint(val)
		}
	case reflect.Float32, reflect.Float64:
		if val, err := strconv.ParseFloat(srcStr, dest.Type().Bits()); err == nil {
			dest.SetFloat(val)
		}
	case reflect.Bool:
		if val, err := strconv.ParseBool(srcStr); err == nil {
			dest.SetBool(val)
		}
	case reflect.Struct:
		if dest.Type() == reflect.TypeOf(time.Time{}) {
			timeFormat := destTag.Get("time_format")
			if timeFormat == "" {
				timeFormat = "2006-01-02 15:04:05"
			}
			if val, err := time.Parse(timeFormat, srcStr); err == nil {
				dest.Set(reflect.ValueOf(val))
			}
		}
	case reflect.Array:
		if dest.Type() == reflect.TypeOf(uuid.UUID{}) {
			trimmed := strings.TrimSpace(src.String())
			if val, err := uuid.Parse(trimmed); err == nil {
				dest.Set(reflect.ValueOf(val))
			} else {
				logger.Errorf("[LoadStruct] Error parsing UUID '%s': %v", srcStr, err)
			}
		}
	}
}

func strToPtr(src *reflect.Value, dest *reflect.Value, destTag reflect.StructTag) {
	if src.Kind() != reflect.String || dest.Kind() != reflect.Ptr {
		return
	}
	srcStr := strings.TrimSpace(src.String())
	destElemType := dest.Type().Elem() // Тип значения, на которое указывает указатель

	switch destElemType.Kind() {
	case reflect.Array:
		if destElemType == reflect.TypeOf(uuid.UUID{}) {
			if parsed, err := uuid.Parse(srcStr); err == nil {
				newVal := reflect.New(destElemType).Elem()
				newVal.Set(reflect.ValueOf(parsed))
				dest.Set(newVal.Addr())
			} else {
				logger.Errorf("[LoadStruct] Error parsing UUID '%s': %v", srcStr, err)
			}
			return
		}
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
			timeFormat := destTag.Get("time_format")
			if timeFormat == "" {
				timeFormat = "2006-01-02 15:04:05"
			}
			if val, err := time.Parse(timeFormat, srcStr); err == nil {
				newVal := reflect.New(destElemType).Elem()
				newVal.Set(reflect.ValueOf(val))
				dest.Set(newVal.Addr())
			}
		}
	default:
		panic("unhandled default case")
	}
}

func typeToPtr(src *reflect.Value, dest *reflect.Value) {
	destElemType := dest.Type().Elem() // Тип значения, на которое указывает указатель
	if dest.Kind() != reflect.Ptr || src.Type() != destElemType {
		return
	}

	newValue := reflect.New(dest.Type().Elem())
	newValue.Elem().Set(*src)
	dest.Set(newValue)
}
