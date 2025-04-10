package service

import (
	"github.com/axlle-com/blog/app/db"
	"strconv"
	"testing"
	"time"
)

// TODO протестировать формат времени
type FormStruct struct {
	Uint                   string
	UintPtr                string
	Int                    string
	IntPtr                 string
	Int8                   string
	Int8Ptr                string
	Int16                  string
	Int16Ptr               string
	Int32                  string
	Int32Ptr               string
	Int64                  string
	Int64Ptr               string
	String                 string
	StringPtr              string
	Float32                string
	Float32Ptr             string
	Float64                string
	Float64Ptr             string
	Bool                   string
	BoolPtr                string
	Time                   string
	TimePtr                string
	Fake                   string
	Inner                  FormStructInner
	InnerPtr               *FormStructInner
	InnerPtrFake           FormStructInner
	Slice                  []string
	SliceString            []string
	SliceIntPtr            []string
	SliceStructInner       []FormStructInner
	SliceStructInnerTwo    []FormStructInner
	SliceStructInnerPtr    []FormStructInner
	SliceStructInnerPtrTwo []*FormStructInner
}

type FormStructInner struct {
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
	Fake       string
}

type SourceStruct struct {
	Uint        uint
	UintPtr     uint
	Int         int
	IntPtr      int
	Int8        int8
	Int8Ptr     int8
	Int16       int16
	Int16Ptr    int16
	Int32       int32
	Int32Ptr    int32
	Int64       int64
	Int64Ptr    int64
	String      string
	StringPtr   string
	Float32     float32
	Float32Ptr  float32
	Float64     float64
	Float64Ptr  float64
	Bool        bool
	BoolPtr     bool
	Time        time.Time
	TimePtr     time.Time
	Fake        string
	Inner       SourceStructInner
	InnerPtr    *SourceStructInner
	Slice       []int
	SliceString []string
	SliceIntPtr []*int
}

type SourceStructInner struct {
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
	Fake       string
}

type DestinationStruct struct {
	Uint                   uint
	UintPtr                *uint
	Int                    int
	IntPtr                 *int
	Int8                   int8
	Int8Ptr                *int8
	Int16                  int16
	Int16Ptr               *int16
	Int32                  int32
	Int32Ptr               *int32
	Int64                  int64
	Int64Ptr               *int64
	String                 string
	StringPtr              *string
	Float32                float32
	Float32Ptr             *float32
	Float64                float64
	Float64Ptr             *float64
	Bool                   bool
	BoolPtr                *bool
	Time                   time.Time
	TimePtr                *time.Time //TODO перепроверить
	FakeDest               string
	Inner                  DestinationStructInner
	InnerPtr               *DestinationStructInner
	InnerPtrFake           *DestinationStructInner
	Slice                  []int
	SliceString            []string
	SliceIntPtr            []*int
	SliceStructInner       []DestinationStructInner
	SliceStructInnerTwo    []DestinationStructInner
	SliceStructInnerPtr    []*DestinationStructInner
	SliceStructInnerPtrTwo []DestinationStructInner
}

type DestinationStructInner struct {
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
	FakeDest   string
	Time       time.Time
	TimePtr    *time.Time
}

func TestLoadFromRequestString(t *testing.T) {
	dest := getDestination()
	form := getForm()

	LoadStruct(&dest, &form)
	checkString(form, dest, t)
	checkStringInner(form, dest, t)
	checkStringInnerPtr(form, dest, t)
	checkStringInnerPtrFake(form, dest, t)
	checkStringSliceInner(form, dest, t)
	checkStringSliceInnerTwo(form, dest, t)
	checkStringSliceInnerPtrTwo(form, dest, t)
}

func TestLoadFromRequestStringEmpty(t *testing.T) {
	dest := getDestination()

	LoadStruct(&dest, &FormStruct{})
	checkStringEmpty(dest, t)
	checkStringInnerEmpty(dest, t)
	checkStringInnerPtrEmpty(dest, t)
}

func TestLoadFromRequestType(t *testing.T) {
	dest := getDestination()
	src := getSource()

	LoadStruct(&dest, &src)
	checkType(src, dest, t)
	checkTypeInner(src.Inner, dest.Inner, t)
	checkTypeInner(*src.InnerPtr, *dest.InnerPtr, t)
}

func TestLoadFromRequestTypeEmpty(t *testing.T) {
	dest := getDestination()
	src := SourceStruct{}

	LoadStruct(&dest, &src)
	checkTypeEmpty(src, dest, t)
	checkTypeInnerEmpty(src.Inner, dest.Inner, t)

	if dest.InnerPtr != nil {
		t.Errorf("Ожидалось dest.InnerPtr: %v, получено: %v", nil, dest.InnerPtr)
	}
}

func getForm() FormStruct {
	inner1 := getFormInner()
	inner2 := getFormInner()
	inner3 := getFormInner()
	inner4 := getFormInner()
	return FormStruct{
		Uint:                   "10",
		UintPtr:                "10",
		Int:                    "20",
		IntPtr:                 "20",
		Int8:                   "8",
		Int8Ptr:                "8",
		Int16:                  "16",
		Int16Ptr:               "16",
		Int32:                  "32",
		Int32Ptr:               "32",
		Int64:                  "64",
		Int64Ptr:               "64",
		String:                 "String",
		StringPtr:              "String",
		Float32:                "32.32",
		Float32Ptr:             "32.32",
		Float64:                "64.64",
		Float64Ptr:             "64.64",
		Bool:                   "true",
		BoolPtr:                "true",
		Time:                   "2000-01-01 00:00:00",
		TimePtr:                "1990-01-01 00:00:00",
		Fake:                   "string",
		Slice:                  []string{"1", "2", "3", "4"},
		SliceString:            []string{"11", "22", "33", "44"},
		SliceIntPtr:            []string{"11", "22", "33", "44"},
		SliceStructInner:       []FormStructInner{getFormInner(), getFormInner(), getFormInner()},
		SliceStructInnerTwo:    []FormStructInner{getFormInner(), getFormInner(), getFormInner()},
		SliceStructInnerPtr:    []FormStructInner{getFormInner(), getFormInner(), getFormInner()},
		SliceStructInnerPtrTwo: []*FormStructInner{&inner2, &inner3, &inner4},
		Inner:                  getFormInner(),
		InnerPtr:               &inner1,
		InnerPtrFake:           getFormInner(),
	}
}

func getFormInner() FormStructInner {
	return FormStructInner{
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
		Fake:       "string",
	}
}

func getSource() SourceStruct {
	return SourceStruct{
		Uint:        10,
		UintPtr:     10,
		Int:         20,
		IntPtr:      20,
		Int8:        8,
		Int8Ptr:     8,
		Int16:       16,
		Int16Ptr:    16,
		Int32:       32,
		Int32Ptr:    32,
		Int64:       64,
		Int64Ptr:    64,
		String:      "String",
		StringPtr:   "StringPtr",
		Float32:     32.32,
		Float32Ptr:  32.32,
		Float64:     64.64,
		Float64Ptr:  64.64,
		Bool:        true,
		BoolPtr:     true,
		Time:        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		TimePtr:     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Fake:        "string",
		Slice:       []int{1, 2, 3},
		SliceString: []string{"11", "22", "33"},
		SliceIntPtr: []*int{db.IntPtr(11), db.IntPtr(22), db.IntPtr(33)},
		Inner: SourceStructInner{
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
			Fake:       "string",
		},
		InnerPtr: &SourceStructInner{
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
			Fake:       "string",
		},
	}
}

func getDestination() DestinationStruct {
	inner1 := getDestinationInner()
	inner2 := getDestinationInner()
	inner3 := getDestinationInner()
	inner4 := getDestinationInner()
	inner5 := getDestinationInner()
	inner6 := getDestinationInner()
	return DestinationStruct{
		Uint:                   100,
		UintPtr:                db.IntToUintPtr(100),
		Int:                    200,
		IntPtr:                 db.IntPtr(200),
		Int8:                   80,
		Int8Ptr:                db.Int8Ptr(80),
		Int16:                  160,
		Int16Ptr:               db.Int16Ptr(160),
		Int32:                  320,
		Int32Ptr:               db.Int32Ptr(320),
		Int64:                  640,
		Int64Ptr:               db.Int64Ptr(640),
		String:                 "String",
		StringPtr:              db.StrPtr("StringPtr"),
		Float32:                320.32,
		Float32Ptr:             db.Float32Ptr(320.32),
		Float64:                640.64,
		Float64Ptr:             db.Float64Ptr(640.64),
		Bool:                   true,
		BoolPtr:                db.BoolToBoolPtr(true),
		Time:                   time.Time{},
		TimePtr:                &time.Time{},
		FakeDest:               "string",
		Slice:                  []int{4, 5, 6},
		SliceString:            []string{"44", "55", "66"},
		SliceIntPtr:            []*int{db.IntPtr(44), db.IntPtr(55), db.IntPtr(66)},
		SliceStructInner:       []DestinationStructInner{getDestinationInner(), getDestinationInner(), getDestinationInner(), getDestinationInner()},
		SliceStructInnerTwo:    []DestinationStructInner{getDestinationInner(), getDestinationInner()},
		SliceStructInnerPtr:    []*DestinationStructInner{&inner3, &inner4, &inner5, &inner6},
		SliceStructInnerPtrTwo: []DestinationStructInner{getDestinationInner(), getDestinationInner()},
		Inner:                  getDestinationInner(),
		InnerPtr:               &inner1,
		InnerPtrFake:           &inner2,
	}
}

func getDestinationInner() DestinationStructInner {
	return DestinationStructInner{
		Uint:       100,
		UintPtr:    db.IntToUintPtr(100),
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
		FakeDest:   "string",
		Time:       time.Time{},
		TimePtr:    &time.Time{},
	}
}

func checkString(form FormStruct, dest DestinationStruct, t *testing.T) {
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

	if dest.Uint != uint(expectedUint) {
		t.Errorf("Ожидалось Uint: %d, получено: %d", expectedUint, dest.Uint)
	}
	if dest.UintPtr == nil || *dest.UintPtr != uint(expectedUintPtr) {
		t.Errorf("Ожидалось IntToUintPtr: %d, получено: %v", expectedUintPtr, *dest.UintPtr)
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
	if dest.FakeDest != "string" {
		t.Errorf("Ожидалось FakeDest: %s, получено: %v", "string", dest.FakeDest)
	}
	if !dest.Time.Equal(expectedTime) {
		t.Errorf("Ожидалось Time: %v, получено: %v", expectedTime, dest.Time)
	}
	if dest.TimePtr == nil || !dest.TimePtr.Equal(expectedTimePtr) {
		t.Errorf("Ожидалось TimePtr: %v, получено: %v", expectedTimePtr, *dest.TimePtr)
	}

	//Slice
	if dest.Slice == nil || dest.Slice[0] != 1 {
		t.Errorf("Ожидалось Slice[0]: %v, получено: %v", 1, dest.Slice[0])
	}
	if dest.Slice == nil || dest.Slice[1] != 2 {
		t.Errorf("Ожидалось Slice[1]: %v, получено: %v", 2, dest.Slice[1])
	}
	if dest.Slice == nil || dest.Slice[2] != 3 {
		t.Errorf("Ожидалось Slice[2]: %v, получено: %v", 3, dest.Slice[2])
	}

	if dest.Slice == nil || dest.Slice[3] != 4 {
		t.Errorf("Ожидалось Slice[3]: %v, получено: %v", 4, dest.Slice[3])
	}

	// SliceString
	if dest.SliceString == nil || dest.SliceString[0] != "11" {
		t.Errorf("Ожидалось SliceString[0]: %v, получено: %v", 1, dest.SliceString[0])
	}
	if dest.SliceString == nil || dest.SliceString[1] != "22" {
		t.Errorf("Ожидалось SliceString[1]: %v, получено: %v", 2, dest.SliceString[1])
	}
	if dest.SliceString == nil || dest.SliceString[2] != "33" {
		t.Errorf("Ожидалось SliceString[2]: %v, получено: %v", 3, dest.SliceString[2])
	}
	if dest.SliceString == nil || dest.SliceString[3] != "44" {
		t.Errorf("Ожидалось SliceString[3]: %v, получено: %v", 4, dest.SliceString[3])
	}

	// SliceIntPtr
	if dest.SliceIntPtr == nil || *dest.SliceIntPtr[0] != 11 {
		t.Errorf("Ожидалось SliceIntPtr[0]: %v, получено: %v", 11, *dest.SliceIntPtr[0])
	}
	if dest.SliceIntPtr == nil || *dest.SliceIntPtr[1] != 22 {
		t.Errorf("Ожидалось SliceIntPtr[1]: %v, получено: %v", 22, *dest.SliceIntPtr[1])
	}
	if dest.SliceIntPtr == nil || *dest.SliceIntPtr[2] != 33 {
		t.Errorf("Ожидалось SliceIntPtr[2]: %v, получено: %v", 33, *dest.SliceIntPtr[2])
	}
	if dest.SliceIntPtr == nil || *dest.SliceIntPtr[3] != 44 {
		t.Errorf("Ожидалось SliceIntPtr[3]: %v, получено: %v", 44, *dest.SliceIntPtr[3])
	}
}

func checkStringInner(form FormStruct, dest DestinationStruct, t *testing.T) {
	expectedUint, _ := strconv.ParseUint(form.Inner.Uint, 10, 64)
	expectedUintPtr, _ := strconv.ParseUint(form.Inner.UintPtr, 10, 64)
	expectedInt, _ := strconv.Atoi(form.Inner.Int)
	expectedIntPtr, _ := strconv.Atoi(form.Inner.IntPtr)
	expectedInt8, _ := strconv.ParseInt(form.Inner.Int8, 10, 8)
	expectedInt8Ptr, _ := strconv.ParseInt(form.Inner.Int8Ptr, 10, 8)
	expectedInt16, _ := strconv.ParseInt(form.Inner.Int16, 10, 16)
	expectedInt16Ptr, _ := strconv.ParseInt(form.Inner.Int16Ptr, 10, 16)
	expectedInt32, _ := strconv.ParseInt(form.Inner.Int32, 10, 32)
	expectedInt32Ptr, _ := strconv.ParseInt(form.Inner.Int32Ptr, 10, 32)
	expectedInt64, _ := strconv.ParseInt(form.Inner.Int64, 10, 64)
	expectedInt64Ptr, _ := strconv.ParseInt(form.Inner.Int64Ptr, 10, 64)
	expectedFloat32, _ := strconv.ParseFloat(form.Inner.Float32, 32)
	expectedFloat32Ptr, _ := strconv.ParseFloat(form.Inner.Float32Ptr, 32)
	expectedFloat64, _ := strconv.ParseFloat(form.Inner.Float64, 64)
	expectedFloat64Ptr, _ := strconv.ParseFloat(form.Inner.Float64Ptr, 64)
	expectedBool, _ := strconv.ParseBool(form.Inner.Bool)
	expectedBoolPtr, _ := strconv.ParseBool(form.Inner.BoolPtr)
	expectedTime, _ := time.Parse("2006-01-02 15:04:05", form.Inner.Time)
	expectedTimePtr, _ := time.Parse("2006-01-02 15:04:05", form.Inner.TimePtr)

	if dest.Inner.Uint != uint(expectedUint) {
		t.Errorf("Ожидалось Inner.Uint: %v, получено: %v", uint(expectedUint), dest.Inner.Uint)
	}
	if dest.Inner.UintPtr == nil || *dest.Inner.UintPtr != uint(expectedUintPtr) {
		t.Errorf("Ожидалось Inner.IntToUintPtr: %d, получено: %v", expectedUintPtr, *dest.Inner.UintPtr)
	}
	if dest.Inner.Int != expectedInt {
		t.Errorf("Ожидалось Inner.Int: %d, получено: %d", expectedInt, dest.Inner.Int)
	}
	if dest.Inner.IntPtr == nil || *dest.Inner.IntPtr != expectedIntPtr {
		t.Errorf("Ожидалось Inner.IntPtr: %d, получено: %v", expectedIntPtr, *dest.Inner.IntPtr)
	}
	if dest.Inner.Int8 != int8(expectedInt8) {
		t.Errorf("Ожидалось Inner.Int8: %d, получено: %d", expectedInt8, dest.Inner.Int8)
	}
	if dest.Inner.Int8Ptr == nil || *dest.Inner.Int8Ptr != int8(expectedInt8Ptr) {
		t.Errorf("Ожидалось Inner.Int8Ptr: %d, получено: %v", expectedInt8Ptr, *dest.Inner.Int8Ptr)
	}
	if dest.Inner.Int16 != int16(expectedInt16) {
		t.Errorf("Ожидалось Inner.Int16: %d, получено: %d", expectedInt16, dest.Inner.Int16)
	}
	if dest.Inner.Int16Ptr == nil || *dest.Inner.Int16Ptr != int16(expectedInt16Ptr) {
		t.Errorf("Ожидалось Inner.Int16Ptr: %d, получено: %v", expectedInt16Ptr, *dest.Inner.Int16Ptr)
	}
	if dest.Inner.Int32 != int32(expectedInt32) {
		t.Errorf("Ожидалось Inner.Int32: %d, получено: %d", expectedInt32, dest.Inner.Int32)
	}
	if dest.Inner.Int32Ptr == nil || *dest.Inner.Int32Ptr != int32(expectedInt32Ptr) {
		t.Errorf("Ожидалось Inner.Int32Ptr: %d, получено: %v", expectedInt32Ptr, *dest.Inner.Int32Ptr)
	}
	if dest.Inner.Int64 != expectedInt64 {
		t.Errorf("Ожидалось Inner.Int64: %d, получено: %d", expectedInt64, dest.Inner.Int64)
	}
	if dest.Inner.Int64Ptr == nil || *dest.Inner.Int64Ptr != expectedInt64Ptr {
		t.Errorf("Ожидалось Inner.Int64Ptr: %d, получено: %v", expectedInt64Ptr, *dest.Inner.Int64Ptr)
	}
	if dest.Inner.String != form.Inner.String {
		t.Errorf("Ожидалось Inner.String: %s, получено: %s", form.String, dest.Inner.String)
	}
	if dest.Inner.StringPtr == nil || *dest.Inner.StringPtr != form.Inner.StringPtr {
		t.Errorf("Ожидалось Inner.StringPtr: %s, получено: %v", form.StringPtr, *dest.Inner.StringPtr)
	}
	if dest.Inner.Float32 != float32(expectedFloat32) {
		t.Errorf("Ожидалось Inner.Float32: %f, получено: %f", expectedFloat32, dest.Inner.Float32)
	}
	if dest.Inner.Float32Ptr == nil || *dest.Inner.Float32Ptr != float32(expectedFloat32Ptr) {
		t.Errorf("Ожидалось Inner.Float32Ptr: %f, получено: %v", expectedFloat32Ptr, *dest.Inner.Float32Ptr)
	}
	if dest.Inner.Float64 != expectedFloat64 {
		t.Errorf("Ожидалось Inner.Float64: %f, получено: %f", expectedFloat64, dest.Inner.Float64)
	}
	if dest.Inner.Float64Ptr == nil || *dest.Inner.Float64Ptr != expectedFloat64Ptr {
		t.Errorf("Ожидалось Inner.Float64Ptr: %f, получено: %v", expectedFloat64Ptr, *dest.Inner.Float64Ptr)
	}
	if dest.Inner.Bool != expectedBool {
		t.Errorf("Ожидалось Inner.Bool: %t, получено: %t", expectedBool, dest.Inner.Bool)
	}
	if dest.Inner.BoolPtr == nil || *dest.Inner.BoolPtr != expectedBoolPtr {
		t.Errorf("Ожидалось Inner.BoolPtr: %t, получено: %v", expectedBoolPtr, *dest.Inner.BoolPtr)
	}
	if !dest.Inner.Time.Equal(expectedTime) {
		t.Errorf("Ожидалось Inner.Time: %v, получено: %v", expectedTime, dest.Inner.Time)
	}
	if dest.Inner.TimePtr == nil || !dest.Inner.TimePtr.Equal(expectedTimePtr) {
		t.Errorf("Ожидалось Inner.TimePtr: %v, получено: %v", expectedTimePtr, *dest.Inner.TimePtr)
	}
}

func checkStringInnerPtr(form FormStruct, dest DestinationStruct, t *testing.T) {
	expectedUint, _ := strconv.ParseUint(form.InnerPtr.Uint, 10, 64)
	expectedUintPtr, _ := strconv.ParseUint(form.InnerPtr.UintPtr, 10, 64)
	expectedInt, _ := strconv.Atoi(form.InnerPtr.Int)
	expectedIntPtr, _ := strconv.Atoi(form.InnerPtr.IntPtr)
	expectedInt8, _ := strconv.ParseInt(form.InnerPtr.Int8, 10, 8)
	expectedInt8Ptr, _ := strconv.ParseInt(form.InnerPtr.Int8Ptr, 10, 8)
	expectedInt16, _ := strconv.ParseInt(form.InnerPtr.Int16, 10, 16)
	expectedInt16Ptr, _ := strconv.ParseInt(form.InnerPtr.Int16Ptr, 10, 16)
	expectedInt32, _ := strconv.ParseInt(form.InnerPtr.Int32, 10, 32)
	expectedInt32Ptr, _ := strconv.ParseInt(form.InnerPtr.Int32Ptr, 10, 32)
	expectedInt64, _ := strconv.ParseInt(form.InnerPtr.Int64, 10, 64)
	expectedInt64Ptr, _ := strconv.ParseInt(form.InnerPtr.Int64Ptr, 10, 64)
	expectedFloat32, _ := strconv.ParseFloat(form.InnerPtr.Float32, 32)
	expectedFloat32Ptr, _ := strconv.ParseFloat(form.InnerPtr.Float32Ptr, 32)
	expectedFloat64, _ := strconv.ParseFloat(form.InnerPtr.Float64, 64)
	expectedFloat64Ptr, _ := strconv.ParseFloat(form.InnerPtr.Float64Ptr, 64)
	expectedBool, _ := strconv.ParseBool(form.InnerPtr.Bool)
	expectedBoolPtr, _ := strconv.ParseBool(form.InnerPtr.BoolPtr)
	expectedTime, _ := time.Parse("2006-01-02 15:04:05", form.InnerPtr.Time)
	expectedTimePtr, _ := time.Parse("2006-01-02 15:04:05", form.InnerPtr.TimePtr)

	if dest.InnerPtr.Uint != uint(expectedUint) {
		t.Errorf("Ожидалось InnerPtr.Uint: %v, получено: %v", uint(expectedUint), dest.InnerPtr.Uint)
	}
	if dest.InnerPtr.UintPtr == nil || *dest.InnerPtr.UintPtr != uint(expectedUintPtr) {
		t.Errorf("Ожидалось InnerPtr.IntToUintPtr: %d, получено: %v", expectedUintPtr, *dest.InnerPtr.UintPtr)
	}
	if dest.InnerPtr.Int != expectedInt {
		t.Errorf("Ожидалось InnerPtr.Int: %d, получено: %d", expectedInt, dest.InnerPtr.Int)
	}
	if dest.InnerPtr.IntPtr == nil || *dest.InnerPtr.IntPtr != expectedIntPtr {
		t.Errorf("Ожидалось InnerPtr.IntPtr: %d, получено: %v", expectedIntPtr, *dest.InnerPtr.IntPtr)
	}
	if dest.InnerPtr.Int8 != int8(expectedInt8) {
		t.Errorf("Ожидалось InnerPtr.Int8: %d, получено: %d", expectedInt8, dest.InnerPtr.Int8)
	}
	if dest.InnerPtr.Int8Ptr == nil || *dest.InnerPtr.Int8Ptr != int8(expectedInt8Ptr) {
		t.Errorf("Ожидалось InnerPtr.Int8Ptr: %d, получено: %v", expectedInt8Ptr, *dest.InnerPtr.Int8Ptr)
	}
	if dest.InnerPtr.Int16 != int16(expectedInt16) {
		t.Errorf("Ожидалось InnerPtr.Int16: %d, получено: %d", expectedInt16, dest.InnerPtr.Int16)
	}
	if dest.InnerPtr.Int16Ptr == nil || *dest.InnerPtr.Int16Ptr != int16(expectedInt16Ptr) {
		t.Errorf("Ожидалось InnerPtr.Int16Ptr: %d, получено: %v", expectedInt16Ptr, *dest.InnerPtr.Int16Ptr)
	}
	if dest.InnerPtr.Int32 != int32(expectedInt32) {
		t.Errorf("Ожидалось InnerPtr.Int32: %d, получено: %d", expectedInt32, dest.InnerPtr.Int32)
	}
	if dest.InnerPtr.Int32Ptr == nil || *dest.InnerPtr.Int32Ptr != int32(expectedInt32Ptr) {
		t.Errorf("Ожидалось InnerPtr.Int32Ptr: %d, получено: %v", expectedInt32Ptr, *dest.InnerPtr.Int32Ptr)
	}
	if dest.InnerPtr.Int64 != expectedInt64 {
		t.Errorf("Ожидалось InnerPtr.Int64: %d, получено: %d", expectedInt64, dest.InnerPtr.Int64)
	}
	if dest.InnerPtr.Int64Ptr == nil || *dest.InnerPtr.Int64Ptr != expectedInt64Ptr {
		t.Errorf("Ожидалось InnerPtr.Int64Ptr: %d, получено: %v", expectedInt64Ptr, *dest.InnerPtr.Int64Ptr)
	}
	if dest.InnerPtr.String != form.InnerPtr.String {
		t.Errorf("Ожидалось InnerPtr.String: %s, получено: %s", form.String, dest.InnerPtr.String)
	}
	if dest.InnerPtr.StringPtr == nil || *dest.InnerPtr.StringPtr != form.InnerPtr.StringPtr {
		t.Errorf("Ожидалось InnerPtr.StringPtr: %s, получено: %v", form.StringPtr, *dest.InnerPtr.StringPtr)
	}
	if dest.InnerPtr.Float32 != float32(expectedFloat32) {
		t.Errorf("Ожидалось InnerPtr.Float32: %f, получено: %f", expectedFloat32, dest.InnerPtr.Float32)
	}
	if dest.InnerPtr.Float32Ptr == nil || *dest.InnerPtr.Float32Ptr != float32(expectedFloat32Ptr) {
		t.Errorf("Ожидалось InnerPtr.Float32Ptr: %f, получено: %v", expectedFloat32Ptr, *dest.InnerPtr.Float32Ptr)
	}
	if dest.InnerPtr.Float64 != expectedFloat64 {
		t.Errorf("Ожидалось InnerPtr.Float64: %f, получено: %f", expectedFloat64, dest.InnerPtr.Float64)
	}
	if dest.InnerPtr.Float64Ptr == nil || *dest.InnerPtr.Float64Ptr != expectedFloat64Ptr {
		t.Errorf("Ожидалось InnerPtr.Float64Ptr: %f, получено: %v", expectedFloat64Ptr, *dest.InnerPtr.Float64Ptr)
	}
	if dest.InnerPtr.Bool != expectedBool {
		t.Errorf("Ожидалось InnerPtr.Bool: %t, получено: %t", expectedBool, dest.InnerPtr.Bool)
	}
	if dest.InnerPtr.BoolPtr == nil || *dest.InnerPtr.BoolPtr != expectedBoolPtr {
		t.Errorf("Ожидалось InnerPtr.BoolPtr: %t, получено: %v", expectedBoolPtr, *dest.InnerPtr.BoolPtr)
	}
	if !dest.InnerPtr.Time.Equal(expectedTime) {
		t.Errorf("Ожидалось InnerPtr.Time: %v, получено: %v", expectedTime, dest.InnerPtr.Time)
	}
	if dest.InnerPtr.TimePtr == nil || !dest.InnerPtr.TimePtr.Equal(expectedTimePtr) {
		t.Errorf("Ожидалось InnerPtr.TimePtr: %v, получено: %v", expectedTimePtr, *dest.InnerPtr.TimePtr)
	}
}

func checkStringInnerPtrFake(form FormStruct, dest DestinationStruct, t *testing.T) {
	expectedUint, _ := strconv.ParseUint(form.InnerPtrFake.Uint, 10, 64)
	expectedUintPtr, _ := strconv.ParseUint(form.InnerPtrFake.UintPtr, 10, 64)
	expectedInt, _ := strconv.Atoi(form.InnerPtrFake.Int)
	expectedIntPtr, _ := strconv.Atoi(form.InnerPtrFake.IntPtr)
	expectedInt8, _ := strconv.ParseInt(form.InnerPtrFake.Int8, 10, 8)
	expectedInt8Ptr, _ := strconv.ParseInt(form.InnerPtrFake.Int8Ptr, 10, 8)
	expectedInt16, _ := strconv.ParseInt(form.InnerPtrFake.Int16, 10, 16)
	expectedInt16Ptr, _ := strconv.ParseInt(form.InnerPtrFake.Int16Ptr, 10, 16)
	expectedInt32, _ := strconv.ParseInt(form.InnerPtrFake.Int32, 10, 32)
	expectedInt32Ptr, _ := strconv.ParseInt(form.InnerPtrFake.Int32Ptr, 10, 32)
	expectedInt64, _ := strconv.ParseInt(form.InnerPtrFake.Int64, 10, 64)
	expectedInt64Ptr, _ := strconv.ParseInt(form.InnerPtrFake.Int64Ptr, 10, 64)
	expectedFloat32, _ := strconv.ParseFloat(form.InnerPtrFake.Float32, 32)
	expectedFloat32Ptr, _ := strconv.ParseFloat(form.InnerPtrFake.Float32Ptr, 32)
	expectedFloat64, _ := strconv.ParseFloat(form.InnerPtrFake.Float64, 64)
	expectedFloat64Ptr, _ := strconv.ParseFloat(form.InnerPtrFake.Float64Ptr, 64)
	expectedBool, _ := strconv.ParseBool(form.InnerPtrFake.Bool)
	expectedBoolPtr, _ := strconv.ParseBool(form.InnerPtrFake.BoolPtr)
	expectedTime, _ := time.Parse("2006-01-02 15:04:05", form.InnerPtrFake.Time)
	expectedTimePtr, _ := time.Parse("2006-01-02 15:04:05", form.InnerPtrFake.TimePtr)

	if dest.InnerPtrFake.Uint != uint(expectedUint) {
		t.Errorf("Ожидалось InnerPtrFake.Uint: %v, получено: %v", uint(expectedUint), dest.InnerPtrFake.Uint)
	}
	if dest.InnerPtrFake.UintPtr == nil || *dest.InnerPtrFake.UintPtr != uint(expectedUintPtr) {
		t.Errorf("Ожидалось InnerPtrFake.IntToUintPtr: %d, получено: %v", expectedUintPtr, *dest.InnerPtrFake.UintPtr)
	}
	if dest.InnerPtrFake.Int != expectedInt {
		t.Errorf("Ожидалось InnerPtrFake.Int: %d, получено: %d", expectedInt, dest.InnerPtrFake.Int)
	}
	if dest.InnerPtrFake.IntPtr == nil || *dest.InnerPtrFake.IntPtr != expectedIntPtr {
		t.Errorf("Ожидалось InnerPtrFake.IntPtr: %d, получено: %v", expectedIntPtr, *dest.InnerPtrFake.IntPtr)
	}
	if dest.InnerPtrFake.Int8 != int8(expectedInt8) {
		t.Errorf("Ожидалось InnerPtrFake.Int8: %d, получено: %d", expectedInt8, dest.InnerPtrFake.Int8)
	}
	if dest.InnerPtrFake.Int8Ptr == nil || *dest.InnerPtrFake.Int8Ptr != int8(expectedInt8Ptr) {
		t.Errorf("Ожидалось InnerPtrFake.Int8Ptr: %d, получено: %v", expectedInt8Ptr, *dest.InnerPtrFake.Int8Ptr)
	}
	if dest.InnerPtrFake.Int16 != int16(expectedInt16) {
		t.Errorf("Ожидалось InnerPtrFake.Int16: %d, получено: %d", expectedInt16, dest.InnerPtrFake.Int16)
	}
	if dest.InnerPtrFake.Int16Ptr == nil || *dest.InnerPtrFake.Int16Ptr != int16(expectedInt16Ptr) {
		t.Errorf("Ожидалось InnerPtrFake.Int16Ptr: %d, получено: %v", expectedInt16Ptr, *dest.InnerPtrFake.Int16Ptr)
	}
	if dest.InnerPtrFake.Int32 != int32(expectedInt32) {
		t.Errorf("Ожидалось InnerPtrFake.Int32: %d, получено: %d", expectedInt32, dest.InnerPtrFake.Int32)
	}
	if dest.InnerPtrFake.Int32Ptr == nil || *dest.InnerPtrFake.Int32Ptr != int32(expectedInt32Ptr) {
		t.Errorf("Ожидалось InnerPtrFake.Int32Ptr: %d, получено: %v", expectedInt32Ptr, *dest.InnerPtrFake.Int32Ptr)
	}
	if dest.InnerPtrFake.Int64 != expectedInt64 {
		t.Errorf("Ожидалось InnerPtrFake.Int64: %d, получено: %d", expectedInt64, dest.InnerPtrFake.Int64)
	}
	if dest.InnerPtrFake.Int64Ptr == nil || *dest.InnerPtrFake.Int64Ptr != expectedInt64Ptr {
		t.Errorf("Ожидалось InnerPtrFake.Int64Ptr: %d, получено: %v", expectedInt64Ptr, *dest.InnerPtrFake.Int64Ptr)
	}
	if dest.InnerPtrFake.String != form.InnerPtrFake.String {
		t.Errorf("Ожидалось InnerPtrFake.String: %s, получено: %s", form.String, dest.InnerPtrFake.String)
	}
	if dest.InnerPtrFake.StringPtr == nil || *dest.InnerPtrFake.StringPtr != form.InnerPtrFake.StringPtr {
		t.Errorf("Ожидалось InnerPtrFake.StringPtr: %s, получено: %v", form.StringPtr, *dest.InnerPtrFake.StringPtr)
	}
	if dest.InnerPtrFake.Float32 != float32(expectedFloat32) {
		t.Errorf("Ожидалось InnerPtrFake.Float32: %f, получено: %f", expectedFloat32, dest.InnerPtrFake.Float32)
	}
	if dest.InnerPtrFake.Float32Ptr == nil || *dest.InnerPtrFake.Float32Ptr != float32(expectedFloat32Ptr) {
		t.Errorf("Ожидалось InnerPtrFake.Float32Ptr: %f, получено: %v", expectedFloat32Ptr, *dest.InnerPtrFake.Float32Ptr)
	}
	if dest.InnerPtrFake.Float64 != expectedFloat64 {
		t.Errorf("Ожидалось InnerPtrFake.Float64: %f, получено: %f", expectedFloat64, dest.InnerPtrFake.Float64)
	}
	if dest.InnerPtrFake.Float64Ptr == nil || *dest.InnerPtrFake.Float64Ptr != expectedFloat64Ptr {
		t.Errorf("Ожидалось InnerPtrFake.Float64Ptr: %f, получено: %v", expectedFloat64Ptr, *dest.InnerPtrFake.Float64Ptr)
	}
	if dest.InnerPtrFake.Bool != expectedBool {
		t.Errorf("Ожидалось InnerPtrFake.Bool: %t, получено: %t", expectedBool, dest.InnerPtrFake.Bool)
	}
	if dest.InnerPtrFake.BoolPtr == nil || *dest.InnerPtrFake.BoolPtr != expectedBoolPtr {
		t.Errorf("Ожидалось InnerPtrFake.BoolPtr: %t, получено: %v", expectedBoolPtr, *dest.InnerPtrFake.BoolPtr)
	}
	if !dest.InnerPtrFake.Time.Equal(expectedTime) {
		t.Errorf("Ожидалось InnerPtrFake.Time: %v, получено: %v", expectedTime, dest.InnerPtrFake.Time)
	}
	if dest.InnerPtrFake.TimePtr == nil || !dest.InnerPtrFake.TimePtr.Equal(expectedTimePtr) {
		t.Errorf("Ожидалось InnerPtrFake.TimePtr: %v, получено: %v", expectedTimePtr, *dest.InnerPtrFake.TimePtr)
	}
}

func checkStringSliceInner(form FormStruct, dest DestinationStruct, t *testing.T) {
	expectedUint, _ := strconv.ParseUint(form.Inner.Uint, 10, 64)
	expectedUintPtr, _ := strconv.ParseUint(form.Inner.UintPtr, 10, 64)
	expectedInt, _ := strconv.Atoi(form.Inner.Int)
	expectedIntPtr, _ := strconv.Atoi(form.Inner.IntPtr)
	expectedInt8, _ := strconv.ParseInt(form.Inner.Int8, 10, 8)
	expectedInt8Ptr, _ := strconv.ParseInt(form.Inner.Int8Ptr, 10, 8)
	expectedInt16, _ := strconv.ParseInt(form.Inner.Int16, 10, 16)
	expectedInt16Ptr, _ := strconv.ParseInt(form.Inner.Int16Ptr, 10, 16)
	expectedInt32, _ := strconv.ParseInt(form.Inner.Int32, 10, 32)
	expectedInt32Ptr, _ := strconv.ParseInt(form.Inner.Int32Ptr, 10, 32)
	expectedInt64, _ := strconv.ParseInt(form.Inner.Int64, 10, 64)
	expectedInt64Ptr, _ := strconv.ParseInt(form.Inner.Int64Ptr, 10, 64)
	expectedFloat32, _ := strconv.ParseFloat(form.Inner.Float32, 32)
	expectedFloat32Ptr, _ := strconv.ParseFloat(form.Inner.Float32Ptr, 32)
	expectedFloat64, _ := strconv.ParseFloat(form.Inner.Float64, 64)
	expectedFloat64Ptr, _ := strconv.ParseFloat(form.Inner.Float64Ptr, 64)
	expectedBool, _ := strconv.ParseBool(form.Inner.Bool)
	expectedBoolPtr, _ := strconv.ParseBool(form.Inner.BoolPtr)
	expectedTime, _ := time.Parse("2006-01-02 15:04:05", form.Inner.Time)
	expectedTimePtr, _ := time.Parse("2006-01-02 15:04:05", form.Inner.TimePtr)

	if len(form.SliceStructInner) != len(dest.SliceStructInner) {
		t.Errorf("Ожидалось len: %v, получено: %v", len(form.SliceStructInner), len(dest.SliceStructInner))
	}

	for _, strct := range dest.SliceStructInner {
		if strct.Uint != uint(expectedUint) {
			t.Errorf("Ожидалось Inner.Uint: %v, получено: %v", uint(expectedUint), strct.Uint)
		}
		if strct.UintPtr == nil || *strct.UintPtr != uint(expectedUintPtr) {
			t.Errorf("Ожидалось Inner.IntToUintPtr: %d, получено: %v", expectedUintPtr, *strct.UintPtr)
		}
		if strct.Int != expectedInt {
			t.Errorf("Ожидалось Inner.Int: %d, получено: %d", expectedInt, strct.Int)
		}
		if strct.IntPtr == nil || *strct.IntPtr != expectedIntPtr {
			t.Errorf("Ожидалось Inner.IntPtr: %d, получено: %v", expectedIntPtr, *strct.IntPtr)
		}
		if strct.Int8 != int8(expectedInt8) {
			t.Errorf("Ожидалось Inner.Int8: %d, получено: %d", expectedInt8, strct.Int8)
		}
		if strct.Int8Ptr == nil || *strct.Int8Ptr != int8(expectedInt8Ptr) {
			t.Errorf("Ожидалось Inner.Int8Ptr: %d, получено: %v", expectedInt8Ptr, *strct.Int8Ptr)
		}
		if strct.Int16 != int16(expectedInt16) {
			t.Errorf("Ожидалось Inner.Int16: %d, получено: %d", expectedInt16, strct.Int16)
		}
		if strct.Int16Ptr == nil || *strct.Int16Ptr != int16(expectedInt16Ptr) {
			t.Errorf("Ожидалось Inner.Int16Ptr: %d, получено: %v", expectedInt16Ptr, *strct.Int16Ptr)
		}
		if strct.Int32 != int32(expectedInt32) {
			t.Errorf("Ожидалось Inner.Int32: %d, получено: %d", expectedInt32, strct.Int32)
		}
		if strct.Int32Ptr == nil || *strct.Int32Ptr != int32(expectedInt32Ptr) {
			t.Errorf("Ожидалось Inner.Int32Ptr: %d, получено: %v", expectedInt32Ptr, *strct.Int32Ptr)
		}
		if strct.Int64 != expectedInt64 {
			t.Errorf("Ожидалось Inner.Int64: %d, получено: %d", expectedInt64, strct.Int64)
		}
		if strct.Int64Ptr == nil || *strct.Int64Ptr != expectedInt64Ptr {
			t.Errorf("Ожидалось Inner.Int64Ptr: %d, получено: %v", expectedInt64Ptr, *strct.Int64Ptr)
		}
		if strct.String != form.Inner.String {
			t.Errorf("Ожидалось Inner.String: %s, получено: %s", form.String, strct.String)
		}
		if strct.StringPtr == nil || *strct.StringPtr != form.Inner.StringPtr {
			t.Errorf("Ожидалось Inner.StringPtr: %s, получено: %v", form.StringPtr, *strct.StringPtr)
		}
		if strct.Float32 != float32(expectedFloat32) {
			t.Errorf("Ожидалось Inner.Float32: %f, получено: %f", expectedFloat32, strct.Float32)
		}
		if strct.Float32Ptr == nil || *strct.Float32Ptr != float32(expectedFloat32Ptr) {
			t.Errorf("Ожидалось Inner.Float32Ptr: %f, получено: %v", expectedFloat32Ptr, *strct.Float32Ptr)
		}
		if strct.Float64 != expectedFloat64 {
			t.Errorf("Ожидалось Inner.Float64: %f, получено: %f", expectedFloat64, strct.Float64)
		}
		if strct.Float64Ptr == nil || *strct.Float64Ptr != expectedFloat64Ptr {
			t.Errorf("Ожидалось Inner.Float64Ptr: %f, получено: %v", expectedFloat64Ptr, *strct.Float64Ptr)
		}
		if strct.Bool != expectedBool {
			t.Errorf("Ожидалось Inner.Bool: %t, получено: %t", expectedBool, strct.Bool)
		}
		if strct.BoolPtr == nil || *strct.BoolPtr != expectedBoolPtr {
			t.Errorf("Ожидалось Inner.BoolPtr: %t, получено: %v", expectedBoolPtr, *strct.BoolPtr)
		}
		if !strct.Time.Equal(expectedTime) {
			t.Errorf("Ожидалось Inner.Time: %v, получено: %v", expectedTime, strct.Time)
		}
		if strct.TimePtr == nil || !strct.TimePtr.Equal(expectedTimePtr) {
			t.Errorf("Ожидалось Inner.TimePtr: %v, получено: %v", expectedTimePtr, *strct.TimePtr)
		}
	}

}

func checkStringSliceInnerTwo(form FormStruct, dest DestinationStruct, t *testing.T) {
	if len(form.SliceStructInnerTwo) != len(dest.SliceStructInnerTwo) {
		t.Errorf("Ожидалось len: %v, получено: %v", len(form.SliceStructInnerPtr), len(dest.SliceStructInnerPtr))
	}
}

func checkStringSliceInnerPtrTwo(form FormStruct, dest DestinationStruct, t *testing.T) {
	expectedUint, _ := strconv.ParseUint(form.Inner.Uint, 10, 64)
	expectedUintPtr, _ := strconv.ParseUint(form.Inner.UintPtr, 10, 64)
	expectedInt, _ := strconv.Atoi(form.Inner.Int)
	expectedIntPtr, _ := strconv.Atoi(form.Inner.IntPtr)
	expectedInt8, _ := strconv.ParseInt(form.Inner.Int8, 10, 8)
	expectedInt8Ptr, _ := strconv.ParseInt(form.Inner.Int8Ptr, 10, 8)
	expectedInt16, _ := strconv.ParseInt(form.Inner.Int16, 10, 16)
	expectedInt16Ptr, _ := strconv.ParseInt(form.Inner.Int16Ptr, 10, 16)
	expectedInt32, _ := strconv.ParseInt(form.Inner.Int32, 10, 32)
	expectedInt32Ptr, _ := strconv.ParseInt(form.Inner.Int32Ptr, 10, 32)
	expectedInt64, _ := strconv.ParseInt(form.Inner.Int64, 10, 64)
	expectedInt64Ptr, _ := strconv.ParseInt(form.Inner.Int64Ptr, 10, 64)
	expectedFloat32, _ := strconv.ParseFloat(form.Inner.Float32, 32)
	expectedFloat32Ptr, _ := strconv.ParseFloat(form.Inner.Float32Ptr, 32)
	expectedFloat64, _ := strconv.ParseFloat(form.Inner.Float64, 64)
	expectedFloat64Ptr, _ := strconv.ParseFloat(form.Inner.Float64Ptr, 64)
	expectedBool, _ := strconv.ParseBool(form.Inner.Bool)
	expectedBoolPtr, _ := strconv.ParseBool(form.Inner.BoolPtr)
	expectedTime, _ := time.Parse("2006-01-02 15:04:05", form.Inner.Time)
	expectedTimePtr, _ := time.Parse("2006-01-02 15:04:05", form.Inner.TimePtr)

	if len(form.SliceStructInnerPtrTwo) != len(dest.SliceStructInnerPtrTwo) {
		t.Errorf("Ожидалось len: %v, получено: %v", len(form.SliceStructInnerPtrTwo), len(dest.SliceStructInnerPtrTwo))
	}

	for _, strct := range dest.SliceStructInnerPtrTwo {
		if strct.Uint != uint(expectedUint) {
			t.Errorf("Ожидалось Inner.Uint: %v, получено: %v", uint(expectedUint), strct.Uint)
		}
		if strct.UintPtr == nil || *strct.UintPtr != uint(expectedUintPtr) {
			t.Errorf("Ожидалось Inner.IntToUintPtr: %d, получено: %v", expectedUintPtr, *strct.UintPtr)
		}
		if strct.Int != expectedInt {
			t.Errorf("Ожидалось Inner.Int: %d, получено: %d", expectedInt, strct.Int)
		}
		if strct.IntPtr == nil || *strct.IntPtr != expectedIntPtr {
			t.Errorf("Ожидалось Inner.IntPtr: %d, получено: %v", expectedIntPtr, *strct.IntPtr)
		}
		if strct.Int8 != int8(expectedInt8) {
			t.Errorf("Ожидалось Inner.Int8: %d, получено: %d", expectedInt8, strct.Int8)
		}
		if strct.Int8Ptr == nil || *strct.Int8Ptr != int8(expectedInt8Ptr) {
			t.Errorf("Ожидалось Inner.Int8Ptr: %d, получено: %v", expectedInt8Ptr, *strct.Int8Ptr)
		}
		if strct.Int16 != int16(expectedInt16) {
			t.Errorf("Ожидалось Inner.Int16: %d, получено: %d", expectedInt16, strct.Int16)
		}
		if strct.Int16Ptr == nil || *strct.Int16Ptr != int16(expectedInt16Ptr) {
			t.Errorf("Ожидалось Inner.Int16Ptr: %d, получено: %v", expectedInt16Ptr, *strct.Int16Ptr)
		}
		if strct.Int32 != int32(expectedInt32) {
			t.Errorf("Ожидалось Inner.Int32: %d, получено: %d", expectedInt32, strct.Int32)
		}
		if strct.Int32Ptr == nil || *strct.Int32Ptr != int32(expectedInt32Ptr) {
			t.Errorf("Ожидалось Inner.Int32Ptr: %d, получено: %v", expectedInt32Ptr, *strct.Int32Ptr)
		}
		if strct.Int64 != expectedInt64 {
			t.Errorf("Ожидалось Inner.Int64: %d, получено: %d", expectedInt64, strct.Int64)
		}
		if strct.Int64Ptr == nil || *strct.Int64Ptr != expectedInt64Ptr {
			t.Errorf("Ожидалось Inner.Int64Ptr: %d, получено: %v", expectedInt64Ptr, *strct.Int64Ptr)
		}
		if strct.String != form.Inner.String {
			t.Errorf("Ожидалось Inner.String: %s, получено: %s", form.String, strct.String)
		}
		if strct.StringPtr == nil || *strct.StringPtr != form.Inner.StringPtr {
			t.Errorf("Ожидалось Inner.StringPtr: %s, получено: %v", form.StringPtr, *strct.StringPtr)
		}
		if strct.Float32 != float32(expectedFloat32) {
			t.Errorf("Ожидалось Inner.Float32: %f, получено: %f", expectedFloat32, strct.Float32)
		}
		if strct.Float32Ptr == nil || *strct.Float32Ptr != float32(expectedFloat32Ptr) {
			t.Errorf("Ожидалось Inner.Float32Ptr: %f, получено: %v", expectedFloat32Ptr, *strct.Float32Ptr)
		}
		if strct.Float64 != expectedFloat64 {
			t.Errorf("Ожидалось Inner.Float64: %f, получено: %f", expectedFloat64, strct.Float64)
		}
		if strct.Float64Ptr == nil || *strct.Float64Ptr != expectedFloat64Ptr {
			t.Errorf("Ожидалось Inner.Float64Ptr: %f, получено: %v", expectedFloat64Ptr, *strct.Float64Ptr)
		}
		if strct.Bool != expectedBool {
			t.Errorf("Ожидалось Inner.Bool: %t, получено: %t", expectedBool, strct.Bool)
		}
		if strct.BoolPtr == nil || *strct.BoolPtr != expectedBoolPtr {
			t.Errorf("Ожидалось Inner.BoolPtr: %t, получено: %v", expectedBoolPtr, *strct.BoolPtr)
		}
		if !strct.Time.Equal(expectedTime) {
			t.Errorf("Ожидалось Inner.Time: %v, получено: %v", expectedTime, strct.Time)
		}
		if strct.TimePtr == nil || !strct.TimePtr.Equal(expectedTimePtr) {
			t.Errorf("Ожидалось Inner.TimePtr: %v, получено: %v", expectedTimePtr, *strct.TimePtr)
		}
	}

}

func checkStringEmpty(dest DestinationStruct, t *testing.T) {
	if dest.Uint != 0 {
		t.Errorf("Ожидалось Uint: %d, получено: %d", 0, dest.Uint)
	}
	if dest.UintPtr != nil {
		t.Errorf("Ожидалось IntToUintPtr: %v, получено: %v", nil, dest.UintPtr)
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

	// Slice
	if dest.Slice != nil {
		t.Errorf("Ожидалось Slice: %v, получено: %v", nil, dest.Slice)
	}
	if dest.SliceString != nil {
		t.Errorf("Ожидалось SliceString: %v, получено: %v", nil, dest.SliceString)
	}
	if dest.SliceIntPtr != nil {
		t.Errorf("Ожидалось SliceIntPtr: %v, получено: %v", nil, dest.SliceIntPtr)
	}
}

func checkStringInnerEmpty(dest DestinationStruct, t *testing.T) {
	if dest.Inner.Uint != 0 {
		t.Errorf("Ожидалось Uint: %d, получено: %d", 0, dest.Inner.Uint)
	}
	if dest.Inner.UintPtr != nil {
		t.Errorf("Ожидалось IntToUintPtr: %v, получено: %v", nil, dest.Inner.UintPtr)
	}
	if dest.Inner.Int != 0 {
		t.Errorf("Ожидалось Int: %d, получено: %d", 0, dest.Inner.Int)
	}
	if dest.Inner.IntPtr != nil {
		t.Errorf("Ожидалось IntPtr: %v, получено: %v", nil, dest.Inner.IntPtr)
	}
	if dest.Inner.Int8 != 0 {
		t.Errorf("Ожидалось Int8: %d, получено: %d", 0, dest.Inner.Int8)
	}
	if dest.Inner.Int8Ptr != nil {
		t.Errorf("Ожидалось Int8Ptr: %v, получено: %v", nil, dest.Inner.Int8Ptr)
	}
	if dest.Inner.Int16 != 0 {
		t.Errorf("Ожидалось Int16: %d, получено: %d", 0, dest.Inner.Int16)
	}
	if dest.Inner.Int16Ptr != nil {
		t.Errorf("Ожидалось Int16Ptr: %v, получено: %v", nil, dest.Inner.Int16Ptr)
	}
	if dest.Inner.Int32 != 0 {
		t.Errorf("Ожидалось Int32: %d, получено: %d", 0, dest.Inner.Int32)
	}
	if dest.Inner.Int32Ptr != nil {
		t.Errorf("Ожидалось Int32Ptr: %v, получено: %v", nil, dest.Inner.Int32Ptr)
	}
	if dest.Inner.Int64 != 0 {
		t.Errorf("Ожидалось Int64: %d, получено: %d", 0, dest.Inner.Int64)
	}
	if dest.Inner.Int64Ptr != nil {
		t.Errorf("Ожидалось Int64Ptr: %v, получено: %v", nil, dest.Inner.Int64Ptr)
	}
	if dest.Inner.String != "" {
		t.Errorf("Ожидалось String: %s, получено: %s", "", dest.Inner.String)
	}
	if dest.Inner.StringPtr != nil {
		t.Errorf("Ожидалось StringPtr: %v, получено: %v", nil, dest.Inner.StringPtr)
	}
	if dest.Inner.Float32 != float32(0) {
		t.Errorf("Ожидалось Float32: %f, получено: %f", float32(0), dest.Inner.Float32)
	}
	if dest.Inner.Float32Ptr != nil {
		t.Errorf("Ожидалось Float32Ptr: %v, получено: %v", nil, dest.Inner.Float32Ptr)
	}
	if dest.Inner.Float64 != float64(0) {
		t.Errorf("Ожидалось Float64: %f, получено: %f", float64(0), dest.Inner.Float64)
	}
	if dest.Inner.Float64Ptr != nil {
		t.Errorf("Ожидалось Float64Ptr: %v, получено: %v", nil, dest.Inner.Float64Ptr)
	}
	if dest.Inner.Bool != false {
		t.Errorf("Ожидалось Bool: %t, получено: %t", false, dest.Inner.Bool)
	}
	if dest.Inner.BoolPtr != nil {
		t.Errorf("Ожидалось BoolPtr: %v, получено: %v", nil, dest.Inner.BoolPtr)
	}
	if !dest.Inner.Time.Equal(time.Time{}) {
		t.Errorf("Ожидалось Time: %v, получено: %v", time.Time{}, dest.Inner.Time)
	}
	if dest.Inner.TimePtr != nil {
		t.Errorf("Ожидалось TimePtr: %v, получено: %v", nil, dest.Inner.TimePtr)
	}
}

func checkStringInnerPtrEmpty(dest DestinationStruct, t *testing.T) {
	if dest.InnerPtr != nil {
		t.Errorf("Ожидалось dest.InnerPtr: %v, получено: %v", nil, dest.InnerPtr)
	}
}

func checkType(src SourceStruct, dest DestinationStruct, t *testing.T) {
	if dest.Uint != src.Uint {
		t.Errorf("Ожидалось Uint: %d, получено: %d", src.Uint, dest.Uint)
	}
	if dest.UintPtr == nil || *dest.UintPtr != src.UintPtr {
		t.Errorf("Ожидалось IntToUintPtr: %d, получено: %v", src.UintPtr, *dest.UintPtr)
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

	//Slice
	if dest.Slice == nil || dest.Slice[0] != 1 {
		t.Errorf("Ожидалось Slice[0]: %v, получено: %v", 1, dest.Slice[0])
	}
	if dest.Slice == nil || dest.Slice[1] != 2 {
		t.Errorf("Ожидалось Slice[1]: %v, получено: %v", 2, dest.Slice[1])
	}
	if dest.Slice == nil || dest.Slice[2] != 3 {
		t.Errorf("Ожидалось Slice[2]: %v, получено: %v", 3, dest.Slice[2])
	}
	if dest.Slice == nil || len(dest.Slice) != len(src.Slice) {
		t.Errorf("Ожидалось Slice.len: %v, получено: %v", len(src.Slice), len(dest.Slice))
	}

	// SliceString
	if dest.SliceString == nil || dest.SliceString[0] != "11" {
		t.Errorf("Ожидалось SliceString[0]: %v, получено: %v", 1, dest.SliceString[0])
	}
	if dest.SliceString == nil || dest.SliceString[1] != "22" {
		t.Errorf("Ожидалось SliceString[1]: %v, получено: %v", 2, dest.SliceString[1])
	}
	if dest.SliceString == nil || dest.SliceString[2] != "33" {
		t.Errorf("Ожидалось SliceString[2]: %v, получено: %v", 3, dest.SliceString[2])
	}
	if dest.SliceString == nil || len(dest.SliceString) != len(src.SliceString) {
		t.Errorf("Ожидалось SliceString.len: %v, получено: %v", len(src.SliceString), len(dest.SliceString))
	}

	// SliceIntPtr
	if dest.SliceIntPtr == nil || *dest.SliceIntPtr[0] != 11 {
		t.Errorf("Ожидалось SliceIntPtr[0]: %v, получено: %v", 11, *dest.SliceIntPtr[0])
	}
	if dest.SliceIntPtr == nil || *dest.SliceIntPtr[1] != 22 {
		t.Errorf("Ожидалось SliceIntPtr[1]: %v, получено: %v", 22, *dest.SliceIntPtr[1])
	}
	if dest.SliceIntPtr == nil || *dest.SliceIntPtr[2] != 33 {
		t.Errorf("Ожидалось SliceIntPtr[2]: %v, получено: %v", 33, *dest.SliceIntPtr[2])
	}
	if dest.SliceIntPtr == nil || len(dest.SliceIntPtr) != len(src.SliceIntPtr) {
		t.Errorf("Ожидалось SliceIntPtr.len: %v, получено: %v", len(src.SliceIntPtr), len(dest.SliceIntPtr))
	}

}

func checkTypeInner(src SourceStructInner, dest DestinationStructInner, t *testing.T) {
	if dest.Uint != src.Uint {
		t.Errorf("Ожидалось Uint: %d, получено: %d", src.Uint, dest.Uint)
	}
	if dest.UintPtr == nil || *dest.UintPtr != src.UintPtr {
		t.Errorf("Ожидалось IntToUintPtr: %d, получено: %v", src.UintPtr, *dest.UintPtr)
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

func checkTypeEmpty(src SourceStruct, dest DestinationStruct, t *testing.T) {
	if dest.Uint != src.Uint {
		t.Errorf("Ожидалось Uint: %d, получено: %d", src.Uint, dest.Uint)
	}
	if dest.UintPtr != nil {
		t.Errorf("Ожидалось IntToUintPtr: %v, получено: %v", nil, dest.UintPtr)
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

	// Slice
	if dest.Slice != nil {
		t.Errorf("Ожидалось Slice: %v, получено: %v", nil, dest.Slice)
	}
	if dest.SliceString != nil {
		t.Errorf("Ожидалось SliceString: %v, получено: %v", nil, dest.SliceString)
	}
	if dest.SliceIntPtr != nil {
		t.Errorf("Ожидалось SliceIntPtr: %v, получено: %v", nil, dest.SliceIntPtr)
	}
}

func checkTypeInnerEmpty(src SourceStructInner, dest DestinationStructInner, t *testing.T) {
	if dest.Uint != src.Uint {
		t.Errorf("Ожидалось Uint: %d, получено: %d", src.Uint, dest.Uint)
	}
	if dest.UintPtr != nil {
		t.Errorf("Ожидалось IntToUintPtr: %v, получено: %v", nil, dest.UintPtr)
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
