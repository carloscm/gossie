package gossie

import (
	"encoding/json"
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

func TestUUIDJsonMarshalling(t *testing.T) {
	s := "00112233-4455-6677-8899-aabbccddeeff"
	v, err := ParseUUID(s)
	if err != nil {
		t.Error("Unexpected error in NewUUID")
	}
	bv, err := v.MarshalJSON()
	if err != nil {
		t.Error("Unexpected error in MarshalJSON")
	}
	if string(bv) != `"`+s+`"` {
		t.Error("Unexpected JSON marshalling result")
	}
}

func TestUUIDJsonUnmarshalling(t *testing.T) {
	js := `"00112233-4455-6677-8899-aabbccddeeff"`
	var v UUID
	err := v.UnmarshalJSON([]byte(js))
	if err != nil {
		t.Error("Unexpected error in UnmarshalJSON", err)
	}
	if v.String() != "00112233-4455-6677-8899-aabbccddeeff" {
		t.Error("Unexpected JSON unmarshalling result")
	}
}

type uuidJsonStruct struct {
	U UUID
}

func TestUUIDJsonStructMarshalling(t *testing.T) {
	s := "00112233-4455-6677-8899-aabbccddeeff"
	v, err := ParseUUID(s)
	if err != nil {
		t.Error("Unexpected error in NewUUID")
	}
	us := &uuidJsonStruct{v}
	bv, err := json.Marshal(us)
	if err != nil {
		t.Error("Unexpected error in JSON Marshalling: ", err)
	}
	us1 := &uuidJsonStruct{}
	err = json.Unmarshal(bv, us1)
	if err != nil {
		t.Error("Unexpected error in JSON Marshalling")
	}
	if us.U != v {
		t.Error("Unexpected JSON struct value")
	}
}
