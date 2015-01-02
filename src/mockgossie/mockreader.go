package mockgossie

import (
	"bytes"

	. "github.com/wadey/gossie/src/cassandra"
	. "github.com/wadey/gossie/src/gossie"
)

type MockReader struct {
	pool        *MockConnectionPool
	columnLimit int
	rowLimit    int
	cf          string
}

var _ Reader = &MockReader{}

func (m *MockReader) ConsistencyLevel(ConsistencyLevel) Reader              { panic("not implemented") }
func (m *MockReader) Slice(*Slice) Reader                                   { panic("not implemented") }
func (m *MockReader) Columns([][]byte) Reader                               { panic("not implemented") }
func (m *MockReader) Where(column []byte, op Operator, value []byte) Reader { panic("not implemented") }
func (m *MockReader) MultiGet(keys [][]byte) ([]*Row, error)                { panic("not implemented") }
func (m *MockReader) Count(key []byte) (int, error)                         { panic("not implemented") }
func (m *MockReader) MultiCount(keys [][]byte) ([]*RowColumnCount, error)   { panic("not implemented") }
func (m *MockReader) RangeGet(*Range) ([]*Row, error)                       { panic("not implemented") }
func (m *MockReader) IndexedGet(*IndexedRange) ([]*Row, error)              { panic("not implemented") }
func (m *MockReader) SetTokenRange(startToken, endToken string) Reader      { panic("not implemented") }
func (m *MockReader) RangeScan() (data <-chan *Row, err <-chan error)       { panic("not implemented") }
func (m *MockReader) WideRowScan(key, startColumn []byte, batchSize int32, callback func(*Column) bool) error {
	panic("not implemented")
}

func newReader(m *MockConnectionPool) *MockReader {
	return &MockReader{
		pool:        m,
		rowLimit:    DEFAULT_ROW_LIMIT,
		columnLimit: DEFAULT_COLUMN_LIMIT,
	}
}

func (m *MockReader) Cf(cf string) Reader {
	m.cf = cf
	return m
}

func (m *MockReader) Get(key []byte) (*Row, error) {
	rows := m.pool.Rows(m.cf)

	for _, r := range rows {
		if bytes.Equal(r.Key, key) {
			checkExpired(r)
			return r, nil
		}
	}
	return nil, nil
}
