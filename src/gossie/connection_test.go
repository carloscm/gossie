package gossie

import (
	"testing"
)

func TestConnection(t *testing.T) {

	c, err := newConnection("127.0.0.1:9999", "NotExists", 3000)
	if err == nil {
		t.Fatal("Invalid connection parameters did not return error")
	}

	c, err = newConnection("127.0.0.1:9160", "NotExists", 3000)
	if err == nil {
		t.Fatal("Invalid keyspace did not return error")
	}

	c, err = newConnection("127.0.0.1:9160", "TestGossie", 3000)
	if err != nil {
		t.Fatal("Error connecting to Cassandra:", err)
	}

	if c.keyspace != "TestGossie" {
		t.Fatal("Invalid keyspace")
	}

	c.close()
}

func TestConnectionPool(t *testing.T) {

	cp, err := NewConnectionPool([]string{"127.0.0.1:9999"}, "NotExists", PoolOptions{Size:50,Timeout:3000})
	if err == nil {
		t.Fatal("Invalid connection parameters did not return error")
	}

	cp, err = NewConnectionPool([]string{"127.0.0.1:9160"}, "NotExists", PoolOptions{Size:50,Timeout:3000})
	if err == nil {
		t.Fatal("Invalid keyspace did not return error")
	}

	cp, err = NewConnectionPool([]string{"127.0.0.1:9160"}, "TestGossie", PoolOptions{Size:50,Timeout:3000})
	if err != nil {
		t.Fatal("Error connecting to Cassandra:", err)
	}

	if cp.Keyspace() != "TestGossie" {
		t.Fatal("Invalid keyspace")
	}

	t.Log("with content")
	row, err := cp.Query().Cf("AllTypes").Key([]byte("a")).GetOne()
	t.Log(row)
	t.Log(err)

	t.Log("without content")
	row, err = cp.Query().Cf("AllTypes").Key([]byte("b")).GetOne()
	t.Log(row)
	t.Log(err)

	t.Fatal("wut")

	cp.Close()
}
