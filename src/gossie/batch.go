package gossie

import (
	"github.com/apesternikov/gossie/src/cassandra"
)

// Batch is a high level interface for Cassandra writes. Simultaneous
// insertions for different column families and keys are possible.
type Batch interface {

	// ConsistencyLevel sets the consistency level for this particular call.
	// It is optional, if left uncalled it will default to your connection
	// pool options value.
	ConsistencyLevel(cassandra.ConsistencyLevel) Batch

	// Ttl sets a time to live for the columns inserted by Insert(). It is 0
	// by default which means no TTL.
	Ttl(int) Batch

	// Insert adds new data to be inserted
	Insert(mapping Mapping, data interface{}) Batch

	// Delete marks only the specific columns of the passed struct to be
	// deleted (respecting the composites).
	Delete(mapping Mapping, data interface{}) Batch

	// DeleteAll marks the entire row of the primary key to be deleted. This
	// will also delete any other struct present in the row if this column
	// family is using composites.
	DeleteAll(mapping Mapping, data interface{}) Batch

	// Run this batch
	Run() error
}

type batch struct {
	pool             *connectionPool
	writer           Writer
	consistencyLevel cassandra.ConsistencyLevel
	ttl              int
	mappingError     error
}

func newBatch(cp *connectionPool) *batch {
	return &batch{
		pool:   cp,
		writer: cp.Writer(),
	}
}

func (b *batch) ConsistencyLevel(c cassandra.ConsistencyLevel) Batch {
	b.writer.ConsistencyLevel(c)
	return b
}

func (b *batch) Ttl(ttl int) Batch {
	b.ttl = ttl
	return b
}

func (b *batch) Insert(mapping Mapping, data interface{}) Batch {
	if b.mappingError == nil {
		row, err := mapping.Map(data)
		if err == nil {
			if b.ttl > 0 {
				b.writer.InsertTtl(mapping.Cf(), row, b.ttl)
			} else {
				b.writer.Insert(mapping.Cf(), row)
			}
		} else {
			b.mappingError = err
		}
	}
	return b
}

func (b *batch) Delete(mapping Mapping, data interface{}) Batch {
	if b.mappingError == nil {
		row, err := mapping.Map(data)
		if err == nil {
			b.writer.DeleteColumns(mapping.Cf(), row.Key, row.ColumnNames())
		} else {
			b.mappingError = err
		}
	}
	return b
}

func (b *batch) DeleteAll(mapping Mapping, data interface{}) Batch {
	if b.mappingError == nil {
		row, err := mapping.Map(data)
		if err == nil {
			b.writer.Delete(mapping.Cf(), row.Key)
		} else {
			b.mappingError = err
		}
	}
	return b
}

func (b *batch) Run() error {
	if b.mappingError != nil {
		return b.mappingError
	}
	return b.writer.Run()
}
