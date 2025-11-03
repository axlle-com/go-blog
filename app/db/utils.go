package db

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/axlle-com/blog/app/logger"
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
	if s == "" {
		return nil
	}
	return &s
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

func IntToUintPtr(u int) *uint {
	v := uint(u)
	return &v
}

func UintPtr(u uint) *uint {
	if u == 0 {
		return nil
	}
	return &u
}

func IntPtr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

func Int8Ptr(i int8) *int8 {
	if i == 0 {
		return nil
	}
	return &i
}

func Int16Ptr(i int16) *int16 {
	if i == 0 {
		return nil
	}
	return &i
}

func Int32Ptr(i int32) *int32 {
	if i == 0 {
		return nil
	}
	return &i
}

func Int64Ptr(i int64) *int64 {
	if i == 0 {
		return nil
	}
	return &i
}

func Float32Ptr(i float32) *float32 {
	if i == 0 {
		return nil
	}
	return &i
}

func Float64Ptr(i float64) *float64 {
	if i == 0 {
		return nil
	}
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

// IndexName формирует единый формат имени индекса: idx_<table>_<col1>[_<col2>...]
func IndexName(table string, columns ...string) string {
	parts := make([]string, 0, 1+len(columns))
	clean := func(s string) string {
		s = strings.TrimSpace(s)
		s = strings.ReplaceAll(s, " ", "_")
		s = strings.ReplaceAll(s, ".", "_")
		s = strings.ReplaceAll(s, ",", "_")
		return s
	}
	parts = append(parts, clean(table))
	for _, c := range columns {
		parts = append(parts, clean(c))
	}
	return "idx_" + strings.Join(parts, "_")
}

// HashIndex создаёт HASH индекс по колонке: CREATE INDEX IF NOT EXISTS idx_<table>_<col> ON <table> USING hash (<col>);
// Важно: HASH индексы в PostgreSQL применяются к одной колонке. Для нескольких колонок создавайте отдельные индексы.
func HashIndex(table string, column string) string {
	name := IndexName(table, column)
	return fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s USING hash (%s);", name, table, column)
}

func UniqueIndex(table string, columns ...string) string {
	name := IndexName(table, columns...)
	columnsStr := strings.Join(columns, ", ")
	return fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS %s ON %s (%s);", name, table, columnsStr)
}

func CompositeIndex(table string, columns ...string) string {
	name := IndexName(table, columns...)
	columnsStr := strings.Join(columns, ", ")
	return fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s (%s);", name, table, columnsStr)
}

// GinIndex создаёт GIN индекс для JSONB полей: CREATE INDEX IF NOT EXISTS idx_<table>_<col> ON <table> USING gin(<col> jsonb_path_ops);
func GinIndex(table string, column string) string {
	name := IndexName(table, column)
	return fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s USING gin(%s jsonb_path_ops);", name, table, column)
}
