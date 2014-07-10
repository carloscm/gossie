package mockgossie

import (
	"bytes"

	. "github.com/wadey/gossie/src/cassandra"
	. "github.com/wadey/gossie/src/gossie"
)

type MockQuery struct {
	pool        *MockConnectionPool
	mapping     Mapping
	components  []interface{}
	columnLimit int
	rowLimit    int
}

var _ Query = &MockQuery{}

func (*MockQuery) ConsistencyLevel(ConsistencyLevel) Query { panic("ConsistencyLevel not implemented") }
func (m *MockQuery) Limit(c, r int) Query                  { m.columnLimit = c; m.rowLimit = r; return m }
func (*MockQuery) Reversed(bool) Query                     { panic("Reversed not implemented") }
func (m *MockQuery) Components(c ...interface{}) Query     { m.components = c; return m }
func (*MockQuery) Between(start, end interface{}) Query    { panic("Between not implemented") }
func (*MockQuery) Where(field string, op Operator, value interface{}) Query {
	panic("Where not implemented")
}

func (m *MockQuery) RangeGet(r *Range) (Result, error) {
	rows := m.pool.Rows(m.mapping.Cf())
	count := r.Count
	started := false
	keys := make([]interface{}, 0)
	for _, row := range rows {
		if started || r.Start == nil || bytes.Equal(row.Key, r.Start) {
			started = true
			keys = append(keys, row.Key)
			count--
		}
		if count <= 0 || (r.End != nil && bytes.Equal(row.Key, r.End)) {
			break
		}
	}
	return m.MultiGet(keys)
}

func (m *MockQuery) MultiGet(keys []interface{}) (Result, error) {
	rows := m.pool.Rows(m.mapping.Cf())

	buffer := make([]*Row, 0)
	for _, key := range keys {
		k, err := m.mapping.MarshalKey(key)
		if err != nil {
			return nil, err
		}

		for _, r := range rows {
			if bytes.Equal(r.Key, k) {
				checkExpired(r)
				buffer = append(buffer, r)
			}
		}
	}

	return &result{
		MockQuery: *m,
		buffer:    buffer,
	}, nil
}

func (m *MockQuery) Get(key interface{}) (Result, error) {
	return m.MultiGet([]interface{}{key})
}

func (m *MockQuery) GetOne(key interface{}, destination interface{}) error {
	res, err := m.Get(key)
	if err != nil {
		return err
	}
	return res.Next(destination)
}

var rangeOne = &Range{Count: 1}

func (m *MockQuery) RangeOne(destination interface{}) error {
	res, err := m.RangeGet(rangeOne)
	if err != nil {
		return err
	}
	return res.Next(destination)
}

type result struct {
	MockQuery
	buffer   []*Row
	row      *Row
	position int
}

func (r *result) feedRow() error {
	if r.row == nil {
		if len(r.buffer) <= 0 {
			return Done
		}
		r.row = r.buffer[0]
		r.position = 0
		r.buffer = r.buffer[1:len(r.buffer)]
	}
	return nil
}

func (r *result) Key() ([]byte, error) {
	if err := r.feedRow(); err != nil {
		return nil, err
	}
	return r.row.Key, nil
}

func (r *result) NextColumn() (*Column, error) {
	if err := r.feedRow(); err != nil {
		return nil, err
	}
	if r.position >= len(r.row.Columns) {
		if r.position >= r.columnLimit {
			return nil, EndAtLimit
		} else {
			return nil, EndBeforeLimit
		}
	}
	c := r.row.Columns[r.position]
	r.position++
	return c, nil
}

func (r *result) Rewind() {
	r.position--
	if r.position < 0 {
		r.position = 0
	}
}

func (r *result) Next(destination interface{}) error {
	err := r.mapping.Unmap(destination, r)
	if err == Done {
		// force new row feed and try again, just once
		r.row = nil
		err = r.mapping.Unmap(destination, r)
	}
	return err
}
