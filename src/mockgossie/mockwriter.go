package mockgossie

import (
	"bytes"
	"sort"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
	. "github.com/wadey/gossie/src/cassandra"
	. "github.com/wadey/gossie/src/gossie"
)

type MockWriter struct {
	pool *MockConnectionPool
}

var _ Writer = &MockWriter{}

func newWriter(cp *MockConnectionPool) *MockWriter {
	return &MockWriter{
		pool: cp,
	}
}

func now() int64 {
	return time.Now().UnixNano() / 1000
}

func (w *MockWriter) ConsistencyLevel(c ConsistencyLevel) Writer {
	return w
}

func (w *MockWriter) Insert(cf string, row *Row) Writer {
	return w.InsertTtl(cf, row, -1)
}

func (w *MockWriter) InsertTtl(cf string, row *Row, ttl int) Writer {
	rows := w.pool.Rows(cf)

	t := thrift.Int64Ptr(now())
	for _, c := range row.Columns {
		if c.Timestamp == nil {
			c.Timestamp = t
		}
		if ttl > 0 {
			c.Ttl = thrift.Int32Ptr(int32(ttl))
		}
		if c.Ttl != nil {
			// reset to the actual time to expire
			c.Ttl = thrift.Int32Ptr(int32(now()/1e6) + *c.Ttl)
		}
	}

	i := sort.Search(len(rows), func(i int) bool { return bytes.Compare(rows[i].Key, row.Key) >= 0 })
	if i < len(rows) && bytes.Equal(rows[i].Key, row.Key) {
		// Row already exists, merge the columns
		e := rows[i]
		checkExpired(e)
		cols := e.Columns
		for _, c := range row.Columns {
			j := sort.Search(len(cols), func(j int) bool { return bytes.Compare(cols[j].Name, c.Name) >= 0 })
			if j < len(cols) && bytes.Equal(cols[j].Name, c.Name) {
				// Column already exists, pick the one with the greater timestamp
				ec := cols[j]
				et := *t
				if ec != nil {
					et = *ec.Timestamp
				}
				if *c.Timestamp >= et {
					ec.Value = c.Value
					ec.Ttl = c.Ttl
					ec.Timestamp = c.Timestamp
				}
			} else {
				// New column, insert sorted
				cols = append(cols, c)
				copy(cols[j+1:], cols[j:])
				cols[j] = c
			}
		}
		e.Columns = cols
	} else {
		// New row, insert sorted
		sort.Sort(Columns(row.Columns))
		rows = append(rows, row)
		copy(rows[i+1:], rows[i:])
		rows[i] = row

		w.pool.Data[cf] = rows
	}

	return w
}

func checkExpired(r *Row) {
	for i := 0; i < len(r.Columns); {
		c := r.Columns[i]
		if isExpired(c) {
			copy(r.Columns[i:], r.Columns[i+1:])
			r.Columns[len(r.Columns)-1] = nil
			r.Columns = r.Columns[:len(r.Columns)-1]
		} else {
			i++
		}
	}
}

func isExpired(c *Column) bool {
	if c.Ttl != nil {
		return int32(now()/1e6) > *c.Ttl
	}
	return false
}

func (w *MockWriter) DeltaCounters(cf string, row *Row) Writer {
	panic("Not Implemented")
}

func (w *MockWriter) Delete(cf string, key []byte) Writer {
	rows := w.pool.Rows(cf)

	t := now()

	i := sort.Search(len(rows), func(i int) bool { return bytes.Compare(rows[i].Key, key) >= 0 })
	if i < len(rows) && bytes.Equal(rows[i].Key, key) {
		// Row exists, delete the columns
		e := rows[i]
		cols := e.Columns
		for index, c := range cols {
			if t >= *c.Timestamp {
				// TODO store tombstone?
				copy(cols[index:], cols[index+1:])
				cols[len(cols)-1] = nil
				cols = cols[:len(cols)-1]
			}
		}
		e.Columns = cols
	}

	return w
}

func (w *MockWriter) DeleteColumns(cf string, key []byte, columns [][]byte) Writer {
	rows := w.pool.Rows(cf)

	t := now()

	i := sort.Search(len(rows), func(i int) bool { return bytes.Compare(rows[i].Key, key) >= 0 })
	if i < len(rows) && bytes.Equal(rows[i].Key, key) {
		// Row exists, delete the columns
		e := rows[i]
		cols := e.Columns
		for _, c := range columns {
			j := sort.Search(len(cols), func(j int) bool { return bytes.Compare(cols[j].Name, c) >= 0 })
			if j < len(cols) && bytes.Equal(cols[j].Name, c) {
				if t >= *cols[j].Timestamp {
					// TODO store tombstone?
					copy(cols[j:], cols[j+1:])
					cols[len(cols)-1] = nil
					cols = cols[:len(cols)-1]
				}
			}
		}
		e.Columns = cols
	}

	return w
}

func (w *MockWriter) Run() error {
	return nil
}
