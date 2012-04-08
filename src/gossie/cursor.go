package gossie

import (
	"errors"
	"fmt"
)

/*

todo:

    buffering for slicing and range get
    Next/Prev for the buffering with autopaging

    isSliceColumn == true

    Search() and interface(s) for indexed get

    multiget
        still need to think abstraction
        or just with a direct call that accepts interface{} and assumes Go slices? what about composites?

*/

// Cursor is a simple cursor-based interface for reading and writing structs from a Cassandra column family.
type Cursor interface {

	// Read rows (or slices of a row) to fill up the struct.
	//
	// For limit == 1 a single get operation will be issued. If a composite is present then this operation will
	// issue a slice with an exact match for the composite value
	//
	// For limit > 1, and for structs with no composite, a range get operation will be issued, the result buffered
	// internally in the Cursor, and the first returned row unmapped into the struct. Next() is then available for
	// paging inside the returned results and will issue new range get operations when neede. Please note that row
	// range get is unordered in most Cassandra configurations. TO DO
	//
	// If the struct has a composite column name then this operation will issue a single row get, but with a slice
	// predicate that allows for iteration over the row with Next()/Prev() TO DO
	Read(limit int) error

	//Next() bool
	//Prev() bool

	//First()

	// Write a single row (or slice of a row) with the data currently in the struct.
	Write() error

	//Delete()
}

type cursor struct {
	source interface{}
	pool   *connectionPool
	//position int
	//buffer []interface{}
}

func makeCursor(cp *connectionPool, source interface{}) (c *cursor) {
	return &cursor{
		source: source,
		pool:   cp,
	}
}

func (c *cursor) Write() error {

	row, ms, err := internalMap(c.source)
	if err != nil {
		return err
	}

	if err = c.pool.Mutation().Insert(ms.sm.cf, row).Run(); err != nil {
		return err
	}

	return nil
}

func (c *cursor) Read(limit int) error {
	if limit < 0 {
		return errors.New("Limit is less than 1, nothing to read")
	}

	// deconstruct the source struct into a reflect.Value and a (cached) struct mapping
	ms, err := newMappedStruct(c.source)
	if err != nil {
		return err
	}

	// sanity checks
	if !ms.sm.isStarNameColumn && !ms.sm.isSliceColumn {
		return errors.New(fmt.Sprint("Struct ", ms.v.Type().Name(), " has no *name nor slice field in its col tag, nothing to read"))
	}
	if ms.sm.isSliceColumn {
		return errors.New(fmt.Sprint("Slice field in col tag is unsuported in Cursor for now, check back soon!"))
	}

	// marshal the key field for the key to look up
	key, err := ms.marshalKey()
	if err != nil {
		return err
	}

	// start building the query
	q := c.pool.Query().Cf(ms.sm.cf)

	//isCompositeColumn bool
	//isSliceColumn bool
	//isStarNameColumn bool

	if limit == 1 {

		if ms.sm.isCompositeColumn {
			return errors.New(fmt.Sprint("Cursor composite support will be implemented soon!"))
			/*

			   // we only want a single result so issue an exact match composite comparator slice
			   start := make([]byte, 0)
			   start = packComposite(start, component []byte, true, true, true)

			   packComposite(start, component []byte, true, sliceStart, inclusive bool) []byte {

			   s := &Slice{
			       Start:
			       End:
			       Count: 1
			   }
			   q.Slice(s)
			*/

		}
		row, err := q.Get(key)
		if err != nil {
			return err
		}
		err = Unmap(row, c.source)
		if err != nil {
			return err
		}
	} else {
		return errors.New(fmt.Sprint("Limit > 1 will be implemented soon!"))
	}

	return nil
}
