package gossie

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/apesternikov/gossie/src/cassandra"
)

// Writer is the interface for all the write operations over Cassandra.
// The method calls support chaining so you can build concise queries
type Writer interface {

	// ConsistencyLevel sets the consistency level for this particular call.
	// It is optional, if left uncalled it will default to your connection
	// pool options value.
	ConsistencyLevel(cassandra.ConsistencyLevel) Writer

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
	pool             connectionRunner
	consistencyLevel cassandra.ConsistencyLevel
	writers          map[string]map[string][]*cassandra.Mutation
	usedCounters     bool
}

func newWriter(cp connectionRunner, cl cassandra.ConsistencyLevel) *writer {
	return &writer{
		pool:             cp,
		consistencyLevel: cl,
		writers:          make(map[string]map[string][]*cassandra.Mutation),
	}
}

func now() int64 {
	return nowfunc().UnixNano() / 1000
}

func (w *writer) addWriter(cf string, key []byte) *cassandra.Mutation {
	tm := cassandra.NewMutation()
	skey := string(key)
	if _, exists := w.writers[skey]; !exists {
		w.writers[skey] = make(map[string][]*cassandra.Mutation, 1)
	}
	w.writers[skey][cf] = append(w.writers[skey][cf], tm)
	return tm
}

func (w *writer) ConsistencyLevel(l cassandra.ConsistencyLevel) Writer {
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
			c.Ttl = thrift.Int32Ptr(int32(ttl))
		}
		if col.Timestamp != nil {
			c.Timestamp = col.Timestamp
		} else {
			c.Timestamp = &t
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
	d.Timestamp = thrift.Int64Ptr(now())
	tm.Deletion = d
	return w
}

func (w *writer) DeleteColumns(cf string, key []byte, columns [][]byte) Writer {
	tm := w.addWriter(cf, key)
	d := cassandra.NewDeletion()
	d.Timestamp = thrift.Int64Ptr(now())
	sp := cassandra.NewSlicePredicate()
	sp.ColumnNames = columns
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
	toRun := func(c *connection) error {
		return c.client.BatchMutate(w.writers, cassandra.ConsistencyLevel(w.consistencyLevel))
	}
	if w.usedCounters {
		return w.pool.runWithRetries(toRun, 1)
	}
	return w.pool.run(toRun)
}
