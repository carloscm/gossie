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
	ErrorNotFound = errors.New("Key (plus any composites) not found")
)

const (
	DEFAULT_LIMIT = 100
)

/*
// Query is a high level interface for Cassandra queries
type Query interface {

	// ConsistencyLevel sets the consistency level for this particular call.
	// It is optional, if left uncalled it will default to your connection pool options value.
	ConsistencyLevel(int) Query

	// Limit sets the object limit to buffer at once. Paging with a Result will internally buffer up to
	// Limit objects. Defaults to 100.
	Limit(int) Query

	// Get looks up a row with the given key (and optionally components for a composite) and returns
	// a Result to it. If the row uses a composite comparator an you only specify the key the Cursor
	// will allow you page and iterate over the entire row.
	Get(key interface{}, components interface{}...) (Result, error)
}

// Result reads Query results into Go objects, internally buffering them and paging them.
type Result interface {

	// Next advances to the next object. It returns false when no more objects are available
	Next() bool

	// Read the current object in the result sequence.
	Read(interface{}) error
}

type cursor struct {
	pool    *connectionPool
	options CursorOptions
}

func (co *CursorOptions) defaults() {
	if co.Limit == 0 {
		co.Limit = DEFAULT_LIMIT
	}
}

func newCursor(cp *connectionPool) *cursor {
	c := &cursor{
		pool: cp,
	}
	c.options.defaults()
	return c
}

func (c *cursor) Options(options CursorOptions) {
	options.defaults()
	c.options = options
}

func (c *cursor) Write(source interface{}) error {

	row, mi, err := internalMap(source)
	if err != nil {
		return err
	}

	m := c.pool.Writer().Insert(mi.m.cf, row)

	if c.options.WriteConsistency != 0 {
		m.ConsistencyLevel(c.options.WriteConsistency)
	}

	return m.Run()
}

func (c *cursor) Read(source interface{}) error {

	// deconstruct the source struct into a reflect.Value and a (cached) struct mapping
	mi, err := newMappedInstance(source)
	if err != nil {
		return err
	}

	// sanity checks
	// marshal the key field
	key, err := mi.m.key.marshalValue(&mi.v)
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
}
*/
