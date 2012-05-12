package gossie

import (
	"errors"
	"github.com/carloscm/gossie/src/cassandra"
	"github.com/pomack/thrift4go/lib/go/src/thrift"
)

/*
	to do:
	support Where for RangeGet in Cassandra 1.1
	figure out what's the deal with get_paged_slice in 1.1 and try to implement it in a sane way
*/

// Columns encapsulate the individual columns from/to Cassandra reads and writes
type Column struct {
	Name      []byte
	Value     []byte
	Ttl       int32
	Timestamp int64
}

// Row is a Cassandra row, including its row key
type Row struct {
	Key     []byte
	Columns []*Column
}

// RowColumnCount stores the number of columns matched in a MultiCount reader
type RowColumnCount struct {
	Key   []byte
	Count int
}

// Slice allows to specify a range of columns to return
// Always specify a Count value since there is an interface-mandated default of 100.
type Slice struct {
	Start    []byte
	End      []byte
	Count    int
	Reversed bool
}

// Range represents a range of rows to return, in order to be able to iterate over their keys.
// The low level token range is not exposed. Use an empty slice to indicate if you want the first
// or the last possible key in a range then pass the last read row key as the new Start key in a
// new RangeGet reader to page results. This will allow you to iterate over an entire CF even when
// using the random partitioner. Always specify a Count value since there is an interface-mandated
// default of 100.
type Range struct {
	Start []byte
	End   []byte
	Count int
}

// IndexedRange represents a range of rows to return for the IndexedGet method.
// The low level token range is not exposed. Use an empty slice to indicate if you want the first key
// in a range, then pass the last read row key as the new Start key in a new IndexedGet reader to page
// results. Always specify a Count value since there is an interface-mandated default of 100.
type IndexedRange struct {
	Start []byte
	Count int
}

// Operator for Where
type Operator int

const (
	EQ  Operator = 0
	GTE Operator = 1
	GT  Operator = 2
	LTE Operator = 3
	LT  Operator = 4
)

// Reader is the interface for all read operations over Cassandra.
// The method calls support chaining so you can build concise queries
type Reader interface {

	// ConsistencyLevel sets the consistency level for this particular call.
	// It is optional, if left uncalled it will default to your connection pool options value.
	ConsistencyLevel(int) Reader

	// Cf sets the column family name for the reader.
	// This method must be always called.
	Cf(string) Reader

	// Slice optionally sets a slice to set a range of column names and potentially iterate over the
	// columns of the returned row(s)
	Slice(*Slice) Reader

	// Columns optionally filters the returned columns to only the passed set of column names
	Columns([][]byte) Reader

	// Each call to this method adds a new comparison to be checked against the returned rows of
	// IndexedGet
	// All the comparisons are checked for every row. In the current Cassandra implementation at
	// least one of the Where calls must use a secondary indexed column with an EQ operator.
	Where(column []byte, op Operator, value []byte) Reader

	// Get looks up a row with the given key and returns it, or nil in case it is not found
	Get(key []byte) (*Row, error)

	// MultiGet performs a parallel Get operation for all the passed keys, and returns a slice of
	// RowColumnCounts pointers to the gathered rows, which may be empty if none were found. It returns
	// nil only on error conditions
	MultiGet(keys [][]byte) ([]*Row, error)

	// Count looks up a row with the given key and returns the number of columns it has
	Count(key []byte) (int, error)

	// MultiGet performs a parallel Count operation for all the passed keys, and returns a slice of Row
	// pointers to the gathered rows, which may be empty if none were found. It returns nil only on
	// error conditions
	MultiCount(keys [][]byte) ([]*RowColumnCount, error)

	// RangeGet performs a sequential Get operation for a range of rows. See the docs for Range for an
	// explanation on how to page results. It returns a slice of Row pointers to the gathered rows, which
	// may be empty if none were found. It returns nil only on error conditions
	RangeGet(*Range) ([]*Row, error)

	// IndexedGet performs a sequential Get operation for a range of rows and returns only those that match
	// the Where clauses. See the docs for Range for an explanation on how to page results. It returns a
	// slice of Row pointers to the gathered rows, which may be empty if none were found. It returns nil only
	// on error conditions
	IndexedGet(*IndexedRange) ([]*Row, error)
}

type reader struct {
	pool             *connectionPool
	consistencyLevel int
	cf               string
	slice            Slice
	setSlice         bool
	columns          [][]byte
	setColumns       bool
	setWhere         bool
	expressions      thrift.TList
}

func newReader(cp *connectionPool, cl int) *reader {
	return &reader{
		pool:             cp,
		consistencyLevel: cl,
	}
}

func (r *reader) ConsistencyLevel(l int) Reader {
	r.consistencyLevel = l
	return r
}

func (r *reader) Cf(cf string) Reader {
	r.cf = cf
	return r
}

func (r *reader) Slice(s *Slice) Reader {
	r.slice = *s
	r.setSlice = true
	return r
}

func (r *reader) Columns(c [][]byte) Reader {
	copy(r.columns, c)
	r.setColumns = true
	return r
}

func (r *reader) Where(column []byte, op Operator, value []byte) Reader {
	if r.expressions == nil {
		r.expressions = thrift.NewTList(thrift.STRUCT, 1)
	}
	exp := cassandra.NewIndexExpression()
	exp.ColumnName = column
	exp.Op = cassandra.IndexOperator(op)
	exp.Value = value
	r.expressions.Push(exp)
	r.setWhere = true
	return r
}

func sliceToCassandra(slice *Slice) *cassandra.SliceRange {
	sr := cassandra.NewSliceRange()
	sr.Start = slice.Start
	sr.Finish = slice.End
	sr.Count = int32(slice.Count)
	sr.Reversed = slice.Reversed
	// workaround some uninitialized slice == nil quirks that trickle down into the generated thrift4go code
	if sr.Start == nil {
		sr.Start = make([]byte, 0)
	}
	if sr.Finish == nil {
		sr.Finish = make([]byte, 0)
	}
	return sr
}

func fullSlice() *cassandra.SliceRange {
	sr := cassandra.NewSliceRange()
	// workaround some uninitialized slice == nil quirks that trickle down into the generated thrift4go code
	sr.Start = make([]byte, 0)
	sr.Finish = make([]byte, 0)
	return sr
}

func (r *reader) buildPredicate() *cassandra.SlicePredicate {
	sp := cassandra.NewSlicePredicate()
	if r.setColumns {
		sp.ColumnNames = thrift.NewTList(thrift.BINARY, 1)
		for _, col := range r.columns {
			sp.ColumnNames.Push(col)
		}
	} else if r.setSlice {
		sp.SliceRange = sliceToCassandra(&r.slice)
	} else {
		sp.SliceRange = fullSlice()
	}
	return sp
}

func (r *reader) buildColumnParent() *cassandra.ColumnParent {
	cp := cassandra.NewColumnParent()
	cp.ColumnFamily = r.cf
	return cp
}

func (q *reader) buildKeyRange(r *Range) *cassandra.KeyRange {
	kr := cassandra.NewKeyRange()
	kr.StartKey = r.Start
	kr.EndKey = r.End
	kr.Count = int32(r.Count)
	// workaround some uninitialized slice == nil quirks that trickle down into the generated thrift4go code
	if kr.StartKey == nil {
		kr.StartKey = make([]byte, 0)
	}
	if kr.EndKey == nil {
		kr.EndKey = make([]byte, 0)
	}
	return kr
}

func (r *reader) buildIndexClause(ir *IndexedRange) *cassandra.IndexClause {
	ic := cassandra.NewIndexClause()
	ic.Expressions = r.expressions
	ic.StartKey = ir.Start
	ic.Count = int32(ir.Count)
	// workaround some uninitialized slice == nil quirks that trickle down into the generated thrift4go code
	if ic.StartKey == nil {
		ic.StartKey = make([]byte, 0)
	}
	return ic
}

func (r *reader) Get(key []byte) (*Row, error) {
	if r.cf == "" {
		return nil, errors.New("No column family specified")
	}

	cp := r.buildColumnParent()
	sp := r.buildPredicate()

	var ret thrift.TList
	err := r.pool.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, error) {
		var ire *cassandra.InvalidRequestException
		var ue *cassandra.UnavailableException
		var te *cassandra.TimedOutException
		var err error
		ret, ire, ue, te, err = c.client.GetSlice(key, cp, sp, cassandra.ConsistencyLevel(r.consistencyLevel))
		return ire, ue, te, err
	})

	if err != nil {
		return nil, err
	}

	return rowFromTListColumns(key, ret), nil
}

func (r *reader) Count(key []byte) (int, error) {
	if r.cf == "" {
		return 0, errors.New("No column family specified")
	}

	cp := r.buildColumnParent()
	sp := r.buildPredicate()

	var ret int32
	err := r.pool.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, error) {
		var ire *cassandra.InvalidRequestException
		var ue *cassandra.UnavailableException
		var te *cassandra.TimedOutException
		var err error
		ret, ire, ue, te, err = c.client.GetCount(key, cp, sp, cassandra.ConsistencyLevel(r.consistencyLevel))
		return ire, ue, te, err
	})

	if err != nil {
		return 0, err
	}

	return int(ret), nil
}

func (r *reader) buildMultiKeys(keys [][]byte) thrift.TList {
	tkeys := thrift.NewTList(thrift.BINARY, 1)
	for _, k := range keys {
		tkeys.Push(k)
	}
	return tkeys
}

func (r *reader) MultiGet(keys [][]byte) ([]*Row, error) {
	if r.cf == "" {
		return nil, errors.New("No column family specified")
	}

	if len(keys) <= 0 {
		return make([]*Row, 0), nil
	}

	cp := r.buildColumnParent()
	sp := r.buildPredicate()
	tk := r.buildMultiKeys(keys)

	var ret thrift.TMap
	err := r.pool.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, error) {
		var ire *cassandra.InvalidRequestException
		var ue *cassandra.UnavailableException
		var te *cassandra.TimedOutException
		var err error
		ret, ire, ue, te, err = c.client.MultigetSlice(tk, cp, sp, cassandra.ConsistencyLevel(r.consistencyLevel))
		return ire, ue, te, err
	})

	if err != nil {
		return nil, err
	}

	return rowsFromTMap(ret), nil
}

func (r *reader) MultiCount(keys [][]byte) ([]*RowColumnCount, error) {
	if r.cf == "" {
		return nil, errors.New("No column family specified")
	}

	if len(keys) <= 0 {
		return make([]*RowColumnCount, 0), nil
	}

	cp := r.buildColumnParent()
	sp := r.buildPredicate()
	tk := r.buildMultiKeys(keys)

	var ret thrift.TMap
	err := r.pool.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, error) {
		var ire *cassandra.InvalidRequestException
		var ue *cassandra.UnavailableException
		var te *cassandra.TimedOutException
		var err error
		ret, ire, ue, te, err = c.client.MultigetCount(tk, cp, sp, cassandra.ConsistencyLevel(r.consistencyLevel))
		return ire, ue, te, err
	})

	if err != nil {
		return nil, err
	}

	return rowsColumnCountFromTMap(ret), nil
}

func (r *reader) RangeGet(rang *Range) ([]*Row, error) {
	if r.cf == "" {
		return nil, errors.New("No column family specified")
	}

	if rang == nil || rang.Count <= 0 {
		return make([]*Row, 0), nil
	}

	kr := r.buildKeyRange(rang)
	cp := r.buildColumnParent()
	sp := r.buildPredicate()

	var ret thrift.TList
	err := r.pool.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, error) {
		var ire *cassandra.InvalidRequestException
		var ue *cassandra.UnavailableException
		var te *cassandra.TimedOutException
		var err error
		ret, ire, ue, te, err = c.client.GetRangeSlices(cp, sp, kr, cassandra.ConsistencyLevel(r.consistencyLevel))
		return ire, ue, te, err
	})

	if err != nil {
		return nil, err
	}

	return rowsFromTListKeySlice(ret), nil
}

func (r *reader) IndexedGet(rang *IndexedRange) ([]*Row, error) {
	if r.cf == "" {
		return nil, errors.New("No column family specified")
	}

	if !r.setWhere {
		return nil, errors.New("At least one Where call must be made")
	}

	if rang == nil || rang.Count <= 0 {
		return make([]*Row, 0), nil
	}

	ic := r.buildIndexClause(rang)
	cp := r.buildColumnParent()
	sp := r.buildPredicate()

	var ret thrift.TList
	err := r.pool.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, error) {
		var ire *cassandra.InvalidRequestException
		var ue *cassandra.UnavailableException
		var te *cassandra.TimedOutException
		var err error
		ret, ire, ue, te, err = c.client.GetIndexedSlices(cp, ic, sp, cassandra.ConsistencyLevel(r.consistencyLevel))
		return ire, ue, te, err
	})

	if err != nil {
		return nil, err
	}

	return rowsFromTListKeySlice(ret), nil
}

func rowFromTListColumns(key []byte, tl thrift.TList) *Row {
	if tl == nil || tl.Len() <= 0 {
		return nil
	}
	r := &Row{Key: key}
	for colI := range tl.Iter() {
		var col *cassandra.ColumnOrSuperColumn = colI.(*cassandra.ColumnOrSuperColumn)
		if col.Column != nil {
			c := &Column{
				Name:      col.Column.Name,
				Value:     col.Column.Value,
				Timestamp: col.Column.Timestamp,
				Ttl:       col.Column.Ttl,
			}
			r.Columns = append(r.Columns, c)
		} else if col.CounterColumn != nil {
			v, _ := Marshal(col.CounterColumn.Value, LongType)
			c := &Column{
				Name:  col.CounterColumn.Name,
				Value: v,
			}
			r.Columns = append(r.Columns, c)
		}
	}
	return r
}

func keyFromTMap(e thrift.TMapElem) []byte {
	// workaround some issues with the way the key->row array gets built by thrift4go and
	// the cassandra IDL wrongly insisting keys are strings
	rawKey := e.Key()
	var key []byte
	switch k := rawKey.(type) {
	case []uint8:
		key = []byte(k)
	case string:
		key = []byte(k)
	}
	return key
}

func rowsFromTMap(tm thrift.TMap) []*Row {
	if tm == nil || tm.Len() <= 0 {
		return make([]*Row, 0)
	}
	r := make([]*Row, 0)
	for rowI := range tm.Iter() {
		key := keyFromTMap(rowI)
		columns := (rowI.Value()).(thrift.TList)
		row := rowFromTListColumns(key, columns)
		if row != nil {
			r = append(r, row)
		}
	}
	return r
}

func rowsColumnCountFromTMap(tm thrift.TMap) []*RowColumnCount {
	if tm == nil || tm.Len() <= 0 {
		return make([]*RowColumnCount, 0)
	}
	r := make([]*RowColumnCount, 0)
	for rowI := range tm.Iter() {
		key := keyFromTMap(rowI)
		count := int((rowI.Value()).(int32))
		if count > 0 {
			r = append(r, &RowColumnCount{Key: key, Count: count})
		}
	}
	return r
}

func rowsFromTListKeySlice(tl thrift.TList) []*Row {
	if tl == nil || tl.Len() <= 0 {
		return make([]*Row, 0)
	}
	r := make([]*Row, 0)
	for keySliceI := range tl.Iter() {
		keySlice := keySliceI.(*cassandra.KeySlice)
		key := keySlice.Key
		row := rowFromTListColumns(key, keySlice.Columns)
		if row != nil {
			r = append(r, row)
		}
	}
	return r
}
