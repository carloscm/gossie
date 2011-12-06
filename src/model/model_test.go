package model

import (
	"testing"
)

func TestLong(t *testing.T) {

	var l Long = 3

	if (l.Len() != 8) {
		t.Error("Long Len() is not 8")
	}

	var b []byte = make([]byte, 8)
	l.Bytes(b)
	if (len(b) != 8) {
		t.Error("Long bytes size is not 8")
	}

	if (b[0] != 0 ||
		b[1] != 0 ||
		b[2] != 0 ||
		b[3] != 0 ||
		b[4] != 0 ||
		b[5] != 0 ||
		b[6] != 0 ||
		b[7] != 3) {
		t.Error("Long serialization is wrong")
	}

	b[0] = 0x40;
	b[7] = 0;
	l.SetBytes(b)
	if (l != 0x4000000000000000) {
		t.Error("Long unserialization is wrong (63 bit)")
	}

	b[0] = 0xff;
	b[1] = 0xff;
	b[2] = 0xff;
	b[3] = 0xff;
	b[4] = 0xff;
	b[5] = 0xff;
	b[6] = 0xff;
	b[7] = 0xfe;
	l.SetBytes(b)
	if (l != -2) {
		t.Error("Long unserialization is wrong (negative)")
	}

}

func TestUTF8(t *testing.T) {

	var u UTF8 = "cáñamo"

	if (u.Len() != 8) {
		t.Error("UTF8 Len() is not the expected value")
	}

	var b []byte = make([]byte, 8)
	u.Bytes(b)
	if (len(b) != 8) {
		t.Error("UTF8 bytes size is not 8")
	}

	if (b[0] != 0 ||
		b[1] != 0 ||
		b[2] != 0 ||
		b[3] != 0 ||
		b[4] != 0 ||
		b[5] != 0 ||
		b[6] != 0 ||
		b[7] != 3) {
		t.Error("UTF8 serialization is wrong")
	}

}
