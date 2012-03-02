package gossie

import (
    "os"
    "thrift"
    "cassandra"
    "time"
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

// RowColumnCount stores the number of columns matched in a MultiCount query
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
// new RangeGet query to page results. This will allow you to iterate over an entire CF even when
// using the random partitioner. Always specify a Count value since there is an interface-mandated
// default of 100.
type Range struct {
    Start []byte
    End   []byte
    Count int
}

// IndexedRange represents a range of rows to return for the IndexedGet method.
// The low level token range is not exposed. Use an empty slice to indicate if you want the first key
// in a range, then pass the last read row key as the new Start key in a new IndexedGet query to page
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

// Query is the interface for all read operations over Cassandra.
// The method calls support chaining so you can build concise queries
type Query interface {

    // ConsistencyLevel sets the consistency level for this particular call.
    // It is optional, if left uncalled it will default to your connection pool options value.
    ConsistencyLevel(int) Query

    // Cf sets the column family name for the query.
    // This method must be always called.
    Cf(string) Query

    // Slice optionally sets a slice to set a range of column names and potentially iterate over the
    // columns of the returned row(s)
    Slice(*Slice) Query

    // Columns optionally filters the returned columns to only the passed set of column names
    Columns([][]byte) Query

    // Each call to this method adds a new comparison to be checked against the returned rows of
    // IndexedGet
    // All the comparisons are checked for every row. In the current Cassandra implementation at
    // least one of the Where calls must use a secondary indexed column with an EQ operator.
    Where(column []byte, op Operator, value []byte) Query

    // Get looks up a row with the given key and returns it, or nil in case it is not found
    Get(key []byte) (*Row, os.Error)

    // MultiGet performs a parallel Get operation for all the passed keys, and returns a slice of
    // RowColumnCounts pointers to the gathered rows, which may be empty if none were found. It returns
    // nil only on error conditions
    MultiGet(keys [][]byte) ([]*Row, os.Error)

    // Count looks up a row with the given key and returns the number of columns it has
    Count(key []byte) (int, os.Error)

    // MultiGet performs a parallel Count operation for all the passed keys, and returns a slice of Row
    // pointers to the gathered rows, which may be empty if none were found. It returns nil only on
    // error conditions
    MultiCount(keys [][]byte) ([]*RowColumnCount, os.Error)

    // MultiGet performs a sequential Get operation for a range of rows. See the docs for Range for an
    // explanation on how to page results. It returns a slice of Row pointers to the gathered rows, which
    // may be empty if none were found. It returns nil only on error conditions
    RangeGet(*Range) ([]*Row, os.Error)

    // IndexedGet performs a sequential Get operation for a range of rows and returns only those that match
    // the Where clauses. See the docs for Range for an explanation on how to page results. It returns a
    // slice of Row pointers to the gathered rows, which may be empty if none were found. It returns nil only
    // on error conditions
    IndexedGet(*IndexedRange) ([]*Row, os.Error)
}

// Mutation is the interface for all the write operations over Cassandra.
// The method calls support chaining so you can build concise queries
type Mutation interface {

    // Insert adds a new row insertion to the mutation
    Insert(cf string, row *Row) Mutation

    // DeltaCounters add a new delta operation over counters
    DeltaCounters(cf string, row *Row) Mutation

    // Delete deletes a single row specified by key
    Delete(cf string, key []byte) Mutation

    // Delete deletes the passed columns from the row specified by key
    DeleteColumns(cf string, key []byte, columns [][]byte) Mutation

    //DeleteSlice(cf string, key []byte, slice *Slice) Mutation

    // Run this mutation
    Run() os.Error
}

type query struct {
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

func (q *query) ConsistencyLevel(l int) Query {
    q.consistencyLevel = l
    return q
}

func (q *query) Cf(cf string) Query {
    q.cf = cf
    return q
}

func (q *query) Slice(s *Slice) Query {
    q.slice = *s
    q.setSlice = true
    return q
}

func (q *query) Columns(c [][]byte) Query {
    copy(q.columns, c)
    q.setColumns = true
    return q
}

func (q *query) Where(column []byte, op Operator, value []byte) Query {
    if q.expressions == nil {
        q.expressions = thrift.NewTList(thrift.STRUCT, 1)
    }
    exp := cassandra.NewIndexExpression()
    exp.ColumnName = column
    exp.Op = cassandra.IndexOperator(op)
    exp.Value = value
    q.expressions.Push(exp)
    q.setWhere = true
    return q
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

func (q *query) buildPredicate() *cassandra.SlicePredicate {
    sp := cassandra.NewSlicePredicate()
    if q.setColumns {
        sp.ColumnNames = thrift.NewTList(thrift.BINARY, 1)
        for _, col := range q.columns {
            sp.ColumnNames.Push(col)
        }
    } else if q.setSlice {
        sp.SliceRange = sliceToCassandra(&q.slice)
    } else {
        sp.SliceRange = fullSlice()
    }
    return sp
}

func (q *query) buildColumnParent() *cassandra.ColumnParent {
    cp := cassandra.NewColumnParent()
    cp.ColumnFamily = q.cf
    return cp
}

func (q *query) buildKeyRange(r *Range) *cassandra.KeyRange {
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

func (q *query) buildIndexClause(r *IndexedRange) *cassandra.IndexClause {
    ic := cassandra.NewIndexClause()
    ic.Expressions = q.expressions
    ic.StartKey = r.Start
    ic.Count = int32(r.Count)
    // workaround some uninitialized slice == nil quirks that trickle down into the generated thrift4go code
    if ic.StartKey == nil {
        ic.StartKey = make([]byte, 0)
    }
    return ic
}

func (q *query) Get(key []byte) (*Row, os.Error) {
    if q.cf == "" {
        return nil, os.NewError("No column family specified")
    }

    cp := q.buildColumnParent()
    sp := q.buildPredicate()

    var ret thrift.TList
    err := q.pool.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, os.Error) {
        var ire *cassandra.InvalidRequestException
        var ue *cassandra.UnavailableException
        var te *cassandra.TimedOutException
        var err os.Error
        ret, ire, ue, te, err = c.client.GetSlice(key, cp, sp, cassandra.ConsistencyLevel(q.consistencyLevel))
        return ire, ue, te, err
    })

    if err != nil {
        return nil, err
    }

    return rowFromTListColumns(key, ret), nil
}

func (q *query) Count(key []byte) (int, os.Error) {
    if q.cf == "" {
        return 0, os.NewError("No column family specified")
    }

    cp := q.buildColumnParent()
    sp := q.buildPredicate()

    var ret int32
    err := q.pool.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, os.Error) {
        var ire *cassandra.InvalidRequestException
        var ue *cassandra.UnavailableException
        var te *cassandra.TimedOutException
        var err os.Error
        ret, ire, ue, te, err = c.client.GetCount(key, cp, sp, cassandra.ConsistencyLevel(q.consistencyLevel))
        return ire, ue, te, err
    })

    if err != nil {
        return 0, err
    }

    return int(ret), nil
}

func (q *query) buildMultiKeys(keys [][]byte) thrift.TList {
    tkeys := thrift.NewTList(thrift.BINARY, 1)
    for _, k := range keys {
        tkeys.Push(k)
    }
    return tkeys
}

func (q *query) MultiGet(keys [][]byte) ([]*Row, os.Error) {
    if q.cf == "" {
        return nil, os.NewError("No column family specified")
    }

    if len(keys) <= 0 {
        return make([]*Row, 0), nil
    }

    cp := q.buildColumnParent()
    sp := q.buildPredicate()
    tk := q.buildMultiKeys(keys)

    var ret thrift.TMap
    err := q.pool.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, os.Error) {
        var ire *cassandra.InvalidRequestException
        var ue *cassandra.UnavailableException
        var te *cassandra.TimedOutException
        var err os.Error
        ret, ire, ue, te, err = c.client.MultigetSlice(tk, cp, sp, cassandra.ConsistencyLevel(q.consistencyLevel))
        return ire, ue, te, err
    })

    if err != nil {
        return nil, err
    }

    return rowsFromTMap(ret), nil
}

func (q *query) MultiCount(keys [][]byte) ([]*RowColumnCount, os.Error) {
    if q.cf == "" {
        return nil, os.NewError("No column family specified")
    }

    if len(keys) <= 0 {
        return make([]*RowColumnCount, 0), nil
    }

    cp := q.buildColumnParent()
    sp := q.buildPredicate()
    tk := q.buildMultiKeys(keys)

    var ret thrift.TMap
    err := q.pool.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, os.Error) {
        var ire *cassandra.InvalidRequestException
        var ue *cassandra.UnavailableException
        var te *cassandra.TimedOutException
        var err os.Error
        ret, ire, ue, te, err = c.client.MultigetCount(tk, cp, sp, cassandra.ConsistencyLevel(q.consistencyLevel))
        return ire, ue, te, err
    })

    if err != nil {
        return nil, err
    }

    return rowsColumnCountFromTMap(ret), nil
}

func (q *query) RangeGet(rang *Range) ([]*Row, os.Error) {
    if q.cf == "" {
        return nil, os.NewError("No column family specified")
    }

    if rang == nil || rang.Count <= 0 {
        return make([]*Row, 0), nil
    }

    kr := q.buildKeyRange(rang)
    cp := q.buildColumnParent()
    sp := q.buildPredicate()

    var ret thrift.TList
    err := q.pool.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, os.Error) {
        var ire *cassandra.InvalidRequestException
        var ue *cassandra.UnavailableException
        var te *cassandra.TimedOutException
        var err os.Error
        ret, ire, ue, te, err = c.client.GetRangeSlices(cp, sp, kr, cassandra.ConsistencyLevel(q.consistencyLevel))
        return ire, ue, te, err
    })

    if err != nil {
        return nil, err
    }

    return rowsFromTListKeySlice(ret), nil
}

func (q *query) IndexedGet(rang *IndexedRange) ([]*Row, os.Error) {
    if q.cf == "" {
        return nil, os.NewError("No column family specified")
    }

    if !q.setWhere {
        return nil, os.NewError("At least one Where call must be made")
    }

    if rang == nil || rang.Count <= 0 {
        return make([]*Row, 0), nil
    }

    ic := q.buildIndexClause(rang)
    cp := q.buildColumnParent()
    sp := q.buildPredicate()

    var ret thrift.TList
    err := q.pool.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, os.Error) {
        var ire *cassandra.InvalidRequestException
        var ue *cassandra.UnavailableException
        var te *cassandra.TimedOutException
        var err os.Error
        ret, ire, ue, te, err = c.client.GetIndexedSlices(cp, ic, sp, cassandra.ConsistencyLevel(q.consistencyLevel))
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

type mutation struct {
    pool             *connectionPool
    consistencyLevel int
    mutations        thrift.TMap
}

func makeMutation(cp *connectionPool, cl int) *mutation {
    return &mutation{
        pool:             cp,
        consistencyLevel: cl,
        mutations:        thrift.NewTMap(thrift.BINARY, thrift.MAP, 1),
    }
}

func now() int64 {
    return time.Nanoseconds() / 1000
}

func (m *mutation) addMutation(cf string, key []byte) *cassandra.Mutation {
    tm := cassandra.NewMutation()
    var cfMuts thrift.TMap
    im, exists := m.mutations.Get(key)
    if !exists {
        cfMuts = thrift.NewTMap(thrift.STRING, thrift.LIST, 1)
        m.mutations.Set(key, cfMuts)
    } else {
        cfMuts = im.(thrift.TMap)
    }
    var mutList thrift.TList
    im, exists = cfMuts.Get(cf)
    if !exists {
        mutList = thrift.NewTList(thrift.STRUCT, 1)
        cfMuts.Set(cf, mutList)
    } else {
        mutList = im.(thrift.TList)
    }
    mutList.Push(tm)
    return tm
}

func (m *mutation) Insert(cf string, row *Row) Mutation {
    t := now()
    for _, col := range row.Columns {
        tm := m.addMutation(cf, row.Key)
        c := cassandra.NewColumn()
        c.Name = col.Name
        c.Value = col.Value
        c.Ttl = col.Ttl
        if col.Timestamp > 0 {
            c.Timestamp = col.Timestamp
        } else {
            c.Timestamp = t
        }
        cs := cassandra.NewColumnOrSuperColumn()
        cs.Column = c
        tm.ColumnOrSupercolumn = cs
    }
    return m
}

func (m *mutation) DeltaCounters(cf string, row *Row) Mutation {
    for _, col := range row.Columns {
        tm := m.addMutation(cf, row.Key)
        c := cassandra.NewCounterColumn()
        c.Name = col.Name
        Unmarshal(col.Value, LongType, &c.Value)
        cs := cassandra.NewColumnOrSuperColumn()
        cs.CounterColumn = c
        tm.ColumnOrSupercolumn = cs
    }
    return m
}

func (m *mutation) Delete(cf string, key []byte) Mutation {
    tm := m.addMutation(cf, key)
    d := cassandra.NewDeletion()
    d.Timestamp = now()
    tm.Deletion = d
    return m
}

func (m *mutation) DeleteColumns(cf string, key []byte, columns [][]byte) Mutation {
    tm := m.addMutation(cf, key)
    d := cassandra.NewDeletion()
    d.Timestamp = now()
    sp := cassandra.NewSlicePredicate()
    sp.ColumnNames = thrift.NewTList(thrift.BINARY, 1)
    for _, name := range columns {
        sp.ColumnNames.Push(name)
    }
    d.Predicate = sp
    tm.Deletion = d
    return m
}

/* InvalidRequestException({TStruct:InvalidRequestException Why:Deletion does not yet support SliceRange predicates.})
func (m *mutation) DeleteSlice(cf string, key []byte, slice *Slice) Mutation {
    tm := m.addMutation(cf, key)
    d := cassandra.NewDeletion()
    d.Timestamp = now()
    sp := cassandra.NewSlicePredicate()
    sp.SliceRange = sliceToCassandra(slice)
    d.Predicate = sp
    tm.Deletion = d
    return m
}
*/

func (m *mutation) Run() os.Error {
    return m.pool.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, os.Error) {
        var ire *cassandra.InvalidRequestException
        var ue *cassandra.UnavailableException
        var te *cassandra.TimedOutException
        var err os.Error
        ire, ue, te, err = c.client.BatchMutate(m.mutations, cassandra.ConsistencyLevel(m.consistencyLevel))
        return ire, ue, te, err
    })
}
