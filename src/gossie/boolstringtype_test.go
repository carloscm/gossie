package gossie

import "testing"

type CustomBoolStringType struct {
	B bool `marshal:"boolstring"`
}

func TestBoolStringMarshaler(t *testing.T) {
	v := &CustomBoolStringType{B: true}
	ret := &CustomBoolStringType{}

	b := []byte(`true`)
	jv := &boolStringMarshaler{&v.B}
	jret := &boolStringMarshaler{&ret.B}
	checkFullMarshal(t, b, BytesType, jv, jret)
}
