package gossie

import "encoding/json"

// jsonType is a GossieType that makes it easy to marshal any struct field to/from JSON
// It is enabled with `marshal:"json"`
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
	v interface{}
}

func (t *jsonType) Marshaler(v interface{}) Marshaler {
	return &jsonMarshaler{v}
}

func (t *jsonType) Unmarshaler(v interface{}) Unmarshaler {
	return &jsonMarshaler{v}
}

func (m *jsonMarshaler) MarshalCassandra() ([]byte, error) {
	return json.Marshal(m.v)
}

func (m *jsonMarshaler) UnmarshalCassandra(b []byte) error {
	return json.Unmarshal(b, m.v)
}
