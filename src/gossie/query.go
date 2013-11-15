package gossie

import (
	"errors"
	"github.com/apesternikov/gossie/src/cassandra"
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
	ConsistencyLevel(cassandra.ConsistencyLevel) Query

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
	pool             *connectionPool
	mapping          Mapping
	consistencyLevel cassandra.ConsistencyLevel
	columnLimit      int
	rowLimit         int
	reversed         bool
	components       []interface{}
	betweenStart     interface{}
	betweenEnd       interface{}
}

func newQuery(cp *connectionPool, m Mapping) *query {
	return &query{
		pool:        cp,
		mapping:     m,
		columnLimit: DEFAULT_COLUMN_LIMIT,
		rowLimit:    DEFAULT_ROW_LIMIT,
		components:  make([]interface{}, 0),
	}
}

func (q *query) ConsistencyLevel(c cassandra.ConsistencyLevel) Query {
	q.consistencyLevel = c
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

	reader := q.pool.Reader().Cf(q.mapping.Cf())

	if q.consistencyLevel != 0 {
		reader.ConsistencyLevel(q.consistencyLevel)
	}

	q.buildSlice(reader)

	rows := make([]*Row, 0)

	if len(keysB) == 1 {
		row, err := reader.Get(keysB[0])
		if err != nil {
			return nil, err
		}
		if row != nil {
			rows = []*Row{row}
		}
	} else {
		rows, err = reader.MultiGet(keysB)
		if err != nil {
			return nil, err
		}
	}

	return &result{query: *q, buffer: rows}, nil
}

func (q *query) buildSlice(reader Reader) error {
	start := make([]byte, 0)
	end := make([]byte, 0)

	components := q.components
	if q.betweenStart != nil {
		components = append(components, q.betweenStart)
	}

	if len(components) > 0 {
		last := len(components) - 1
		for i, c := range components {
			b, err := q.mapping.MarshalComponent(c, i)
			if err != nil {
				return err
			}
			start = append(start, packComposite(b, eocEquals)...)
			if i == last {
				if q.betweenEnd != nil {
					b, err := q.mapping.MarshalComponent(q.betweenEnd, i)
					if err != nil {
						return err
					}
					end = append(end, packComposite(b, eocEquals)...)
				} else {
					end = append(end, packComposite(b, eocGreater)...)
				}
			} else {
				end = append(end, packComposite(b, eocEquals)...)
			}
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
