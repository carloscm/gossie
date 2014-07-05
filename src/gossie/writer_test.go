package gossie

import (
	"code.google.com/p/gomock/gomock"
	. "github.com/apesternikov/gossie/src/cassandra"
	"github.com/apesternikov/gossie/src/gossie/mock_cassandra"
	"testing"
)

type stubTransactionRunner struct {
	conn *connection
}

func (s *stubTransactionRunner) Run(t func(client Cassandra) error) error {
	t(s.conn.client)
	return nil
}
func (s *stubTransactionRunner) RunWithRetries(t func(client Cassandra) error, retries int) error {
	t(s.conn.client)
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
