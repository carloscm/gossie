package mockgossie

import (
	"bytes"
	. "github.com/wadey/gossie/src/gossie"
	"sort"
	"time"
)

type MockWriter struct {
	pool *MockConnectionPool
}

func newWriter(cp *MockConnectionPool) *MockWriter {
	return &MockWriter{
		pool: cp,
	}
}

func now() int64 {
	return time.Now().UnixNano() / 1000
}

func (w *MockWriter) ConsistencyLevel(l int) Writer {
	return w
}

func (w *MockWriter) Insert(cf string, row *Row) Writer {
	rows := w.pool.Rows(cf)

	i := sort.Search(len(rows), func(i int) bool { return bytes.Compare(rows[i].Key, row.Key) >= 0 })
	if i < len(rows) && bytes.Equal(rows[i].Key, row.Key) {
		// Row already exists, merge the columns
		e := rows[i]
		cols := e.Columns
		for _, c := range row.Columns {
			j := sort.Search(len(cols), func(j int) bool { return bytes.Compare(cols[j].Name, c.Name) >= 0 })
			if j < len(cols) && bytes.Equal(cols[j].Name, c.Name) {
				// Column already exists, pick the one with the greater timestamp
				ec := cols[j]
				if c.Timestamp >= ec.Timestamp {
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

func (w *MockWriter) InsertTtl(cf string, row *Row, ttl int) Writer {
	panic("Not Implemented")
}

func (w *MockWriter) DeltaCounters(cf string, row *Row) Writer {
	panic("Not Implemented")
}

func (w *MockWriter) Delete(cf string, key []byte) Writer {
	panic("Not Implemented")
}

func (w *MockWriter) DeleteColumns(cf string, key []byte, columns [][]byte) Writer {
	panic("Not Implemented")
}

func (w *MockWriter) Run() error {
	return nil
}
