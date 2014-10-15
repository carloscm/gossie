package gossie

import "encoding/json"

// jsonType is a GossieType that makes it easy to marshal any struct field to/from JSON
// It is enabled with `marshal:"json"`
//
// The value to use when the field is empty can be provided as an argument,
// such as `marshal:"json,{}"` to use `{}` for nil.
//
// Example:
//
//  type Round struct {
//      ID   string `cf:"rounds" key:"ID"`
//      Bets []*Bet `name:"bets" marshal:"json"`
//  }
type jsonType struct {
}

type jsonMarshaler struct {
	v     interface{}
	empty *string
}

func (t *jsonType) Marshaler(v interface{}, tagArgs *string) Marshaler {
	return &jsonMarshaler{v, tagArgs}
}

func (t *jsonType) Unmarshaler(v interface{}, tagArgs *string) Unmarshaler {
	return &jsonMarshaler{v, tagArgs}
}

func (m *jsonMarshaler) MarshalCassandra() ([]byte, error) {
	b, err := json.Marshal(m.v)
	if err != nil {
		return b, err
	}
	if string(b) == "null" && m.empty != nil {
		b = []byte(*m.empty)
	}
	return b, err
}

func (m *jsonMarshaler) UnmarshalCassandra(b []byte) error {
	if m.empty != nil && string(b) == *m.empty {
		b = []byte("null")
	}
	return json.Unmarshal(b, m.v)
}
