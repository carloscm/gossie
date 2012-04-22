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
}

type cursor struct {
	pool *connectionPool
}

func makeCursor(cp *connectionPool) (c *cursor) {
	return &cursor{
		pool: cp,
	}
}

func (c *cursor) Write(source interface{}) error {

	row, mi, err := internalMap(source)
	if err != nil {
		return err
	}

	if err = c.pool.Mutation().Insert(mi.m.cf, row).Run(); err != nil {
		return err
	}

	return nil
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

	// build a slice composite comparator if needed
	if len(mi.m.composite) > 0 {
		// iterate over the components and set an equality comparison for every simple field
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
		// TODO: fix hardcoded number of columns
		q.Slice(&Slice{Start: start, End: end, Count: 100})
	}

	row, err := q.Get(key)

	if err != nil {
		return err
	}

	if row == nil {
		return ErrorNotFound
	}

	err = Unmap(row, source)
	if err != nil {
		return err
	}

	return nil
}
