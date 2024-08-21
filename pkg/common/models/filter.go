package models

import (
	"errors"
	"fmt"
	errorsForm "github.com/axlle-com/blog/pkg/common/errors"
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type Filter struct {
	array       map[string]string // TODO map[string][]string
	queryString template.URL
	query       url.Values
	context     *gin.Context
}

func (f *Filter) ValidateForm(ctx *gin.Context, model interface{}) *errorsForm.Errors {
	f.context = ctx
	err := ctx.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		return &errorsForm.Errors{Message: "Форма не валидная!"}
	}

	if len(ctx.Request.PostForm) == 0 {
		return &errorsForm.Errors{Message: "Форма не валидная!"}
	}

	if err := ctx.ShouldBind(model); err != nil {
		errBind := errorsForm.ParseBindErrorToMap(err)
		return errBind
	}
	f.SetEmptyPointersToNil(model)
	f.setQuery().
		setQueryString().
		setMap(model).
		addQueryString(f.mapToQueryString())

	if f.IsEmpty() {
		return &errorsForm.Errors{Message: "Форма пустая"}
	}
	return nil
}

func (f *Filter) ValidateQuery(ctx *gin.Context, model interface{}) *errorsForm.Errors {
	f.context = ctx
	if err := ctx.ShouldBindQuery(model); err != nil {
		errBind := errorsForm.ParseBindErrorToMap(err)
		return errBind
	}
	f.SetEmptyPointersToNil(model)
	f.setQuery().setQueryString().setMap(model)
	return nil
}

func (f *Filter) IsEmpty() bool {
	return f.array == nil || len(f.array) == 0
}

func (f *Filter) GetMap() map[string]string {
	if f.array == nil {
		return nil
	}
	return f.array
}

func (f *Filter) GetQueryString() template.URL {
	return f.queryString
}

func (f *Filter) setMap(model interface{}) *Filter {
	result := make(map[string]string)
	v := reflect.ValueOf(model)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return f
	}

	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)
		if !structField.IsExported() {
			continue
		}

		if tagValue, ok := structField.Tag.Lookup("ignore"); ok && tagValue == "true" {
			continue
		}

		field := v.Field(i)
		if field.Kind() == reflect.Ptr && field.IsNil() {
			continue
		}

		fieldName, ok := structField.Tag.Lookup("json")
		if !ok {
			fieldName = errorsForm.ToSnakeCase(v.Type().Field(i).Name)
		}

		if s := f.processFieldValue(field); s != "" {
			result[fieldName] = s
		}
	}
	f.array = result
	return f
}

func (f *Filter) setQuery() *Filter {
	f.query = f.context.Request.URL.Query()
	return f
}

func (f *Filter) setQueryString() *Filter {
	q := make(url.Values)
	for f, v := range f.context.Request.URL.Query() {
		if f == "page" || f == "pageSize" {
			continue
		}
		q[f] = v
	}
	f.queryString = template.URL(q.Encode())
	return f
}

func (f *Filter) addQueryString(s string) *Filter {
	f.queryString += template.URL(s)
	return f
}

func (f *Filter) processFieldValue(field reflect.Value) string {
	var fieldValue string

	switch field.Kind() {
	case reflect.Ptr:
		if field.IsNil() {
			return ""
		}
		return f.processFieldValue(field.Elem())

	case reflect.String:
		fieldValue = field.String()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fieldValue = strconv.FormatInt(field.Int(), 10)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fieldValue = strconv.FormatUint(field.Uint(), 10)

	case reflect.Float32, reflect.Float64:
		fieldValue = strconv.FormatFloat(field.Float(), 'f', -1, 64)

	case reflect.Bool:
		fieldValue = strconv.FormatBool(field.Bool())

	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.Uint8 {
			fieldValue = string(field.Bytes())
		} else {
			fieldValue = "unsupported slice type"
		}

	case reflect.Map:
		keys := field.MapKeys()
		mapResult := "{"
		for i, key := range keys {
			mapKey := f.processFieldValue(key)
			mapValue := f.processFieldValue(field.MapIndex(key))
			if i > 0 {
				mapResult += ", "
			}
			mapResult += mapKey + ": " + mapValue
		}
		mapResult += "}"
		fieldValue = mapResult

	case reflect.Struct:
		return ""

	default:
		logger.Error(errors.New(fmt.Sprintf("Unsupported field type: %v", field.Kind())))
	}

	return fieldValue
}

func (f *Filter) mapToQueryString() string {
	var queryParts []string

	for key, value := range f.array {
		queryParts = append(queryParts, url.QueryEscape(key)+"="+url.QueryEscape(value))
	}

	return strings.Join(queryParts, "&")
}

func (f *Filter) SetEmptyPointersToNil(v interface{}) {
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
		case reflect.Bool:
			if !elem.Bool() {
				field.Set(reflect.Zero(field.Type()))
			}
		}
	}
}