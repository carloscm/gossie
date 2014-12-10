package gossie

import (
	"errors"
	. "github.com/wadey/gossie/src/cassandra"
)

/*
todo:
    autopaging?
    Search() and interface(s) for indexed get
*/

var (
	Done = errors.New("No more results found")
)

const (
	DEFAULT_COLUMN_LIMIT = 10000
	DEFAULT_ROW_LIMIT    = 500
)

// Query is a high level interface for Cassandra queries
type Query interface {

	// ConsistencyLevel sets the consistency level for this particular call.
	// It is optional, if left uncalled it will default to your connection
	// pool options value.
	ConsistencyLevel(ConsistencyLevel) Query

	// Limit sets the column and rows to buffer at once.
	Limit(columns, rows int) Query

	// Reverse set to true will reverse the order of the columns in the result.
	Reversed(bool) Query

	// Components buils a column slice for Get operations with fixed values
	// for the passed components. It fills those components in the same order
	// as in the method arguments.
	Components(components ...interface{}) Query

	// Between allows to pass different values for the last components of the
	// Start and End columns in the column slice for Get operations. Do not
	// pass the last component to Component() when usign Between(). start is
	// inclusive, end is exclusive.
	Between(start, end interface{}) Query

	// Where adds filter expression (see Reader.Where)
	// first field must be a part of secondary index and first op must be EQ
	// use Range to abtain result
	Where(field string, op Operator, value interface{}) Query

	// scan range filtered by Where statement(s)
	RangeGet(*Range) (Result, error)

	// scan range for one record and unmarshal it into destination
	RangeOne(destination interface{}) error

	// Get looks up a row with the given key. If the row uses a composite
	// column names the Result will allow you to iterate over the entire row.
	Get(key interface{}) (Result, error)

	// Like Get() but returns exactly one record or error
	GetOne(key interface{}, destination interface{}) error

	// MultiGet looks up multiple rows given the keys.
	MultiGet(keys []interface{}) (Result, error)
}

// Result reads Query results into Go objects, internally buffering them.
type Result interface {

	// Next sets destination with the contents of the current object in the
	// Result buffer and advances to the next object. It returns Done when
	// no more objects are available.
	Next(destination interface{}) error
}

type query struct {
	pool         *connectionPool
	mapping      Mapping
	reader       Reader
	columnLimit  int
	rowLimit     int
	reversed     bool
	components   []interface{}
	betweenStart interface{}
	betweenEnd   interface{}
}

func newQuery(cp *connectionPool, m Mapping) *query {
	return &query{
		pool:        cp,
		mapping:     m,
		reader:      cp.Reader().Cf(m.Cf()),
		columnLimit: DEFAULT_COLUMN_LIMIT,
		rowLimit:    DEFAULT_ROW_LIMIT,
		components:  make([]interface{}, 0),
	}
}

func (q *query) ConsistencyLevel(c ConsistencyLevel) Query {
	q.reader.ConsistencyLevel(c)
	return q
}

func (q *query) Limit(columns, rows int) Query {
	q.columnLimit = columns
	q.rowLimit = rows
	return q
}

func (q *query) Reversed(r bool) Query {
	q.reversed = r
	return q
}

func (q *query) Components(components ...interface{}) Query {
	q.components = components
	return q
}

func (q *query) Between(start, end interface{}) Query {
	q.betweenStart = start
	q.betweenEnd = end
	return q
}

func (q *query) Get(key interface{}) (Result, error) {
	return q.MultiGet([]interface{}{key})
}

func (q *query) GetOne(key interface{}, destination interface{}) error {
	res, err := q.Get(key)
	if err != nil {
		return err
	}
	return res.Next(destination)
}

var rangeOne = &Range{Count: 1}

func (q *query) RangeOne(destination interface{}) error {
	res, err := q.RangeGet(rangeOne)
	if err != nil {
		return err
	}
	return res.Next(destination)
}

func (q *query) Where(field string, op Operator, val interface{}) Query {
	bv, err := q.mapping.MarshalField(field, val)
	if err != nil {
		panic(err)
	}
	q.reader.Where([]byte(field), op, bv)
	return q
}

func (q *query) MultiGet(keys []interface{}) (Result, error) {
	var err error

	keysB := make([][]byte, 0)

	for _, key := range keys {
		keyB, err := q.mapping.MarshalKey(key)
		if err != nil {
			return nil, err
		}
		keysB = append(keysB, keyB)
	}

	q.buildSlice(q.reader)

	var rows []*Row

	if len(keysB) == 1 {
		row, err := q.reader.Get(keysB[0])
		if err != nil {
			return nil, err
		}
		if row != nil {
			rows = []*Row{row}
		}
	} else {
		rows, err = q.reader.MultiGet(keysB)
		if err != nil {
			return nil, err
		}
	}

	return &result{query: *q, buffer: rows}, nil
}

func (q *query) RangeGet(r *Range) (Result, error) {
	q.buildSlice(q.reader)
	rows, err := q.reader.RangeGet(r)
	if err != nil {
		return nil, err
	}
	return &result{query: *q, buffer: rows}, nil
}

func (q *query) buildSlice(reader Reader) error {
	var start, end []byte

	components := q.components

	if q.mapping.Compact() && len(q.mapping.Components()) == 1 {
		if len(components) == 1 {
			c := components[0]
			b, err := q.mapping.MarshalComponent(c, 0)
			if err != nil {
				return err
			}
			start = b
			end = b
		} else {
			if q.betweenStart != nil {
				b, err := q.mapping.MarshalComponent(q.betweenStart, 0)
				if err != nil {
					return err
				}
				start = b
			}
			if q.betweenEnd != nil {
				b, err := q.mapping.MarshalComponent(q.betweenEnd, 0)
				if err != nil {
					return err
				}
				end = b
			}
		}
	} else if len(components) > 0 {
		last := len(components) - 1
		for i, c := range components {
			b, err := q.mapping.MarshalComponent(c, i)
			if err != nil {
				return err
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
		if q.betweenStart != nil {
			b, err := q.mapping.MarshalComponent(q.betweenStart, len(components))
			if err != nil {
				return err
			}
			start = append(end, packComposite(b, eocEquals)...)
		}
		if q.betweenEnd != nil {
			b, err := q.mapping.MarshalComponent(q.betweenEnd, len(components))
			if err != nil {
				return err
			}
			end = append(end, packComposite(b, eocEquals)...)
		}
	}

	reader.Slice(&Slice{Start: start, End: end, Count: q.columnLimit, Reversed: q.reversed})
	return nil
}

type result struct {
	query
	buffer   []*Row
	row      *Row
	position int
}

func (r *result) feedRow() error {
	if r.row == nil {
		if len(r.buffer) == 0 {
			return Done
		}
		r.row = r.buffer[0]
		r.position = 0
		r.buffer = r.buffer[1:]
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

type StreamingResult struct {
	Mapping
	RowsChannel <-chan *Row
	row         *Row
	position    int
}

func (r *StreamingResult) feedRow() error {
	if r.row == nil {
		var ok bool
		r.row, ok = <-r.RowsChannel
		if !ok {
			return Done
		}
		r.position = 0
	}
	return nil
}

func (r *StreamingResult) Key() ([]byte, error) {
	if err := r.feedRow(); err != nil {
		return nil, err
	}
	return r.row.Key, nil
}

func (r *StreamingResult) NextColumn() (*Column, error) {
	if err := r.feedRow(); err != nil {
		return nil, err
	}
	if r.position >= len(r.row.Columns) {
		return nil, EndBeforeLimit
	}
	c := r.row.Columns[r.position]
	r.position++
	return c, nil
}

func (r *StreamingResult) Rewind() {
	r.position--
	if r.position < 0 {
		r.position = 0
	}
}

func (r *StreamingResult) Next(destination interface{}) error {
	err := r.Unmap(destination, r)
	if err == Done {
		// force new row feed and try again, just once
		r.row = nil
		err = r.Unmap(destination, r)
	}
	return err
}
