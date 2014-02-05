package gossie

import (
	"github.com/hailocab/gossie/src/cassandra"
	"github.com/hailocab/thrift4go/lib/go/src/thrift"
	"time"
)

// Writer is the interface for all the write operations over Cassandra.
// The method calls support chaining so you can build concise queries
type Writer interface {

	// ConsistencyLevel sets the consistency level for this particular call.
	// It is optional, if left uncalled it will default to your connection
	// pool options value.
	ConsistencyLevel(int) Writer

	// Insert adds a new row insertion to the mutation
	Insert(cf string, row *Row) Writer

	// InsertTtl adds a new row insertion to the mutation, overriding the
	// columns Ttl with the passed value
	InsertTtl(cf string, row *Row, ttl int) Writer

	// DeltaCounters add a new delta operation over counters. When using this
	// function the number of retries for the connection is temporally set to
	// 1.
	DeltaCounters(cf string, row *Row) Writer

	// Delete deletes a single row specified by key
	Delete(cf string, key []byte) Writer

	// DeleteColumns deletes the passed columns from the row specified by key.
	DeleteColumns(cf string, key []byte, columns [][]byte) Writer

	// Run this mutation
	Run() error
}

type writer struct {
	pool             *connectionPool
	consistencyLevel int
	writers          thrift.TMap
	usedCounters     bool
}

func newWriter(cp *connectionPool, cl int) *writer {
	return &writer{
		pool:             cp,
		consistencyLevel: cl,
		writers:          thrift.NewTMap(thrift.BINARY, thrift.MAP, 1),
	}
}

func now() int64 {
	return time.Now().UnixNano() / 1000
}

func (w *writer) addWriter(cf string, key []byte) *cassandra.Mutation {
	tm := cassandra.NewMutation()
	var cfMuts thrift.TMap
	im, exists := w.writers.Get(key)
	if !exists {
		cfMuts = thrift.NewTMap(thrift.STRING, thrift.LIST, 1)
		w.writers.Set(key, cfMuts)
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

func (w *writer) ConsistencyLevel(l int) Writer {
	w.consistencyLevel = l
	return w
}

func (w *writer) Insert(cf string, row *Row) Writer {
	return w.InsertTtl(cf, row, -1)
}

func (w *writer) InsertTtl(cf string, row *Row, ttl int) Writer {
	t := now()
	for _, col := range row.Columns {
		tm := w.addWriter(cf, row.Key)
		c := cassandra.NewColumn()
		c.Name = col.Name
		c.Value = col.Value
		if ttl > 0 {
			c.Ttl = int32(ttl)
		} else {
			c.Ttl = col.Ttl
		}
		if col.Timestamp > 0 {
			c.Timestamp = col.Timestamp
		} else {
			c.Timestamp = t
		}
		cs := cassandra.NewColumnOrSuperColumn()
		cs.Column = c
		tm.ColumnOrSupercolumn = cs
	}
	return w
}

func (w *writer) DeltaCounters(cf string, row *Row) Writer {
	for _, col := range row.Columns {
		tm := w.addWriter(cf, row.Key)
		c := cassandra.NewCounterColumn()
		c.Name = col.Name
		Unmarshal(col.Value, LongType, &c.Value)
		cs := cassandra.NewColumnOrSuperColumn()
		cs.CounterColumn = c
		tm.ColumnOrSupercolumn = cs
	}
	w.usedCounters = true
	return w
}

func (w *writer) Delete(cf string, key []byte) Writer {
	tm := w.addWriter(cf, key)
	d := cassandra.NewDeletion()
	d.Timestamp = now()
	tm.Deletion = d
	return w
}

func (w *writer) DeleteColumns(cf string, key []byte, columns [][]byte) Writer {
	tm := w.addWriter(cf, key)
	d := cassandra.NewDeletion()
	d.Timestamp = now()
	sp := cassandra.NewSlicePredicate()
	sp.ColumnNames = thrift.NewTList(thrift.BINARY, 1)
	for _, name := range columns {
		sp.ColumnNames.Push(name)
	}
	d.Predicate = sp
	tm.Deletion = d
	return w
}

/* InvalidRequestException({TStruct:InvalidRequestException Why:Deletion does not yet support SliceRange predicates.})
func (w *writer) DeleteSlice(cf string, key []byte, slice *Slice) Writer {
    tm := w.addWriter(cf, key)
    d := cassandra.NewDeletion()
    d.Timestamp = now()
    sp := cassandra.NewSlicePredicate()
    sp.SliceRange = sliceToCassandra(slice)
    d.Predicate = sp
    tm.Deletion = d
    return w
}
*/

func (w *writer) Run() error {
	toRun := func(c *connection) *transactionError {
		ire, ue, te, err := c.client.BatchMutate(
			w.writers, cassandra.ConsistencyLevel(w.consistencyLevel))
		return &transactionError{ire, ue, te, err}
	}
	if w.usedCounters {
		return w.pool.runWithRetries(toRun, 1)
	}
	return w.pool.run(toRun)
}
