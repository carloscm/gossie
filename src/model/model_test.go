package model

import (
	"testing"
)

func TestLong(t *testing.T) {

	var l Long = 3

	b := l.Bytes()
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

	b = make([]byte, 8)
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
	var u2 BaseValue = u

	_, ok := (interface{}(u)).(BaseValue)
	if (!ok) {
		t.Error("UTF8 fails BaseValue type assertion")
	}

	b := u.Bytes()
	if (len(b) != 8) {
		t.Error("UTF8 bytes size is not 8")
	}

	if (b[0] != 0x63 ||
		b[1] != 0xc3 ||
		b[2] != 0xa1 ||
		b[3] != 0xc3 ||
		b[4] != 0xb1 ||
		b[5] != 0x61 ||
		b[6] != 0x6d ||
		b[7] != 0x6f) {
		t.Error("UTF8 serialization is wrong")
	}

	b = make([]byte, 8)
	b[0] = 0x6d
	b[1] = 0x6f
	b[2] = 0x63
	b[3] = 0x61
	b[4] = 0xc3
	b[5] = 0xb1
	b[6] = 0xc3
	b[7] = 0xa1

	u.SetBytes(b)
	if (u != "mocañá") {
		t.Error("UTF8 unserialization is wrong")
	}

}
