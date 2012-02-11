package gossie

import (
	"testing"
)

func checkMarshal(t *testing.T, value interface{}, good []byte, typeDesc TypeDesc) {
	b, err := Marshal(value, typeDesc)
	if err != nil {
		t.Error("Error marshalling integer: ", err)
	}
	if len(good) != len(b) {
		t.Error("Marshalled integer has wrong size, expected ", len(good), " actual ", len(b))
	}
	t.Log("--")
	for i := 0; i < len(good); i++ {
		t.Log(i, good[i], b[i])
		if good[i] != b[i] {
			t.Error("Marshalled integer has wrong serialization, expected ", good[i], " actual ", b[i])
		}
	}
}

func errorMarshal(t *testing.T, value interface{}, typeDesc TypeDesc) {
	_, err := Marshal(value, typeDesc)
	if err == nil {
		t.Error("Error expected for marshalling, got none")
	}
}

func TestMarshalBool(t *testing.T) {
	var b []byte
	var v bool

	v = false

	b = []byte {0}
	checkMarshal(t, v, b, BytesType)
	checkMarshal(t, v, b, BooleanType)

	b = []byte {'0'}
	checkMarshal(t, v, b, AsciiType)
	checkMarshal(t, v, b, UTF8Type)

	b = []byte {0, 0, 0, 0, 0, 0, 0, 0}
	checkMarshal(t, v, b, LongType)

	v = true

	b = []byte {1}
	checkMarshal(t, v, b, BytesType)
	checkMarshal(t, v, b, BooleanType)

	b = []byte {'1'}
	checkMarshal(t, v, b, AsciiType)
	checkMarshal(t, v, b, UTF8Type)

	b = []byte {0, 0, 0, 0, 0, 0, 0, 1}
	checkMarshal(t, v, b, LongType)

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

	// positive

	v64 = 4611686018427387907
	v32 = 1073741827
	vi = 1073741827
	v16 = 16387
	v8 = 67

	b = []byte {0x40, 0, 0, 0, 0, 0, 0, 3}
	checkMarshal(t, v64, b, LongType)
	checkMarshal(t, v64, b, DateType)
	checkMarshal(t, v64, b, BytesType)

	b = []byte {0, 0, 0, 0, 0x40, 0, 0, 3}
	checkMarshal(t, v32, b, LongType)
	checkMarshal(t, vi, b, LongType)

	b = []byte {0, 0, 0, 0, 0, 0, 0x40, 3}
	checkMarshal(t, v16, b, LongType)

	b = []byte {0, 0, 0, 0, 0, 0, 0, 0x43}
	checkMarshal(t, v8, b, LongType)

	b = []byte {0x40, 0, 0, 3}
	checkMarshal(t, v32, b, BytesType)

	b = []byte {0x40, 0, 0, 3}
	checkMarshal(t, vi, b, BytesType)

	b = []byte {0x40, 3}
	checkMarshal(t, v16, b, BytesType)

	b = []byte {0x43}
	checkMarshal(t, v8, b, BytesType)

	b = []byte {'4', '6', '1', '1', '6', '8', '6', '0', '1', '8', '4', '2', '7', '3', '8', '7', '9', '0', '7'}
	checkMarshal(t, v64, b, AsciiType)
	checkMarshal(t, v64, b, UTF8Type)

	b = []byte {'1', '0', '7', '3', '7', '4', '1', '8', '2', '7'}
	checkMarshal(t, v32, b, AsciiType)
	checkMarshal(t, v32, b, UTF8Type)
	checkMarshal(t, vi, b, AsciiType)
	checkMarshal(t, vi, b, UTF8Type)

	b = []byte {'1', '6', '3', '8', '7'}
	checkMarshal(t, v16, b, AsciiType)
	checkMarshal(t, v16, b, UTF8Type)

	b = []byte {'6', '7'}
	checkMarshal(t, v8, b, AsciiType)
	checkMarshal(t, v8, b, UTF8Type)

	// negative

	v64 = -9223372036854775805
	v32 = -2147483645
	vi = -2147483645
	v16 = -32765
	v8 = -125

	b = []byte {0x80, 0, 0, 0, 0, 0, 0, 3}
	checkMarshal(t, v64, b, LongType)
	checkMarshal(t, v64, b, DateType)
	checkMarshal(t, v64, b, BytesType)

	b = []byte {0xff, 0xff, 0xff, 0xff, 0x80, 0, 0, 3}
	checkMarshal(t, v32, b, LongType)
	checkMarshal(t, vi, b, LongType)

	b = []byte {0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x80, 3}
	checkMarshal(t, v16, b, LongType)

	b = []byte {0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x83}
	checkMarshal(t, v8, b, LongType)

	b = []byte {0x80, 0, 0, 3}
	checkMarshal(t, v32, b, BytesType)

	b = []byte {0x80, 0, 0, 3}
	checkMarshal(t, vi, b, BytesType)

	b = []byte {0x80, 3}
	checkMarshal(t, v16, b, BytesType)

	b = []byte {0x83}
	checkMarshal(t, v8, b, BytesType)

	b = []byte {'-', '9', '2', '2', '3', '3', '7', '2', '0', '3', '6', '8', '5', '4', '7', '7', '5', '8', '0', '5'}
	checkMarshal(t, v64, b, AsciiType)
	checkMarshal(t, v64, b, UTF8Type)

	b = []byte {'-', '2', '1', '4', '7', '4', '8', '3', '6', '4', '5'}
	checkMarshal(t, v32, b, AsciiType)
	checkMarshal(t, v32, b, UTF8Type)
	checkMarshal(t, vi, b, AsciiType)
	checkMarshal(t, vi, b, UTF8Type)

	b = []byte {'-', '3', '2', '7', '6', '5'}
	checkMarshal(t, v16, b, AsciiType)
	checkMarshal(t, v16, b, UTF8Type)

	b = []byte {'-', '1', '2', '5'}
	checkMarshal(t, v8, b, AsciiType)
	checkMarshal(t, v8, b, UTF8Type)

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

	b = []byte {'c', 0xc3, 0xa1, 0xc3, 0xb1, 'a', 'm', 'o'}
	checkMarshal(t, v, b, BytesType)
	checkMarshal(t, v, b, AsciiType) // NOTE: this lib does not perform utf8 checking for now
	checkMarshal(t, v, b, UTF8Type)

	errorMarshal(t, v, LongType)

	v = "4611686018427387907"
	b = []byte {0x40, 0, 0, 0, 0, 0, 0, 3}
	checkMarshal(t, v, b, LongType)

	v = "-9223372036854775805"
	b = []byte {0x80, 0, 0, 0, 0, 0, 0, 3}
	checkMarshal(t, v, b, LongType)

    errorMarshal(t, v, IntegerType)
    errorMarshal(t, v, DecimalType)
    errorMarshal(t, v, UUIDType)
    errorMarshal(t, v, BooleanType)
    errorMarshal(t, v, FloatType)
    errorMarshal(t, v, DoubleType)
    errorMarshal(t, v, DateType)
}
