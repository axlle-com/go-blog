package db

import (
	"log"
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

func ParseDate(dateStr string) *time.Time {
	layout := "02.01.2006"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		log.Println("Error parsing date:", err)
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

func IDStrPtr(id string) *uint {
	v, err := strconv.Atoi(id)
	if err != nil {
		return nil
	}
	newID := uint(v)
	return &newID
}

func CheckStr(num string) bool {
	return num == "1"
}

func IntStr(num string) int {
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
