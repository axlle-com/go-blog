package db

import (
	"fmt"
	"github.com/axlle-com/blog/pkg/app/logger"
	"math/rand"
	"strconv"
	"time"
)

func RandBool() bool {
	value := rand.Intn(2)
	result := value != 0
	return result
}

func IntToBoolPtr() *bool {
	value := rand.Intn(2)
	result := value != 0
	return &result
}

func BoolToBoolPtr(b bool) *bool {
	return &b
}

func ParseDate(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}

	layout := "02.01.2006"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return &date
}

func StrPtr(s string) *string {
	return &s
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

func UintPtr(u int) *uint {
	v := uint(u)
	return &v
}

func IntPtr(i int) *int {
	return &i
}

func Int8Ptr(i int8) *int8 {
	return &i
}

func Int16Ptr(i int16) *int16 {
	return &i
}

func Int32Ptr(i int32) *int32 {
	return &i
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func Float32Ptr(i float32) *float32 {
	return &i
}

func Float64Ptr(i float64) *float64 {
	return &i
}

func IntStr(num int) string {
	return fmt.Sprintf("%d", num)
}

func StrInt(num string) int {
	v, err := strconv.ParseInt(num, 10, 0)
	if err != nil {
		return 0
	}
	newNum := int(v)
	return newNum
}

func IntStrPtr(num string) *int {
	v, err := strconv.ParseInt(num, 10, 0)
	if err != nil {
		return nil
	}
	newNum := int(v)
	return &newNum
}

func FormatDate(date time.Time) string {
	return fmt.Sprintf("%02d.%02d.%d", date.Day(), date.Month(), date.Year())
}

func RandomDate() time.Time {
	startDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC) // Начальная дата
	endDate := time.Now()
	delta := endDate.Sub(startDate)
	randomDuration := time.Duration(rand.Int63n(int64(delta)))
	return startDate.Add(randomDuration)
}
