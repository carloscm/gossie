package gossie

import (
	"errors"
)

/*

todo:

    autopaging?

    Search() and interface(s) for indexed get

    multiget
        still need to think abstraction
        or just with a direct call that accepts interface{} and assumes Go slices? what about composites?

*/

var (
	Done = errors.New("No more results found")
)

const (
	DEFAULT_LIMIT = 10000
)

// Query is a high level interface for Cassandra queries
type Query interface {

	// ConsistencyLevel sets the consistency level for this particular call.
	// It is optional, if left uncalled it will default to your connection
	// pool options value.
	ConsistencyLevel(int) Query

	// Limit sets the column limit to buffer at once.
	Limit(int) Query

	// Reverse set to true will reverse the order of the columns in the result.
	Reversed(bool) Query

	// Get looks up a row with the given key (and optionally components for a
	// composite) and returns a Result to it. If the row uses a composite
	// comparator and you only specify the key and zero or more comparator
	// components the Result will allow you to iterate over the entire row.
	Get(key interface{}, components ...interface{}) (Result, error)

	// GetBetween is like Get, but the last two passed components values are
	// used as the last values for the slice Start and End composite values.
	GetBetween(key interface{}, components ...interface{}) (Result, error)
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
	consistencyLevel int
	limit            int
	reversed         bool
}

func newQuery(cp *connectionPool, m Mapping) *query {
	return &query{
		pool:    cp,
		mapping: m,
		limit:   DEFAULT_LIMIT,
	}
}

func (q *query) ConsistencyLevel(c int) Query {
	q.consistencyLevel = c
	return q
}

func (q *query) Limit(l int) Query {
	q.limit = l
	return q
}

func (q *query) Reversed(r bool) Query {
	q.reversed = r
	return q
}

func (q *query) Get(key interface{}, components ...interface{}) (Result, error) {
	keyB, err := q.mapping.MarshalKey(key)
	if err != nil {
		return nil, err
	}

	reader := q.pool.Reader().Cf(q.mapping.Cf())

	if q.consistencyLevel != 0 {
		reader.ConsistencyLevel(q.consistencyLevel)
	}

	return &result{*q, keyB, reader, components, nil, 0, nil}, nil
}

func (q *query) GetBetween(key interface{}, components ...interface{}) (Result, error) {
	if len(components) < 2 {
		return nil, errors.New("GetBetween requires at least 2 component values")
	}
	r, err := q.Get(key, components[:len(components)-1]...)
	if err != nil {
		return nil, err
	}
	r.(*result).between = components[len(components)-1]
	return r, nil
}

type result struct {
	query
	key        []byte
	reader     Reader
	components []interface{}
	row        *Row
	position   int
	between    interface{}
}

func (r *result) Key() []byte {
	return r.key
}

func (r *result) buildFixedSlice() error {
	start := make([]byte, 0)
	end := make([]byte, 0)
	if len(r.components) > 0 {
		last := len(r.components) - 1
		for i, c := range r.components {
			b, err := r.mapping.MarshalComponent(c, i)
			if err != nil {
				return err
			}
			start = append(start, packComposite(b, eocEquals)...)
			if i == last {
				if r.between != nil {
					b, err := r.mapping.MarshalComponent(r.between, i)
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
	r.reader.Slice(&Slice{Start: start, End: end, Count: r.limit, Reversed: r.reversed})
	return nil
}

func (r *result) NextColumn() (*Column, error) {
	if r.row == nil {
		err := r.buildFixedSlice()
		if err != nil {
			return nil, err
		}
		row, err := r.reader.Get(r.key)
		if err != nil {
			return nil, err
		}
		if row == nil {
			return nil, Done
		}
		r.row = row
		r.position = 0
	}
	if r.position >= len(r.row.Columns) {
		if r.position >= r.limit {
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
	return r.mapping.Unmap(destination, r)
}
