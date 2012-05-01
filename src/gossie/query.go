package gossie

import (
	"errors"
)

/*

todo:

    buffering for slicing and range get
    Next/Prev for the buffering with autopaging

    Search() and interface(s) for indexed get

    multiget
        still need to think abstraction
        or just with a direct call that accepts interface{} and assumes Go slices? what about composites?

*/

var (
	Done = errors.New("No more results found")
)

const (
	DEFAULT_LIMIT = 100
)

// Query is a high level interface for Cassandra queries
type Query interface {

	// ConsistencyLevel sets the consistency level for this particular call.
	// It is optional, if left uncalled it will default to your connection
	// pool options value.
	ConsistencyLevel(int) Query

	// Limit sets the column limit to buffer at once. Paging with a Result
	// will internally query Cassandra in Limit sized column slices.
	Limit(int) Query

	// Get looks up a row with the given key (and optionally components for a
	// composite) and returns a Result to it. If the row uses a composite
	// comparator an you only specify the key the Result will allow you to
	// page and iterate over the entire row.
	Get(key interface{}, components ...interface{}) (Result, error)
}

// Result reads Query results into Go objects, internally buffering them and
// paging them.
type Result interface {

	// Next sets destination with the contents of the current object in the
	// Result buffer and advances to the next object one. It returns Done when
	// no more objects are available.
	Next(destination interface{}) error
}

type query struct {
	pool             *connectionPool
	mapping          Mapping
	consistencyLevel int
	limit            int
	offset           int
	row              *Row
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

func (q *query) Get(key interface{}, components ...interface{}) (Result, error) {
	return nil, nil
	/*
		// sanity checks
		// marshal the key field
		vk := reflect.Value(key)


		marshal especifico para VALUE de una key arbitraria, no intentar pillar un campo de struct


		key, err := q.mapping.key.marshalValue(&vk)
		if err != nil {
			return err
		}

		// start building the query
		q := c.pool.Reader().Cf(mi.m.cf)

		if c.options.ReadConsistency != 0 {
			q.ConsistencyLevel(c.options.ReadConsistency)
		}

		// build a slice composite comparator if needed
		if len(mi.m.composite) > 0 {
			// iterate over the components and set an equality comparison for every field
			start := make([]byte, 0)
			end := make([]byte, 0)
			for _, f := range mi.m.composite {
				b, err := f.marshalValue(&mi.v)
				if err != nil {
					return err
				}
				start = append(start, packComposite(b, eocEquals)...)
				end = append(end, packComposite(b, eocGreater)...)
			}
			q.Slice(&Slice{Start: start, End: end, Count: c.options.Limit})
		}

		row, err := q.Get(key)

		if err != nil {
			return err
		}

		if row == nil {
			return ErrorNotFound
		}

		return Unmap(row, source)
	*/
}
