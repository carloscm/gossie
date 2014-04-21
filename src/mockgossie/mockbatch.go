package mockgossie

import (
	. "github.com/wadey/gossie/src/cassandra"
	. "github.com/wadey/gossie/src/gossie"
)

type MockBatch struct {
	pool             *MockConnectionPool
	writer           Writer
	consistencyLevel int
	ttl              int
	mappingError     error
}

var _ Batch = &MockBatch{}

func newBatch(cp *MockConnectionPool) *MockBatch {
	return &MockBatch{
		pool:   cp,
		writer: cp.Writer(),
	}
}

func (b *MockBatch) ConsistencyLevel(c ConsistencyLevel) Batch {
	b.writer.ConsistencyLevel(c)
	return b
}

func (b *MockBatch) Ttl(ttl int) Batch {
	b.ttl = ttl
	return b
}

func (b *MockBatch) Insert(mapping Mapping, data interface{}) Batch {
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

func (b *MockBatch) Delete(mapping Mapping, data interface{}) Batch {
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

func (b *MockBatch) DeleteAll(mapping Mapping, data interface{}) Batch {
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

func (b *MockBatch) GetWriter() Writer {
	return b.writer
}

func (b *MockBatch) Run() error {
	if b.mappingError != nil {
		return b.mappingError
	}
	return b.writer.Run()
}
