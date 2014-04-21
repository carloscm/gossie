package gossie

import (
	"testing"

	"code.google.com/p/gomock/gomock"
	. "github.com/wadey/gossie/src/cassandra"
	"github.com/wadey/gossie/src/gossie/mock_cassandra"
)

type stubTransactionRunner struct {
	conn *connection
}

func (s *stubTransactionRunner) run(t transaction) error {
	t(s.conn)
	return nil
}
func (s *stubTransactionRunner) runWithRetries(t transaction, retries int) error {
	t(s.conn)
	return nil
}

func NewBytes(in []byte) *[]byte { return &in }

func TestWriterInsert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cli := mock_cassandra.NewMockCassandra(ctrl)
	//expectingBatch := thrift.NewTMap(k, v, s)
	cli.EXPECT().BatchMutate(gomock.Any(), ConsistencyLevel_ONE)
	conn := &connection{
		socket:    nil,
		transport: nil,
		client:    cli,
		node:      &node{node: "node"},
	}
	cp := &stubTransactionRunner{conn: conn}
	w := newWriter(cp, CONSISTENCY_ONE)
	row := &Row{
		Key: []byte("rowkey"),
		Columns: []*Column{
			&Column{Name: []byte("name1"), Value: []byte("value1")},
			&Column{Name: []byte("name2"), Value: []byte("value2")},
		},
	}
	w.Insert("cf", row)
	e := w.Run()
	if e != nil {
		t.Error("Error", e)
	}
}
