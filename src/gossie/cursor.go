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
	ErrorNotFound = errors.New("Row (plus any composites) not found")
)

const (
	DEFAULT_LIMIT_COLUMNS = 100
)

// Cursor is a simple cursor-based interface for reading and writing structs from a Cassandra column family.
type Cursor interface {

	//Search(...)
	//Next() bool
	//Prev() bool
	//First()

	// Read a single row (or slice of a row) to fill up the struct.
	// A single get operation will be issued. If a composite is present then this operation will issue a slice with
	// an exact match for the composite value with the field values present in the struct
	Read(interface{}) error

	// Write a single row (or slice of a row) with the data currently in the struct.
	Write(interface{}) error

	//Delete()

	// Options sets the options for this Cursor
	Options(CursorOptions)
}

// CursorOptions stores some options that modify the behaviour of the queries a Cursor performs
type CursorOptions struct {

	// LimitColumns is the max number of columns that will be read from Cassandra. Default is 100.
	LimitColumns int

	// WriteConsistency overrides the default write consistency level for the underlying connection pool.
	// Default is 0 which means to not override the consistenct level.
	WriteConsistency int

	// ReadConsistency overrides the default read consistency level for the underlying connection pool.
	// Default is 0 which means to not override the consistenct level.
	ReadConsistency int
}

type cursor struct {
	pool    *connectionPool
	options CursorOptions
}

func (co *CursorOptions) defaults() {
	if co.LimitColumns == 0 {
		co.LimitColumns = DEFAULT_LIMIT_COLUMNS
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

	m := c.pool.Mutation().Insert(mi.m.cf, row)

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
	q := c.pool.Query().Cf(mi.m.cf)

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
		q.Slice(&Slice{Start: start, End: end, Count: c.options.LimitColumns})
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
