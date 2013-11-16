package gossie

import (
	"bytes"
	"errors"
	. "github.com/apesternikov/gossie/src/cassandra"
	"github.com/golang/glog"
)

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
type Operator IndexOperator

const (
	EQ  Operator = Operator(IndexOperator_EQ)
	GTE Operator = Operator(IndexOperator_GTE)
	GT  Operator = Operator(IndexOperator_GT)
	LTE Operator = Operator(IndexOperator_LTE)
	LT  Operator = Operator(IndexOperator_LT)
)

// Reader is the interface for all read operations over Cassandra.
// The method calls support chaining so you can build concise queries
type Reader interface {

	// ConsistencyLevel sets the consistency level for this particular call.
	// It is optional, if left uncalled it will default to your connection pool options value.
	ConsistencyLevel(ConsistencyLevel) Reader

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
	// Deprecated, please use RangeGet with Where() instead
	IndexedGet(*IndexedRange) ([]*Row, error)

	//Set range to use with RangeScan
	//Default token range is -1 to 170141183460469231731687303715884105728
	SetTokenRange(startToken, endToken string) Reader

	// Scan a range This function will call the callback
	// function with data read. Callback should return true to continue scanning or false to stop
	// Usage:
	// pool.Reader().SetTokenRange("1234", "4567").RangeScan(func(func(r *Row) bool)
	RangeScan() (data <-chan *Row, err <-chan error)

	//WideRowScan performs sequential scan for a range of columns in a single row. It will call the callback
	// function with data read. Callback should return true to continue scanning or false to stop
	WideRowScan(key, startColumn []byte, batchSize int32, callback func(*Column) bool) error
}

const (
	DEF_START_TOKEN = "-1"
	DEF_END_TOKEN   = "170141183460469231731687303715884105728"
)

type reader struct {
	pool             *connectionPool
	consistencyLevel ConsistencyLevel
	cf               string
	slice            Slice
	setSlice         bool
	columns          [][]byte
	setColumns       bool
	expressions      []*IndexExpression
	startToken       string
	endToken         string
}

func newReader(cp *connectionPool, cl ConsistencyLevel) *reader {
	return &reader{
		pool:             cp,
		consistencyLevel: cl,
	}
}

func (r *reader) SetTokenRange(startToken, endToken string) Reader {
	r.startToken, r.endToken = startToken, endToken
	return r
}

func (r *reader) ConsistencyLevel(l ConsistencyLevel) Reader {
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
	r.columns = make([][]byte, len(c))
	copy(r.columns, c)
	r.setColumns = true
	return r
}

func (r *reader) Where(column []byte, op Operator, value []byte) Reader {
	exp := NewIndexExpression()
	exp.ColumnName = column
	exp.Op = IndexOperator(op)
	exp.Value = value
	r.expressions = append(r.expressions, exp)
	return r
}

func sliceToCassandra(slice *Slice) *SliceRange {
	sr := NewSliceRange()
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

func fullSlice() *SliceRange {
	sr := NewSliceRange()
	// workaround some uninitialized slice == nil quirks that trickle down into the generated thrift4go code
	sr.Start = make([]byte, 0)
	sr.Finish = make([]byte, 0)
	return sr
}

func (r *reader) buildPredicate() *SlicePredicate {
	sp := NewSlicePredicate()
	if r.setColumns {
		sp.ColumnNames = r.columns
	} else if r.setSlice {
		sp.SliceRange = sliceToCassandra(&r.slice)
	} else {
		sp.SliceRange = fullSlice()
	}
	return sp
}

func (r *reader) buildColumnParent() *ColumnParent {
	cp := NewColumnParent()
	cp.ColumnFamily = r.cf
	return cp
}

func (q *reader) buildKeyRange(r *Range) *KeyRange {
	kr := NewKeyRange()
	kr.StartKey = r.Start
	kr.EndKey = r.End
	kr.Count = int32(r.Count)
	// workaround some uninitialized slice == nil quirks that trickle down into the generated thrift4go code
	//TODO: this is no-op, because zero slices are exactly that: empty slices
	if kr.StartKey == nil {
		kr.StartKey = make([]byte, 0)
	}
	if kr.EndKey == nil {
		kr.EndKey = make([]byte, 0)
	}
	kr.RowFilter = q.expressions
	return kr
}

func (r *reader) buildIndexClause(ir *IndexedRange) *IndexClause {
	ic := NewIndexClause()
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

	var ret []*ColumnOrSuperColumn
	err := r.pool.run(func(c *connection) (*InvalidRequestException, *UnavailableException, *TimedOutException, error) {
		var ire *InvalidRequestException
		var ue *UnavailableException
		var te *TimedOutException
		var err error
		ret, ire, ue, te, err = c.client.GetSlice(key, cp, sp, r.consistencyLevel)
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
	err := r.pool.run(func(c *connection) (*InvalidRequestException, *UnavailableException, *TimedOutException, error) {
		var ire *InvalidRequestException
		var ue *UnavailableException
		var te *TimedOutException
		var err error
		ret, ire, ue, te, err = c.client.GetCount(key, cp, sp, r.consistencyLevel)
		return ire, ue, te, err
	})

	if err != nil {
		return 0, err
	}

	return int(ret), nil
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

	var ret map[string][]*ColumnOrSuperColumn
	err := r.pool.run(func(c *connection) (*InvalidRequestException, *UnavailableException, *TimedOutException, error) {
		var ire *InvalidRequestException
		var ue *UnavailableException
		var te *TimedOutException
		var err error
		ret, ire, ue, te, err = c.client.MultigetSlice(keys, cp, sp, r.consistencyLevel)
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

	var ret map[string]int32
	err := r.pool.run(func(c *connection) (*InvalidRequestException, *UnavailableException, *TimedOutException, error) {
		var ire *InvalidRequestException
		var ue *UnavailableException
		var te *TimedOutException
		var err error
		ret, ire, ue, te, err = c.client.MultigetCount(keys, cp, sp, r.consistencyLevel)
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
		return nil, nil
	}

	kr := r.buildKeyRange(rang)
	cp := r.buildColumnParent()
	sp := r.buildPredicate()

	var ret []*KeySlice
	err := r.pool.run(func(c *connection) (*InvalidRequestException, *UnavailableException, *TimedOutException, error) {
		var ire *InvalidRequestException
		var ue *UnavailableException
		var te *TimedOutException
		var err error
		ret, ire, ue, te, err = c.client.GetRangeSlices(cp, sp, kr, r.consistencyLevel)
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

	if len(r.expressions) == 0 {
		return nil, errors.New("At least one Where call must be made")
	}

	if rang == nil || rang.Count <= 0 {
		return make([]*Row, 0), nil
	}

	ic := r.buildIndexClause(rang)
	cp := r.buildColumnParent()
	sp := r.buildPredicate()

	var ret []*KeySlice
	err := r.pool.run(func(c *connection) (*InvalidRequestException, *UnavailableException, *TimedOutException, error) {
		var ire *InvalidRequestException
		var ue *UnavailableException
		var te *TimedOutException
		var err error
		ret, ire, ue, te, err = c.client.GetIndexedSlices(cp, ic, sp, r.consistencyLevel)
		return ire, ue, te, err
	})

	if err != nil {
		return nil, err
	}

	return rowsFromTListKeySlice(ret), nil
}

func (r *reader) RangeScan() (<-chan *Row, <-chan error) {
	if r.cf == "" {
		panic(errors.New("No column family specified"))
	}
	kr := NewKeyRange()
	if len(r.startToken) != 0 {
		kr.StartToken = r.startToken
	} else {
		kr.StartToken = DEF_START_TOKEN
	}
	if len(r.endToken) != 0 {
		kr.EndToken = r.endToken
	} else {
		kr.EndToken = DEF_END_TOKEN
	}
	kr.RowFilter = r.expressions
	cp := r.buildColumnParent()
	sp := r.buildPredicate()

	data := make(chan *Row)
	rerr := make(chan error)

	go func() {
		defer close(rerr)
		defer close(data)

		for {
			var ksv []*KeySlice
			err := r.pool.run(func(c *connection) (*InvalidRequestException, *UnavailableException, *TimedOutException, error) {
				var ire *InvalidRequestException
				var ue *UnavailableException
				var te *TimedOutException
				var err error
				ksv, ire, ue, te, err = c.client.GetRangeSlices(cp, sp, kr, r.consistencyLevel)
				return ire, ue, te, err
			})

			if err != nil {
				glog.Error("Error in GetRangeSlices ", err)
				rerr <- err
				return
			}
			glog.V(2).Infof("Key slice vector size ", len(ksv))
			if len(ksv) == 0 {
				//phew. done
				return
			}
			if bytes.Equal(ksv[0].Key, kr.StartKey) {
				//AP: I'm sending a diarrhea beam your way, dear designer of cassandra iteration
				ksv = ksv[1:]
			}
			if len(ksv) == 0 {
				//phew. done
				return
			}
			kr.StartToken = ""
			kr.StartKey = ksv[len(ksv)-1].Key
			glog.V(2).Infof("Next batch starts with %q", kr.StartKey)
			for _, ks := range ksv {
				glog.V(2).Infof("Raw row key %s columns %v", string(ks.Key), ks.Columns)
				row := rowFromTListColumns(ks.Key, ks.Columns)
				glog.V(2).Infof("Row %q", row)
				if row != nil {
					data <- row
				}
			}
		}
	}()
	return data, rerr
}

func (r *reader) WideRowScan(key, startColumn []byte, batchSize int32, callback func(*Column) bool) error {
	keyRange := NewKeyRange()
	keyRange.StartKey = key
	keyRange.EndKey = key
	keyRange.Count = batchSize //yes, it is weird but this count means columns count for GetPagedSlice

	var ret []*KeySlice
	for {
		err := r.pool.run(func(c *connection) (*InvalidRequestException, *UnavailableException, *TimedOutException, error) {
			var ire *InvalidRequestException
			var ue *UnavailableException
			var te *TimedOutException
			var err error
			ret, ire, ue, te, err = c.client.GetPagedSlice(r.cf, keyRange, startColumn, r.consistencyLevel)
			return ire, ue, te, err
		})

		if err != nil {
			return err
		}
		if len(ret) == 0 {
			return nil //finished
		}
		if len(ret) != 1 {
			return errors.New("Unexpected return vector length")
		}
		columns := ret[0].Columns
		if len(columns) != 0 && bytes.Equal(columns[0].Column.Name, startColumn) {
			//skip the column we saw already
			columns = columns[1:]
		}
		if len(columns) == 0 {
			return nil //finished
		}
		startColumn = columns[len(columns)-1].Column.Name
		for _, col := range columns {
			if !callback(col.Column) {
				return nil // aborted by callback
			}
		}
	}
}

func rowFromTListColumns(key []byte, tl []*ColumnOrSuperColumn) *Row {
	if tl == nil || len(tl) <= 0 {
		return nil
	}
	r := &Row{Key: key}
	if len(tl) == 0 {
		return r
	}
	r.Columns = make([]*Column, 0, len(tl))
	for _, col := range tl {
		if col.Column != nil {
			r.Columns = append(r.Columns, col.Column)
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

func rowsFromTMap(tm map[string][]*ColumnOrSuperColumn) []*Row {
	if tm == nil || len(tm) <= 0 {
		return nil
	}
	r := make([]*Row, 0, len(tm))
	for skey, columns := range tm {
		row := rowFromTListColumns([]byte(skey), columns)
		if row != nil {
			r = append(r, row)
		}
	}
	return r
}

func rowsColumnCountFromTMap(tm map[string]int32) []*RowColumnCount {
	if tm == nil || len(tm) <= 0 {
		return nil
	}
	r := make([]*RowColumnCount, 0, len(tm))
	for skey, count := range tm {
		if count > 0 {
			r = append(r, &RowColumnCount{Key: []byte(skey), Count: int(count)})
		}
	}
	return r
}

func rowsFromTListKeySlice(tl []*KeySlice) []*Row {
	if tl == nil || len(tl) <= 0 {
		return nil
	}
	r := make([]*Row, 0)
	for _, keySlice := range tl {
		key := keySlice.Key
		row := rowFromTListColumns(key, keySlice.Columns)
		if row != nil {
			r = append(r, row)
		}
	}
	return r
}

func (r *Row) ColumnNames() [][]byte {
	names := [][]byte{}
	for _, col := range r.Columns {
		names = append(names, col.Name)
	}
	return names
}
