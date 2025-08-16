package errutil

import (
	"fmt"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func message(fieldErr validator.FieldError) (message string) {
	fieldName := ToSnakeCase(fieldErr.Field())
	tag := fieldErr.Tag()
	switch tag {
	case "required":
		message = fmt.Sprintf("%s is required", fieldName)
	case "email":
		message = fmt.Sprintf("%s must be a valid email", fieldName)
	case "min":
		message = fmt.Sprintf("%s must be at least %s characters long", fieldName, fieldErr.Param())
	case "max":
		message = fmt.Sprintf("%s must be at most %s characters long", fieldName, fieldErr.Param())
	default:
		message = fmt.Sprintf("%s is invalid", fieldName)
	}
	return
}

func ToSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if unicode.IsUpper(r) {
			// Добавляем нижнее подчеркивание перед каждой заглавной буквой, если это не начало строки
			if i > 0 && (unicode.IsLower(rune(str[i-1])) || (i < len(str)-1 && unicode.IsLower(rune(str[i+1])))) {
				result = append(result, '_')
			}
			r = unicode.ToLower(r)
		}
		result = append(result, r)
	}
	return string(result)
}
