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

	row, ms, err := internalMap(source)
	if err != nil {
		return err
	}

	if err = c.pool.Mutation().Insert(ms.sm.cf, row).Run(); err != nil {
		return err
	}

	return nil
}

func (c *cursor) Read(source interface{}) error {

	// deconstruct the source struct into a reflect.Value and a (cached) struct mapping
	ms, err := newMappedStruct(source)
	if err != nil {
		return err
	}

	// sanity checks
	if ms.sm.isSliceColumn {
		return errors.New(fmt.Sprint("Slice field in col tag is unsuported in Cursor for now, check back soon!"))
	}

	// marshal the key field
	key, err := ms.marshalKey()
	if err != nil {
		return err
	}

	// start building the query
	q := c.pool.Query().Cf(ms.sm.cf)

	// build a slice composite comparator if needed
	if ms.sm.isCompositeColumn {
		// iterate over the components and set an equality comparison for every simple field
		start := make([]byte, 0)
		end := make([]byte, 0)
		var component int
		for component = 0; component < len(ms.sm.columns); component++ {
			fm := ms.sm.columns[component]
			if fm.fieldKind != baseTypeField {
				break
			}
			b, err := ms.mapColumn(baseTypeField, fm, 0)
			if err != nil {
				return err
			}
			start = packComposite(start, b, eocEquals)
			end = packComposite(end, b, eocGreater)
		}

		/*if component < len(ms.sm.columns) {
		    // we still got one to go, this means the last one was an iterable non-fixed type (*name or go slice)
		    //fm := ms.sm.columns[component]
		    // TODO: this will only work for *name
		    b := make([]byte, 0)
		    start = packComposite(start, b, o[6], o[7], o[8])
		    end = packComposite(end, b, o[9], o[10], o[11])
		}*/

		// TODO: fix hardcoded number of columns
		q.Slice(&Slice{Start: start, End: end, Count: 100})
	}

	//isCompositeColumn bool
	//isSliceColumn bool
	//isStarNameColumn bool

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
