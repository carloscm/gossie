package gossie

import (
	"testing"
)

/*
	to do: test New*
*/

func TestUUID(t *testing.T) {
	var b []byte
	var r UUID

	b = []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}

	s := "00112233-4455-6677-8899-aabbccddeeff"
	v, err := ParseUUID(s)
	if err != nil {
		t.Error("Unexpected error in NewUUID")
	}
	checkFullMarshal(t, b, BytesType, &v, &r)

	if v.String() != s {
		t.Error("Wrong UUID to string conversion ", v.String())
	}
}
