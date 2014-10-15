package gossie

import (
	"fmt"
	"strconv"
)

// boolStringType is a GossieType that makes it easy to marshal from true/false to "true"/"false"
// It is enabled with `marshal:"boolstring"`
//
// Example:
//
//  type Round struct {
//      ID     string `cf:"rounds" key:"ID"`
//      Backed bool   `name:"backed" marshal:"boolstring"`
//  }
type boolStringType struct {
}

type boolStringMarshaler struct {
	b bool
}
type boolStringUnmarshaler struct {
	b *bool
}

func (b *boolStringType) Marshaler(v interface{}, tagArgs *string) Marshaler {
	return &boolStringMarshaler{v.(bool)}
}
func (b *boolStringType) Unmarshaler(v interface{}, tagArgs *string) Unmarshaler {
	return &boolStringUnmarshaler{v.(*bool)}
}

func (m *boolStringMarshaler) MarshalCassandra() ([]byte, error) {
	return []byte(strconv.FormatBool(m.b)), nil
}

func (m *boolStringUnmarshaler) UnmarshalCassandra(b []byte) error {
	switch string(b) {
	case "true":
		*m.b = true
	case "", "false":
		*m.b = false
	default:
		return fmt.Errorf("invalid boolstring: %v", b)
	}
	return nil
}
