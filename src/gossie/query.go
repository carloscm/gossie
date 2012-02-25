package gossie

import (
    "os"
    "thrift"
    "cassandra"
    "time"
//    "fmt"
)

type Column struct {
    Name []byte
    Value []byte
    Ttl int32
    Timestamp int64
}

type Row struct {
    Key []byte
    Columns []*Column
}

type Slice struct {
    Start []byte
    Finish []byte
    Count int
    Reversed bool
}

type Range struct {
    Start []byte
    Finish []byte
    Count int
}

type Query interface {
    ConsistencyLevel(int) Query
    Cf(string) Query
    Key([]byte) Query
    Keys([][]byte) Query
    Slice(*Slice) Query
    Columns([][]byte) Query
    Range(*Range) Query
    // Index
    GetOne() (*Row, os.Error)
    //GetMany() ([]Row, os.Error)
    //CountOne() (int, os.Error)
    //CountMany() ([]int, os.Error)
}

type query struct {
    pool *connectionPool
    consistencyLevel int
    cf string
    key []byte
    setKey bool
    keys [][]byte
    setKeys bool
    slice Slice
    setSlice bool
    columns [][]byte
    setColumns bool
    qrange Range
    setRange bool
}

func (q *query) ConsistencyLevel(l int) Query {
    q.consistencyLevel = l
    return q
}

func (q *query) Cf(cf string) Query {
    q.cf = cf
    return q
}

func (q *query) Key(k []byte) Query {
    q.key = k
    q.setKey = true
    return q
}

func (q *query) Keys(k [][]byte) Query {
    q.keys = k
    q.setKeys = true
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

func (q *query) Range(r *Range) Query {
    q.qrange = *r
    q.setRange = true
    return q
}

func (q *query) GetOne() (*Row, os.Error) {
    if q.cf == "" {
        return nil, os.NewError("No column family specified")
    }
    if !q.setKey {
        return nil, os.NewError("No key or keys specified")
    }

    sp := cassandra.NewSlicePredicate()
    if q.setColumns {
        for _, col := range q.columns {
            sp.ColumnNames.Push(col)
        }
    } else if q.setSlice {
        sr := cassandra.NewSliceRange()
        sr.Start = q.slice.Start
        sr.Finish = q.slice.Finish
        sr.Count = int32(q.slice.Count)
        sr.Reversed = q.slice.Reversed
        sp.SliceRange = sr
    } else {
        sp.SliceRange = cassandra.NewSliceRange()
    }

    // workaround some uninitialized slice == nil quirks that trickle down into the generated thrift4go code
    if sp.SliceRange != nil {
        if sp.SliceRange.Start == nil {
            sp.SliceRange.Start = make([]byte, 0)
        }
        if sp.SliceRange.Finish == nil {
            sp.SliceRange.Finish = make([]byte, 0)
        }
    }

    cp := cassandra.NewColumnParent()
    cp.ColumnFamily = q.cf

    var ret thrift.TList
    err := q.pool.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, os.Error) {
        var ire *cassandra.InvalidRequestException
        var ue *cassandra.UnavailableException
        var te *cassandra.TimedOutException
        var err os.Error
        ret, ire, ue, te, err = c.client.GetSlice(q.key, cp, sp, cassandra.ConsistencyLevel(q.consistencyLevel))
        return ire, ue, te, err
    })

    if err != nil {
        return nil, err
    }

    if ret != nil {
        return rowFromTList(q.key, ret), nil
    }

    return nil, nil
}

func rowFromTList(key []byte, tl thrift.TList) *Row {
    r := &Row{Key:key}
    for colI := range tl.Iter() {
        var col *cassandra.ColumnOrSuperColumn = colI.(*cassandra.ColumnOrSuperColumn)
        if col.Column != nil {
            c := &Column{
                Name:col.Column.Name,
                Value:col.Column.Value,
                Timestamp:col.Column.Timestamp,
                Ttl:col.Column.Ttl,
            }
            r.Columns = append(r.Columns, c)
        } else if col.CounterColumn != nil {
            v, _ := Marshal(col.CounterColumn.Value, LongType)
            c := &Column{
                Name:col.CounterColumn.Name,
                Value:v,
            }
            r.Columns = append(r.Columns, c)
        }
    }
    return r
}

type Mutation interface {
    Insert(cf string, row *Row) Mutation
    //DeltaCounters(cf string, row Row) Mutation
    //Delete(cf string, key []byte) Mutation
    //DeleteColumns(cf string, row Row) Mutation
    //DeleteSlice(cf string, key []byte, slice *Slice) Mutation
    Run() os.Error
}

type mutation struct {
    pool *connectionPool
    consistencyLevel int
    mutations thrift.TMap
}

func makeMutation(cp *connectionPool, cl int) *mutation {
    return &mutation {
        pool: cp,
        consistencyLevel: cl,
        mutations: thrift.NewTMap(thrift.BINARY, thrift.MAP, 1),
    }
}

func now() int64 {
    return time.Nanoseconds()/1000
}

func (m *mutation) Insert(cf string, row *Row) Mutation {
    t := now()
    key := row.Key
    for _, col := range row.Columns {
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
        tm := cassandra.NewMutation()
        tm.ColumnOrSupercolumn = cs

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
    }

    return m
}

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
