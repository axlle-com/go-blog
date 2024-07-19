package db

import (
	"fmt"
	"math/rand"
	"time"
)

func IntToBoolPtr() *bool {
	value := rand.Intn(2)
	result := value != 0
	return &result
}

func ParseDate(dateStr string) *time.Time {
	layout := "2006-01-02"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
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
	v := uint(rand.Intn(u))
	return &v
}

func IntPtr(i int) *int {
	return &i
}

func Float32Ptr(f float32) *float32 {
	return &f
}
