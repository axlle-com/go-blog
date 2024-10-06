package service

import (
	"github.com/axlle-com/blog/pkg/common/db"
	"strconv"
	"testing"
	"time"
)

type formStruct struct {
	Uint       string
	UintPtr    string
	Int        string
	IntPtr     string
	Int8       string
	Int8Ptr    string
	Int16      string
	Int16Ptr   string
	Int32      string
	Int32Ptr   string
	Int64      string
	Int64Ptr   string
	String     string
	StringPtr  string
	Float32    string
	Float32Ptr string
	Float64    string
	Float64Ptr string
	Bool       string
	BoolPtr    string
	Time       string
	TimePtr    string
}

type sourceStruct struct {
	Uint       uint
	UintPtr    uint
	Int        int
	IntPtr     int
	Int8       int8
	Int8Ptr    int8
	Int16      int16
	Int16Ptr   int16
	Int32      int32
	Int32Ptr   int32
	Int64      int64
	Int64Ptr   int64
	String     string
	StringPtr  string
	Float32    float32
	Float32Ptr float32
	Float64    float64
	Float64Ptr float64
	Bool       bool
	BoolPtr    bool
	Time       time.Time
	TimePtr    time.Time
}

type destinationStruct struct {
	Uint       uint
	UintPtr    *uint
	Int        int
	IntPtr     *int
	Int8       int8
	Int8Ptr    *int8
	Int16      int16
	Int16Ptr   *int16
	Int32      int32
	Int32Ptr   *int32
	Int64      int64
	Int64Ptr   *int64
	String     string
	StringPtr  *string
	Float32    float32
	Float32Ptr *float32
	Float64    float64
	Float64Ptr *float64
	Bool       bool
	BoolPtr    *bool
	Time       time.Time
	TimePtr    *time.Time
}

func TestLoadFromRequestString(t *testing.T) {
	dest := getDestination()
	form := getForm()

	expectedUint, _ := strconv.ParseUint(form.Uint, 10, 64)
	expectedUintPtr, _ := strconv.ParseUint(form.UintPtr, 10, 64)
	expectedInt, _ := strconv.Atoi(form.Int)
	expectedIntPtr, _ := strconv.Atoi(form.IntPtr)
	expectedInt8, _ := strconv.ParseInt(form.Int8, 10, 8)
	expectedInt8Ptr, _ := strconv.ParseInt(form.Int8Ptr, 10, 8)
	expectedInt16, _ := strconv.ParseInt(form.Int16, 10, 16)
	expectedInt16Ptr, _ := strconv.ParseInt(form.Int16Ptr, 10, 16)
	expectedInt32, _ := strconv.ParseInt(form.Int32, 10, 32)
	expectedInt32Ptr, _ := strconv.ParseInt(form.Int32Ptr, 10, 32)
	expectedInt64, _ := strconv.ParseInt(form.Int64, 10, 64)
	expectedInt64Ptr, _ := strconv.ParseInt(form.Int64Ptr, 10, 64)
	expectedFloat32, _ := strconv.ParseFloat(form.Float32, 32)
	expectedFloat32Ptr, _ := strconv.ParseFloat(form.Float32Ptr, 32)
	expectedFloat64, _ := strconv.ParseFloat(form.Float64, 64)
	expectedFloat64Ptr, _ := strconv.ParseFloat(form.Float64Ptr, 64)
	expectedBool, _ := strconv.ParseBool(form.Bool)
	expectedBoolPtr, _ := strconv.ParseBool(form.BoolPtr)
	expectedTime, _ := time.Parse("2006-01-02 15:04:05", form.Time)
	expectedTimePtr, _ := time.Parse("2006-01-02 15:04:05", form.TimePtr)

	LoadFromRequest(&dest, &form)

	if dest.Uint != uint(expectedUint) {
		t.Errorf("Ожидалось Uint: %d, получено: %d", expectedUint, dest.Uint)
	}
	// TODO почему тест проходил?
	if dest.UintPtr == nil || *dest.UintPtr != uint(expectedUintPtr) {
		t.Errorf("Ожидалось UintPtr: %d, получено: %v", expectedUintPtr, *dest.UintPtr)
	}
	if dest.Int != expectedInt {
		t.Errorf("Ожидалось Int: %d, получено: %d", expectedInt, dest.Int)
	}
	if dest.IntPtr == nil || *dest.IntPtr != expectedIntPtr {
		t.Errorf("Ожидалось IntPtr: %d, получено: %v", expectedIntPtr, *dest.IntPtr)
	}
	if dest.Int8 != int8(expectedInt8) {
		t.Errorf("Ожидалось Int8: %d, получено: %d", expectedInt8, dest.Int8)
	}
	if dest.Int8Ptr == nil || *dest.Int8Ptr != int8(expectedInt8Ptr) {
		t.Errorf("Ожидалось Int8Ptr: %d, получено: %v", expectedInt8Ptr, *dest.Int8Ptr)
	}
	if dest.Int16 != int16(expectedInt16) {
		t.Errorf("Ожидалось Int16: %d, получено: %d", expectedInt16, dest.Int16)
	}
	if dest.Int16Ptr == nil || *dest.Int16Ptr != int16(expectedInt16Ptr) {
		t.Errorf("Ожидалось Int16Ptr: %d, получено: %v", expectedInt16Ptr, *dest.Int16Ptr)
	}
	if dest.Int32 != int32(expectedInt32) {
		t.Errorf("Ожидалось Int32: %d, получено: %d", expectedInt32, dest.Int32)
	}
	if dest.Int32Ptr == nil || *dest.Int32Ptr != int32(expectedInt32Ptr) {
		t.Errorf("Ожидалось Int32Ptr: %d, получено: %v", expectedInt32Ptr, *dest.Int32Ptr)
	}
	if dest.Int64 != expectedInt64 {
		t.Errorf("Ожидалось Int64: %d, получено: %d", expectedInt64, dest.Int64)
	}
	if dest.Int64Ptr == nil || *dest.Int64Ptr != expectedInt64Ptr {
		t.Errorf("Ожидалось Int64Ptr: %d, получено: %v", expectedInt64Ptr, *dest.Int64Ptr)
	}
	if dest.String != form.String {
		t.Errorf("Ожидалось String: %s, получено: %s", form.String, dest.String)
	}
	if dest.StringPtr == nil || *dest.StringPtr != form.StringPtr {
		t.Errorf("Ожидалось StringPtr: %s, получено: %v", form.StringPtr, *dest.StringPtr)
	}
	if dest.Float32 != float32(expectedFloat32) {
		t.Errorf("Ожидалось Float32: %f, получено: %f", expectedFloat32, dest.Float32)
	}
	if dest.Float32Ptr == nil || *dest.Float32Ptr != float32(expectedFloat32Ptr) {
		t.Errorf("Ожидалось Float32Ptr: %f, получено: %v", expectedFloat32Ptr, *dest.Float32Ptr)
	}
	if dest.Float64 != expectedFloat64 {
		t.Errorf("Ожидалось Float64: %f, получено: %f", expectedFloat64, dest.Float64)
	}
	if dest.Float64Ptr == nil || *dest.Float64Ptr != expectedFloat64Ptr {
		t.Errorf("Ожидалось Float64Ptr: %f, получено: %v", expectedFloat64Ptr, *dest.Float64Ptr)
	}
	if dest.Bool != expectedBool {
		t.Errorf("Ожидалось Bool: %t, получено: %t", expectedBool, dest.Bool)
	}
	if dest.BoolPtr == nil || *dest.BoolPtr != expectedBoolPtr {
		t.Errorf("Ожидалось BoolPtr: %t, получено: %v", expectedBoolPtr, *dest.BoolPtr)
	}
	if !dest.Time.Equal(expectedTime) {
		t.Errorf("Ожидалось Time: %v, получено: %v", expectedTime, dest.Time)
	}
	if dest.TimePtr == nil || !dest.TimePtr.Equal(expectedTimePtr) {
		t.Errorf("Ожидалось TimePtr: %v, получено: %v", expectedTimePtr, *dest.TimePtr)
	}
}

func TestLoadFromRequestStringEmpty(t *testing.T) {
	dest := getDestination()
	form := formStruct{}

	LoadFromRequest(&dest, &form)
	if dest.Uint != 0 {
		t.Errorf("Ожидалось Uint: %d, получено: %d", 0, dest.Uint)
	}
	if dest.UintPtr != nil {
		t.Errorf("Ожидалось UintPtr: %v, получено: %v", nil, dest.UintPtr)
	}
	if dest.Int != 0 {
		t.Errorf("Ожидалось Int: %d, получено: %d", 0, dest.Int)
	}
	if dest.IntPtr != nil {
		t.Errorf("Ожидалось IntPtr: %v, получено: %v", nil, dest.IntPtr)
	}
	if dest.Int8 != 0 {
		t.Errorf("Ожидалось Int8: %d, получено: %d", 0, dest.Int8)
	}
	if dest.Int8Ptr != nil {
		t.Errorf("Ожидалось Int8Ptr: %v, получено: %v", nil, dest.Int8Ptr)
	}
	if dest.Int16 != 0 {
		t.Errorf("Ожидалось Int16: %d, получено: %d", 0, dest.Int16)
	}
	if dest.Int16Ptr != nil {
		t.Errorf("Ожидалось Int16Ptr: %v, получено: %v", nil, dest.Int16Ptr)
	}
	if dest.Int32 != 0 {
		t.Errorf("Ожидалось Int32: %d, получено: %d", 0, dest.Int32)
	}
	if dest.Int32Ptr != nil {
		t.Errorf("Ожидалось Int32Ptr: %v, получено: %v", nil, dest.Int32Ptr)
	}
	if dest.Int64 != 0 {
		t.Errorf("Ожидалось Int64: %d, получено: %d", 0, dest.Int64)
	}
	if dest.Int64Ptr != nil {
		t.Errorf("Ожидалось Int64Ptr: %v, получено: %v", nil, dest.Int64Ptr)
	}
	if dest.String != "" {
		t.Errorf("Ожидалось String: %s, получено: %s", "", dest.String)
	}
	if dest.StringPtr != nil {
		t.Errorf("Ожидалось StringPtr: %v, получено: %v", nil, dest.StringPtr)
	}
	if dest.Float32 != float32(0) {
		t.Errorf("Ожидалось Float32: %f, получено: %f", float32(0), dest.Float32)
	}
	if dest.Float32Ptr != nil {
		t.Errorf("Ожидалось Float32Ptr: %v, получено: %v", nil, dest.Float32Ptr)
	}
	if dest.Float64 != float64(0) {
		t.Errorf("Ожидалось Float64: %f, получено: %f", float64(0), dest.Float64)
	}
	if dest.Float64Ptr != nil {
		t.Errorf("Ожидалось Float64Ptr: %v, получено: %v", nil, dest.Float64Ptr)
	}
	if dest.Bool != false {
		t.Errorf("Ожидалось Bool: %t, получено: %t", false, dest.Bool)
	}
	if dest.BoolPtr != nil {
		t.Errorf("Ожидалось BoolPtr: %v, получено: %v", nil, dest.BoolPtr)
	}
	if !dest.Time.Equal(time.Time{}) {
		t.Errorf("Ожидалось Time: %v, получено: %v", time.Time{}, dest.Time)
	}
	if dest.TimePtr != nil {
		t.Errorf("Ожидалось TimePtr: %v, получено: %v", nil, dest.TimePtr)
	}
}

func TestLoadFromRequestType(t *testing.T) {
	dest := getDestination()
	src := getSource()

	LoadFromRequest(&dest, &src)

	if dest.Uint != src.Uint {
		t.Errorf("Ожидалось Uint: %d, получено: %d", src.Uint, dest.Uint)
	}
	if dest.UintPtr == nil || *dest.UintPtr != src.UintPtr {
		t.Errorf("Ожидалось UintPtr: %d, получено: %v", src.UintPtr, *dest.UintPtr)
	}
	if dest.Int != src.Int {
		t.Errorf("Ожидалось Int: %d, получено: %d", src.Int, dest.Int)
	}
	if dest.IntPtr == nil || *dest.IntPtr != src.IntPtr {
		t.Errorf("Ожидалось IntPtr: %d, получено: %v", src.IntPtr, *dest.IntPtr)
	}
	if dest.Int8 != src.Int8 {
		t.Errorf("Ожидалось Int8: %d, получено: %d", src.Int8, dest.Int8)
	}
	if dest.Int8Ptr == nil || *dest.Int8Ptr != src.Int8Ptr {
		t.Errorf("Ожидалось Int8Ptr: %d, получено: %v", src.Int8Ptr, *dest.Int8Ptr)
	}
	if dest.Int16 != src.Int16 {
		t.Errorf("Ожидалось Int16: %d, получено: %d", src.Int16, dest.Int16)
	}
	if dest.Int16Ptr == nil || *dest.Int16Ptr != src.Int16Ptr {
		t.Errorf("Ожидалось Int16Ptr: %d, получено: %v", src.Int16Ptr, *dest.Int16Ptr)
	}
	if dest.Int32 != src.Int32 {
		t.Errorf("Ожидалось Int32: %d, получено: %d", src.Int32, dest.Int32)
	}
	if dest.Int32Ptr == nil || *dest.Int32Ptr != src.Int32Ptr {
		t.Errorf("Ожидалось Int32Ptr: %d, получено: %v", src.Int32Ptr, *dest.Int32Ptr)
	}
	if dest.Int64 != src.Int64 {
		t.Errorf("Ожидалось Int64: %d, получено: %d", src.Int64, dest.Int64)
	}
	if dest.Int64Ptr == nil || *dest.Int64Ptr != src.Int64Ptr {
		t.Errorf("Ожидалось Int64Ptr: %d, получено: %v", src.Int64Ptr, *dest.Int64Ptr)
	}
	if dest.String != src.String {
		t.Errorf("Ожидалось String: %s, получено: %s", src.String, dest.String)
	}
	if dest.StringPtr == nil || *dest.StringPtr != src.StringPtr {
		t.Errorf("Ожидалось StringPtr: %s, получено: %v", src.StringPtr, *dest.StringPtr)
	}
	if dest.Float32 != src.Float32 {
		t.Errorf("Ожидалось Float32: %f, получено: %f", src.Float32, dest.Float32)
	}
	if dest.Float32Ptr == nil || *dest.Float32Ptr != src.Float32Ptr {
		t.Errorf("Ожидалось Float32Ptr: %f, получено: %v", src.Float32Ptr, *dest.Float32Ptr)
	}
	if dest.Float64 != src.Float64 {
		t.Errorf("Ожидалось Float64: %f, получено: %f", src.Float64, dest.Float64)
	}
	if dest.Float64Ptr == nil || *dest.Float64Ptr != src.Float64Ptr {
		t.Errorf("Ожидалось Float64Ptr: %f, получено: %v", src.Float64Ptr, *dest.Float64Ptr)
	}
	if dest.Bool != src.Bool {
		t.Errorf("Ожидалось Bool: %t, получено: %t", src.Bool, dest.Bool)
	}
	if dest.BoolPtr == nil || *dest.BoolPtr != src.BoolPtr {
		t.Errorf("Ожидалось BoolPtr: %t, получено: %v", src.BoolPtr, *dest.BoolPtr)
	}
	if !dest.Time.Equal(src.Time) {
		t.Errorf("Ожидалось Time: %v, получено: %v", src.Time, dest.Time)
	}
	if dest.TimePtr == nil || !dest.TimePtr.Equal(src.TimePtr) {
		t.Errorf("Ожидалось TimePtr: %v, получено: %v", src.TimePtr, *dest.TimePtr)
	}
}

func TestLoadFromRequestTypeEmpty(t *testing.T) {
	dest := getDestination()
	src := sourceStruct{}

	LoadFromRequest(&dest, &src)

	if dest.Uint != src.Uint {
		t.Errorf("Ожидалось Uint: %d, получено: %d", src.Uint, dest.Uint)
	}
	if dest.UintPtr != nil {
		t.Errorf("Ожидалось UintPtr: %v, получено: %v", nil, dest.UintPtr)
	}
	if dest.Int != src.Int {
		t.Errorf("Ожидалось Int: %d, получено: %d", src.Int, dest.Int)
	}
	if dest.IntPtr != nil {
		t.Errorf("Ожидалось IntPtr: %v, получено: %v", nil, dest.IntPtr)
	}
	if dest.Int8 != src.Int8 {
		t.Errorf("Ожидалось Int8: %d, получено: %d", src.Int8, dest.Int8)
	}
	if dest.Int8Ptr != nil {
		t.Errorf("Ожидалось Int8Ptr: %v, получено: %v", nil, dest.Int8Ptr)
	}
	if dest.Int16 != src.Int16 {
		t.Errorf("Ожидалось Int16: %d, получено: %d", src.Int16, dest.Int16)
	}
	if dest.Int16Ptr != nil {
		t.Errorf("Ожидалось Int16Ptr: %v, получено: %v", nil, dest.Int16Ptr)
	}
	if dest.Int32 != src.Int32 {
		t.Errorf("Ожидалось Int32: %d, получено: %d", src.Int32, dest.Int32)
	}
	if dest.Int32Ptr != nil {
		t.Errorf("Ожидалось Int32Ptr: %v, получено: %v", nil, dest.Int32Ptr)
	}
	if dest.Int64 != src.Int64 {
		t.Errorf("Ожидалось Int64: %d, получено: %d", src.Int64, dest.Int64)
	}
	if dest.Int64Ptr != nil {
		t.Errorf("Ожидалось Int64Ptr: %v, получено: %v", nil, dest.Int64Ptr)
	}
	if dest.String != src.String {
		t.Errorf("Ожидалось String: %s, получено: %s", src.String, dest.String)
	}
	if dest.StringPtr != nil {
		t.Errorf("Ожидалось StringPtr: %v, получено: %v", nil, dest.StringPtr)
	}
	if dest.Float32 != src.Float32 {
		t.Errorf("Ожидалось Float32: %f, получено: %f", src.Float32, dest.Float32)
	}
	if dest.Float32Ptr != nil {
		t.Errorf("Ожидалось Float32Ptr: %v, получено: %v", nil, dest.Float32Ptr)
	}
	if dest.Float64 != src.Float64 {
		t.Errorf("Ожидалось Float64: %f, получено: %f", src.Float64, dest.Float64)
	}
	if dest.Float64Ptr != nil {
		t.Errorf("Ожидалось Float64Ptr: %v, получено: %v", nil, dest.Float64Ptr)
	}
	if dest.Bool != src.Bool {
		t.Errorf("Ожидалось Bool: %t, получено: %t", src.Bool, dest.Bool)
	}
	if dest.BoolPtr != nil {
		t.Errorf("Ожидалось BoolPtr: %v, получено: %v", nil, dest.BoolPtr)
	}
	if !dest.Time.Equal(src.Time) {
		t.Errorf("Ожидалось Time: %v, получено: %v", src.Time, dest.Time)
	}
	if dest.TimePtr != nil {
		t.Errorf("Ожидалось TimePtr: %v, получено: %v", nil, dest.TimePtr)
	}
}

func getForm() formStruct {
	return formStruct{
		Uint:       "10",
		UintPtr:    "10",
		Int:        "20",
		IntPtr:     "20",
		Int8:       "8",
		Int8Ptr:    "8",
		Int16:      "16",
		Int16Ptr:   "16",
		Int32:      "32",
		Int32Ptr:   "32",
		Int64:      "64",
		Int64Ptr:   "64",
		String:     "String",
		StringPtr:  "String",
		Float32:    "32.32",
		Float32Ptr: "32.32",
		Float64:    "64.64",
		Float64Ptr: "64.64",
		Bool:       "true",
		BoolPtr:    "true",
		Time:       "2000-01-01 00:00:00",
		TimePtr:    "1990-01-01 00:00:00",
	}
}

func getSource() sourceStruct {
	return sourceStruct{
		Uint:       10,
		UintPtr:    10,
		Int:        20,
		IntPtr:     20,
		Int8:       8,
		Int8Ptr:    8,
		Int16:      16,
		Int16Ptr:   16,
		Int32:      32,
		Int32Ptr:   32,
		Int64:      64,
		Int64Ptr:   64,
		String:     "String",
		StringPtr:  "StringPtr",
		Float32:    32.32,
		Float32Ptr: 32.32,
		Float64:    64.64,
		Float64Ptr: 64.64,
		Bool:       true,
		BoolPtr:    true,
		Time:       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		TimePtr:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func getDestination() destinationStruct {
	return destinationStruct{
		Uint:       100,
		UintPtr:    db.UintPtr(100),
		Int:        200,
		IntPtr:     db.IntPtr(200),
		Int8:       80,
		Int8Ptr:    db.Int8Ptr(80),
		Int16:      160,
		Int16Ptr:   db.Int16Ptr(160),
		Int32:      320,
		Int32Ptr:   db.Int32Ptr(320),
		Int64:      640,
		Int64Ptr:   db.Int64Ptr(640),
		String:     "String",
		StringPtr:  db.StrPtr("StringPtr"),
		Float32:    320.32,
		Float32Ptr: db.Float32Ptr(320.32),
		Float64:    640.64,
		Float64Ptr: db.Float64Ptr(640.64),
		Bool:       true,
		BoolPtr:    db.BoolToBoolPtr(true),
		Time:       time.Time{},
		TimePtr:    &time.Time{},
	}
}
