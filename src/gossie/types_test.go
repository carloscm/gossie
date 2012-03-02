package gossie

import (
    "testing"
    "reflect"
)

func checkMarshal(t *testing.T, value interface{}, good []byte, typeDesc TypeDesc) {
    b, err := Marshal(value, typeDesc)
    if err != nil {
        t.Error("Error marshalling value: ", err)
    }
    if !reflect.DeepEqual(b, good) {
        t.Error("Marshalled value does not match expected ", good, " actual ", b)
    }
}

func checkUnmarshal(t *testing.T, b []byte, typeDesc TypeDesc, good interface{}, value interface{}) {
    err := Unmarshal(b, typeDesc, value)
    if err != nil {
        t.Error("Error unmarshalling value: ", err)
    }
    if !reflect.DeepEqual(value, good) {
        t.Error("Unmarshalled value does not match expected ", good, " actual ", value)
    }
}

func checkFullMarshal(t *testing.T, marshalled []byte, typeDesc TypeDesc, goodValue interface{}, retValue interface{}) {
    checkMarshal(t, goodValue, marshalled, typeDesc)
    checkUnmarshal(t, marshalled, typeDesc, goodValue, retValue)
    checkMarshal(t, retValue, marshalled, typeDesc)
}

func errorMarshal(t *testing.T, value interface{}, typeDesc TypeDesc) {
    _, err := Marshal(value, typeDesc)
    if err == nil {
        t.Error("Error expected for marshalling, got none")
    }
}

func TestMarshalWrongType(t *testing.T) {
    type no int
    var v no = 1

    errorMarshal(t, v, BytesType)
    errorMarshal(t, v, AsciiType)
    errorMarshal(t, v, UTF8Type)
    errorMarshal(t, v, LongType)
    errorMarshal(t, v, IntegerType)
    errorMarshal(t, v, DecimalType)
    errorMarshal(t, v, UUIDType)
    errorMarshal(t, v, BooleanType)
    errorMarshal(t, v, FloatType)
    errorMarshal(t, v, DoubleType)
    errorMarshal(t, v, DateType)
}

func TestMarshalBytes(t *testing.T) {
    var b []byte
    var v []byte
    var r []byte

    v = []byte{4, 2}
    b = []byte{4, 2}
    checkFullMarshal(t, b, BytesType, &v, &r)
    checkFullMarshal(t, b, AsciiType, &v, &r)
    checkFullMarshal(t, b, UTF8Type, &v, &r)
    checkFullMarshal(t, b, LongType, &v, &r)
    checkFullMarshal(t, b, IntegerType, &v, &r)
    checkFullMarshal(t, b, DecimalType, &v, &r)
    checkFullMarshal(t, b, UUIDType, &v, &r)
    checkFullMarshal(t, b, BooleanType, &v, &r)
    checkFullMarshal(t, b, FloatType, &v, &r)
    checkFullMarshal(t, b, DoubleType, &v, &r)
    checkFullMarshal(t, b, DateType, &v, &r)
}

func TestMarshalBool(t *testing.T) {
    var b []byte
    var v bool
    var r bool

    v = false

    b = []byte{0}
    checkFullMarshal(t, b, BytesType, &v, &r)
    checkFullMarshal(t, b, BooleanType, &v, &r)

    b = []byte{'0'}
    checkFullMarshal(t, b, AsciiType, &v, &r)
    checkFullMarshal(t, b, UTF8Type, &v, &r)

    b = []byte{0, 0, 0, 0, 0, 0, 0, 0}
    checkFullMarshal(t, b, LongType, &v, &r)

    v = true

    b = []byte{1}
    checkFullMarshal(t, b, BytesType, &v, &r)
    checkFullMarshal(t, b, BooleanType, &v, &r)

    b = []byte{'1'}
    checkFullMarshal(t, b, AsciiType, &v, &r)
    checkFullMarshal(t, b, UTF8Type, &v, &r)

    b = []byte{0, 0, 0, 0, 0, 0, 0, 1}
    checkFullMarshal(t, b, LongType, &v, &r)

    errorMarshal(t, v, IntegerType)
    errorMarshal(t, v, DecimalType)
    errorMarshal(t, v, UUIDType)
    errorMarshal(t, v, FloatType)
    errorMarshal(t, v, DoubleType)
    errorMarshal(t, v, DateType)
}

func TestMarshalInt(t *testing.T) {
    var b []byte
    var v64 int64
    var v32 int32
    var vi int
    var v16 int16
    var v8 int8
    var r64 int64
    var r32 int32
    var ri int
    var r16 int16
    var r8 int8

    // positive

    v64 = 4611686018427387907
    v32 = 1073741827
    vi = 1073741827
    v16 = 16387
    v8 = 67

    b = []byte{0x40, 0, 0, 0, 0, 0, 0, 3}
    checkFullMarshal(t, b, LongType, &v64, &r64)
    checkFullMarshal(t, b, DateType, &v64, &r64)
    checkFullMarshal(t, b, BytesType, &v64, &r64)

    b = []byte{0, 0, 0, 0, 0x40, 0, 0, 3}
    checkFullMarshal(t, b, LongType, &v32, &r32)
    checkFullMarshal(t, b, LongType, &vi, &ri)

    b = []byte{0, 0, 0, 0, 0, 0, 0x40, 3}
    checkFullMarshal(t, b, LongType, &v16, &r16)

    b = []byte{0, 0, 0, 0, 0, 0, 0, 0x43}
    checkFullMarshal(t, b, LongType, &v8, &r8)

    b = []byte{0x40, 0, 0, 3}
    checkFullMarshal(t, b, BytesType, &v32, &r32)
    checkFullMarshal(t, b, BytesType, &vi, &ri)

    b = []byte{0x40, 3}
    checkFullMarshal(t, b, BytesType, &v16, &r16)

    b = []byte{0x43}
    checkFullMarshal(t, b, BytesType, &v8, &r8)

    b = []byte{'4', '6', '1', '1', '6', '8', '6', '0', '1', '8', '4', '2', '7', '3', '8', '7', '9', '0', '7'}
    checkFullMarshal(t, b, AsciiType, &v64, &r64)
    checkFullMarshal(t, b, UTF8Type, &v64, &r64)

    b = []byte{'1', '0', '7', '3', '7', '4', '1', '8', '2', '7'}
    checkFullMarshal(t, b, AsciiType, &v32, &r32)
    checkFullMarshal(t, b, UTF8Type, &v32, &r32)
    checkFullMarshal(t, b, AsciiType, &vi, &ri)
    checkFullMarshal(t, b, UTF8Type, &vi, &ri)

    b = []byte{'1', '6', '3', '8', '7'}
    checkFullMarshal(t, b, AsciiType, &v16, &r16)
    checkFullMarshal(t, b, UTF8Type, &v16, &r16)

    b = []byte{'6', '7'}
    checkFullMarshal(t, b, AsciiType, &v8, &r8)
    checkFullMarshal(t, b, UTF8Type, &v8, &r8)

    // negative

    v64 = -9223372036854775805
    v32 = -2147483645
    vi = -2147483645
    v16 = -32765
    v8 = -125

    b = []byte{0x80, 0, 0, 0, 0, 0, 0, 3}
    checkFullMarshal(t, b, LongType, &v64, &r64)
    checkFullMarshal(t, b, DateType, &v64, &r64)
    checkFullMarshal(t, b, BytesType, &v64, &r64)

    b = []byte{0xff, 0xff, 0xff, 0xff, 0x80, 0, 0, 3}
    checkFullMarshal(t, b, LongType, &v32, &r32)
    checkFullMarshal(t, b, LongType, &vi, &ri)

    b = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x80, 3}
    checkMarshal(t, v16, b, LongType)

    b = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x83}
    checkMarshal(t, v8, b, LongType)

    b = []byte{0x80, 0, 0, 3}
    checkFullMarshal(t, b, BytesType, &v32, &r32)
    checkFullMarshal(t, b, BytesType, &vi, &ri)

    b = []byte{0x80, 3}
    checkFullMarshal(t, b, BytesType, &v16, &r16)

    b = []byte{0x83}
    checkFullMarshal(t, b, BytesType, &v8, &r8)

    b = []byte{'-', '9', '2', '2', '3', '3', '7', '2', '0', '3', '6', '8', '5', '4', '7', '7', '5', '8', '0', '5'}
    checkFullMarshal(t, b, AsciiType, &v64, &r64)
    checkFullMarshal(t, b, UTF8Type, &v64, &r64)

    b = []byte{'-', '2', '1', '4', '7', '4', '8', '3', '6', '4', '5'}
    checkFullMarshal(t, b, AsciiType, &v32, &r32)
    checkFullMarshal(t, b, UTF8Type, &v32, &r32)
    checkFullMarshal(t, b, AsciiType, &vi, &ri)
    checkFullMarshal(t, b, UTF8Type, &vi, &ri)

    b = []byte{'-', '3', '2', '7', '6', '5'}
    checkFullMarshal(t, b, AsciiType, &v16, &r16)
    checkFullMarshal(t, b, UTF8Type, &v16, &r16)

    b = []byte{'-', '1', '2', '5'}
    checkFullMarshal(t, b, AsciiType, &v8, &r8)
    checkFullMarshal(t, b, UTF8Type, &v8, &r8)

    // errors

    errorMarshal(t, vi, IntegerType)
    errorMarshal(t, vi, DecimalType)
    errorMarshal(t, vi, UUIDType)
    errorMarshal(t, vi, FloatType)
    errorMarshal(t, vi, DoubleType)

    errorMarshal(t, v32, DateType)
    errorMarshal(t, vi, DateType)
    errorMarshal(t, v16, DateType)
    errorMarshal(t, v8, DateType)
}

func TestMarshalString(t *testing.T) {
    var b []byte
    var v string = "cáñamo"
    var r string

    b = []byte{'c', 0xc3, 0xa1, 0xc3, 0xb1, 'a', 'm', 'o'}
    checkFullMarshal(t, b, BytesType, &v, &r)
    checkFullMarshal(t, b, AsciiType, &v, &r) // NOTE: this lib does not perform ascii/utf8 checking for now
    checkFullMarshal(t, b, UTF8Type, &v, &r)  // NOTE: this lib does not perform ascii/utf8 checking for now

    errorMarshal(t, v, LongType)

    v = "4611686018427387907"
    b = []byte{0x40, 0, 0, 0, 0, 0, 0, 3}
    checkFullMarshal(t, b, LongType, &v, &r)

    v = "-9223372036854775805"
    b = []byte{0x80, 0, 0, 0, 0, 0, 0, 3}
    checkFullMarshal(t, b, LongType, &v, &r)

    errorMarshal(t, v, IntegerType)
    errorMarshal(t, v, DecimalType)
    errorMarshal(t, v, UUIDType)
    errorMarshal(t, v, BooleanType)
    errorMarshal(t, v, FloatType)
    errorMarshal(t, v, DoubleType)
    errorMarshal(t, v, DateType)
}

func TestMarshalUUID(t *testing.T) {
    var b []byte
    var v UUID = [16]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
    var r UUID

    b = []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
    checkFullMarshal(t, b, BytesType, &v, &r)
    checkFullMarshal(t, b, UUIDType, &v, &r)

    errorMarshal(t, v, LongType)
    errorMarshal(t, v, AsciiType)
    errorMarshal(t, v, UTF8Type)
    errorMarshal(t, v, IntegerType)
    errorMarshal(t, v, DecimalType)
    errorMarshal(t, v, BooleanType)
    errorMarshal(t, v, FloatType)
    errorMarshal(t, v, DoubleType)
    errorMarshal(t, v, DateType)

    // test utility functions

    s := "00112233-4455-6677-8899-aabbccddeeff"
    v2, err := NewUUID(s)
    if err != nil {
        t.Error("Unexpected error in NewUUID")
    }
    checkFullMarshal(t, b, BytesType, &v2, &r)

    if v2.String() != s {
        t.Error("Wrong UUID to string conversion ", v2.String())
    }
}

func TestMarshalFloat(t *testing.T) {
    var b []byte
    var v32 float32 = float32(-1.01)
    var r32 float32
    var v64 float64 = float64(-1.01)
    var r64 float64

    b = []byte{0xbf, 0x81, 0x47, 0xae}
    checkFullMarshal(t, b, BytesType, &v32, &r32)
    checkFullMarshal(t, b, FloatType, &v32, &r32)

    //checkFullMarshal(t, b, FloatType, &v64, &r32)

    b = []byte{0xbf, 0xf0, 0x28, 0xf5, 0xc2, 0x8f, 0x5c, 0x29}
    checkFullMarshal(t, b, BytesType, &v64, &r64)
    checkFullMarshal(t, b, DoubleType, &v64, &r64)

    //bLossy := []byte {0xbf, 0xf0, 0x28, 0xf5, 0xc0, 0, 0, 0}
    //checkFullMarshal(t, bLossy, DoubleType, &v32, &r64)

    errorMarshal(t, v32, LongType)
    errorMarshal(t, v32, AsciiType)
    errorMarshal(t, v32, UTF8Type)
    errorMarshal(t, v32, IntegerType)
    errorMarshal(t, v32, DecimalType)
    errorMarshal(t, v32, BooleanType)
    errorMarshal(t, v32, DoubleType)
    errorMarshal(t, v32, DateType)

    errorMarshal(t, v64, LongType)
    errorMarshal(t, v64, AsciiType)
    errorMarshal(t, v64, UTF8Type)
    errorMarshal(t, v64, IntegerType)
    errorMarshal(t, v64, DecimalType)
    errorMarshal(t, v64, BooleanType)
    errorMarshal(t, v64, FloatType)
    errorMarshal(t, v64, DateType)
}
