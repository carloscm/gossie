package gossie

import (
	"testing"
	"cassandra"
	"os"
	"time"
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

func TestNewConnectionPool(t *testing.T) {

	cp, err := NewConnectionPool([]string{"127.0.0.1:9999"}, "NotExists", PoolOptions{Size:50,Timeout:3000})
	if err == nil {
		t.Fatal("Invalid connection parameters did not return error")
	}

	cp, err = NewConnectionPool([]string{"127.0.0.1:9160"}, "NotExists", PoolOptions{Size:50,Timeout:3000})
	if err == nil {
		t.Fatal("Invalid keyspace did not return error")
	}

	cp, err = NewConnectionPool([]string{"127.0.0.1:9160", "127.0.0.1:9170", "127.0.0.1:9180"}, "TestGossie", PoolOptions{Size:50,Timeout:3000})
	if err != nil {
		t.Fatal("Error connecting to Cassandra:", err)
	}

	if cp.Keyspace() != "TestGossie" {
		t.Fatal("Invalid keyspace")
	}
/*
	t.Log("with content")
	row, err := cp.Query().Cf("AllTypes").Key([]byte("a")).GetOne()
	t.Log(string(row.Key))
	t.Log(string(row.Columns[0].Name))
	t.Log(string(row.Columns[0].Value))
	t.Log(row.Columns[0].Ttl)
	t.Log(row.Columns[0].Timestamp)
	t.Log(err)

	t.Log("without content")
	row, err = cp.Query().Cf("AllTypes").Key([]byte("b")).GetOne()
	t.Log(row)
	t.Log(err)

	//t.Fatal("wut")
*/
	cp.Close()
}

func TestAcquireRelease(t *testing.T) {
	cpI, _ := NewConnectionPool([]string{"127.0.0.1:9160"}, "TestGossie", PoolOptions{Size:1,Grace:1})
	cp := cpI.(*connectionPool)

	if len(cp.available) != 1 {
		t.Error("Available connection slots chan has wrong size")
	}

	c, err := cp.acquire()
	if len(cp.available) != 0 {
		t.Error("Available connection slots chan has wrong size after acquire")
	}
	if c == nil || err != nil {
		t.Error("Normal acquire returned error")
	}

	cp.release(c)
	if len(cp.available) != 1 {
		t.Error("Available connection slots chan has wrong size after release")
	}

	c, err = cp.acquire()
	if len(cp.available) != 0 {
		t.Error("Available connection slots chan has wrong size after acquire")
	}
	if c == nil || err != nil {
		t.Error("Normal acquire returned error")
	}

	cp.blacklist(c)
	if len(cp.available) != 1 {
		t.Error("Available connection slots chan has wrong size after blacklist")
	}

	c, err = cp.acquire()
	if len(cp.available) != 1 {
		t.Error("Available connection slots chan has wrong size after acquire with blacklist")
	}
	if c != nil || err == nil {
		t.Error("Aquire with all nodes down did not return error")
	}

	time.Sleep(2000000000)

	c, err = cp.acquire()
	if len(cp.available) != 0 {
		t.Error("Available connection slots chan has wrong size after acquire")
	}
	if c == nil || err != nil {
		t.Error("Normal acquire after blacklist grace period returned error")
	}

}

func TestRun(t *testing.T) {
	cpI, _ := NewConnectionPool([]string{"127.0.0.1:9160"}, "TestGossie", PoolOptions{Size:1})

	cp := cpI.(*connectionPool)

    var gotConnection bool

    err := cp.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, os.Error) {
    	gotConnection = c.client != nil
        var ire *cassandra.InvalidRequestException
        var ue *cassandra.UnavailableException
        var te *cassandra.TimedOutException
        var err os.Error
        return ire, ue, te, err
    })
	if !gotConnection || err != nil {
		t.Error("Unexpected error in normal run() call")
	}

    err = cp.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, os.Error) {
    	gotConnection = c.client != nil
        var ire *cassandra.InvalidRequestException = &cassandra.InvalidRequestException{}
        var ue *cassandra.UnavailableException
        var te *cassandra.TimedOutException
        var err os.Error
        return ire, ue, te, err
    })
	if !gotConnection || err == nil {
		t.Error("Expected error in run() call did not trigger (ire)")
	}

    err = cp.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, os.Error) {
    	gotConnection = c.client != nil
        var ire *cassandra.InvalidRequestException
        var ue *cassandra.UnavailableException = &cassandra.UnavailableException{}
        var te *cassandra.TimedOutException
        var err os.Error
        return ire, ue, te, err
    })
	if !gotConnection || err == nil {
		t.Error("Expected error in run() call did not trigger (ue)")
	}

    err = cp.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, os.Error) {
    	gotConnection = c.client != nil
        var ire *cassandra.InvalidRequestException
        var ue *cassandra.UnavailableException
        var te *cassandra.TimedOutException = &cassandra.TimedOutException{}
        var err os.Error
        return ire, ue, te, err
    })
	if !gotConnection || err == nil {
		t.Error("Expected error in run() call did not trigger (te)")
	}

    err = cp.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, os.Error) {
    	gotConnection = c.client != nil
        var ire *cassandra.InvalidRequestException
        var ue *cassandra.UnavailableException
        var te *cassandra.TimedOutException
        var err os.Error = os.NewError("uhh")
        return ire, ue, te, err
    })
	if !gotConnection || err == nil {
		t.Error("Expected error in run() call did not trigger (err)")
	}

}


func BenchmarkGetOne(b *testing.B) {
	b.StopTimer()
	cp, _ := NewConnectionPool([]string{"127.0.0.1:9160", "127.0.0.1:9170", "127.0.0.1:9180"}, "TestGossie", PoolOptions{Size:50,Timeout:3000})
	b.StartTimer()
    for i := 0; i < b.N; i++ {
		cp.Query().Cf("AllTypes").Key([]byte("a")).GetOne()
    }
}
