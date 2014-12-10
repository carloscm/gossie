package gossie

import "testing"

type CustomJsonType struct {
	B map[string]string `marshal:"json"`
}

func TestJsonMarshaler(t *testing.T) {
	v := &CustomJsonType{B: map[string]string{"foo": "bar"}}
	ret := &CustomJsonType{}

	b := []byte(`{"foo":"bar"}`)
	jv := &jsonMarshaler{&v.B, nil}
	jret := &jsonMarshaler{&ret.B, nil}
	checkFullMarshal(t, b, BytesType, jv, jret)
}
