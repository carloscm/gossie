// Mock in-memory implementation for gossie. use NewMockConnectionPool to
// create a gossie.ConnectionPool that stores Batch() mutations in an internal
// map[string][]*gossie.Row.
//
// TODO:
//   - Not all methods are implemented
//   - Change the internal map key from string to []byte
//
// Warning: API not finalized, subject to change.
//
// See the example.
package mockgossie

import (
	"bytes"
	. "github.com/wadey/gossie/src/gossie"
)

type MockConnectionPool struct {
	Data map[string]Rows
}

func NewMockConnectionPool() *MockConnectionPool {
	return &MockConnectionPool{
		Data: make(map[string]Rows),
	}
}

func (*MockConnectionPool) Keyspace() string { return "MockKeyspace" }
func (*MockConnectionPool) Schema() *Schema  { panic("Schema Not Implemented") }
func (*MockConnectionPool) Reader() Reader   { panic("Reader Not Implemented") }
func (m *MockConnectionPool) Writer() Writer { return newWriter(m) }
func (m *MockConnectionPool) Batch() Batch   { return newBatch(m) }
func (*MockConnectionPool) Close() error     { return nil }

func (m *MockConnectionPool) Query(mapping Mapping) Query {
	return &MockQuery{
		pool:        m,
		mapping:     mapping,
		rowLimit:    DEFAULT_ROW_LIMIT,
		columnLimit: DEFAULT_COLUMN_LIMIT,
	}
}

type Rows []*Row

func (r Rows) Len() int           { return len(r) }
func (r Rows) Less(i, j int) bool { return bytes.Compare(r[i].Key, r[j].Key) < 0 }
func (r Rows) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

type Columns []*Column

func (r Columns) Len() int           { return len(r) }
func (r Columns) Less(i, j int) bool { return bytes.Compare(r[i].Name, r[j].Name) < 0 }
func (r Columns) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

func (m *MockConnectionPool) Rows(cf string) Rows {
	rows, ok := m.Data[cf]
	if !ok {
		rows = make([]*Row, 0)
		m.Data[cf] = rows
	}
	return rows
}
