package mockgossie

import (
	"bytes"
	enc "encoding/binary"
	"reflect"

	. "github.com/wadey/gossie/src/cassandra"
	. "github.com/wadey/gossie/src/gossie"
)

const (
	eocEquals  byte = 0
	eocGreater byte = 1
	eocLower   byte = 0xff
)

type MockQuery struct {
	pool         *MockConnectionPool
	mapping      Mapping
	components   []interface{}
	betweenStart interface{}
	betweenEnd   interface{}
	columnLimit  int
	rowLimit     int
	reversed     bool
}

var _ Query = &MockQuery{}

func (*MockQuery) ConsistencyLevel(ConsistencyLevel) Query { panic("ConsistencyLevel not implemented") }
func (m *MockQuery) Limit(c, r int) Query                  { m.columnLimit = c; m.rowLimit = r; return m }
func (m *MockQuery) Components(c ...interface{}) Query     { m.components = c; return m }
func (m *MockQuery) Reversed(r bool) Query {
	m.reversed = r
	return m
}
func (*MockQuery) Where(field string, op Operator, value interface{}) Query {
	panic("Where not implemented")
}

func (m *MockQuery) Between(start, end interface{}) Query {
	if start != nil && reflect.ValueOf(start).IsNil() {
		start = nil
	}
	if end != nil && reflect.ValueOf(end).IsNil() {
		end = nil
	}
	m.betweenStart, m.betweenEnd = start, end
	return m
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
				r, err = m.sliceRow(r)
				if err != nil {
					return nil, err
				}
				buffer = append(buffer, r)
			}
		}
	}

	return &result{
		MockQuery: *m,
		buffer:    buffer,
	}, nil
}

func (m *MockQuery) sliceRow(r *Row) (*Row, error) {
	if m.components != nil || m.betweenStart != nil || m.betweenEnd != nil || m.columnLimit != 0 {
		slice, err := m.buildSlice()
		if err != nil {
			return nil, err
		}
		if m.reversed {
			slice.Start, slice.End = slice.End, slice.Start
		}
		cr := *r
		cr.Columns = []*Column{}
		for _, c := range r.Columns {
			if len(slice.Start) > 0 && bytes.Compare(slice.Start, c.Name) > 0 {
				continue
			}
			if len(slice.End) > 0 && bytes.Compare(slice.End, c.Name) < 0 {
				continue
			}
			cr.Columns = append(cr.Columns, c)
		}
		if slice.Count != 0 && len(cr.Columns) > slice.Count {
			if m.reversed {
				cr.Columns = cr.Columns[(len(cr.Columns) - slice.Count):len(cr.Columns)]
			} else {
				cr.Columns = cr.Columns[0:slice.Count]
			}
		}
		r = &cr
	}
	return r, nil
}

func (q *MockQuery) buildSlice() (*Slice, error) {
	var start, end []byte

	components := q.components

	if q.mapping.Compact() && len(q.mapping.Components()) == 1 {
		if len(components) == 1 {
			c := components[0]
			b, err := q.mapping.MarshalComponent(c, 0)
			if err != nil {
				return nil, err
			}
			start = b
			end = b
		} else {
			if q.betweenStart != nil {
				b, err := q.mapping.MarshalComponent(q.betweenStart, 0)
				if err != nil {
					return nil, err
				}
				start = b
			}
			if q.betweenEnd != nil {
				b, err := q.mapping.MarshalComponent(q.betweenEnd, 0)
				if err != nil {
					return nil, err
				}
				end = b
			}
		}
	} else if len(components) > 0 {
		last := len(components) - 1
		for i, c := range components {
			b, err := q.mapping.MarshalComponent(c, i)
			if err != nil {
				return nil, err
			}
			start = append(start, packComposite(b, eocEquals)...)
			if i == last {
				if q.betweenEnd == nil {
					end = append(end, packComposite(b, eocGreater)...)
				} else {
					end = append(end, packComposite(b, eocEquals)...)
				}
			} else {
				end = append(end, packComposite(b, eocEquals)...)
			}
		}
	}
	if q.betweenStart != nil {
		b, err := q.mapping.MarshalComponent(q.betweenStart, len(components))
		if err != nil {
			return nil, err
		}
		start = append(start, packComposite(b, eocEquals)...)
	}
	if q.betweenEnd != nil {
		b, err := q.mapping.MarshalComponent(q.betweenEnd, len(components))
		if err != nil {
			return nil, err
		}
		end = append(end, packComposite(b, eocEquals)...)
	}

	return &Slice{Start: start, End: end, Count: q.columnLimit, Reversed: q.reversed}, nil
}

func packComposite(component []byte, eoc byte) []byte {
	r := make([]byte, 2)
	enc.BigEndian.PutUint16(r, uint16(len(component)))
	r = append(r, component...)
	return append(r, eoc)
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
	var c *Column
	if r.MockQuery.reversed {
		c = r.row.Columns[len(r.row.Columns)-1-r.position]
	} else {
		c = r.row.Columns[r.position]
	}
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
