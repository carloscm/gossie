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
	"sort"

	"github.com/apache/thrift/lib/go/thrift"
	. "github.com/wadey/gossie/src/cassandra"
	. "github.com/wadey/gossie/src/gossie"
)

type MockConnectionPool struct {
	Data map[string]Rows
}

var _ ConnectionPool = &MockConnectionPool{}

func NewMockConnectionPool() *MockConnectionPool {
	return &MockConnectionPool{
		Data: make(map[string]Rows),
	}
}

func (*MockConnectionPool) Keyspace() string { return "MockKeyspace" }
func (*MockConnectionPool) Schema() *Schema  { panic("Schema Not Implemented") }
func (m *MockConnectionPool) Reader() Reader { return newReader(m) }
func (m *MockConnectionPool) Writer() Writer { return newWriter(m) }
func (m *MockConnectionPool) Batch() Batch   { return newBatch(m) }
func (*MockConnectionPool) Close()           {}

func (m *MockConnectionPool) WithTracer(Tracer) ConnectionPool { return m }

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

type RowDump map[string][]byte
type CFDump map[string]RowDump
type Dump map[string]CFDump

// Utility method to make validating tests easier
// result is: map[cf]map[rowKey]map[columnName]columnValue
func (m *MockConnectionPool) Dump() Dump {
	d := Dump{}

	for cf, _ := range m.Data {
		d[cf] = m.DumpCF(cf)
	}

	return d
}

// Utility method to make validating tests easier
// result is: map[rowKey]map[columnName]columnValue
func (m *MockConnectionPool) DumpCF(cf string) CFDump {
	d := CFDump{}

	rows, ok := m.Data[cf]
	if !ok {
		return d
	}

	for _, row := range rows {
		rowMap := map[string][]byte{}
		d[string(row.Key)] = rowMap
		for _, column := range row.Columns {
			rowMap[string(column.Name)] = column.Value
		}
	}

	return d
}

// Utility method for loading data in for tests
func (m *MockConnectionPool) Load(dump Dump) {
	for cf, d := range dump {
		m.LoadCF(cf, d)
	}
}

// Utility method for loading data in for tests
func (m *MockConnectionPool) LoadCF(cf string, dump CFDump) {
	rows := []*Row{}
	t := thrift.Int64Ptr(now())
	for key, columns := range dump {
		cols := Columns{}

		for name, value := range columns {
			cols = append(cols, &Column{
				Name:      []byte(name),
				Value:     value,
				Timestamp: t,
			})
		}
		sort.Sort(cols)
		rows = append(rows, &Row{
			Key:     []byte(key),
			Columns: cols,
		})
	}
	m.Data[cf] = rows
}
