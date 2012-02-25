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

	c, err = newConnection("127.0.0.1:9160", "NotExists", 1000)
	if err == nil {
		t.Fatal("Invalid keyspace did not return error")
	}

	c, err = newConnection("127.0.0.1:9160", "TestGossie", 1000)
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
	var err os.Error
	var c *connection

	cpI, err := NewConnectionPool([]string{"127.0.0.1:9160"}, "TestGossie", PoolOptions{Size:1,Grace:1})
	if err != nil {
		t.Fatal("Error connecting to Cassandra:", err)
	}
	cp := cpI.(*connectionPool)

    check := func (expectedAvailable int, expectedError bool) {
		if len(cp.available) != expectedAvailable {
			t.Error("Available connection slots chan has wrong size")
		}
		if !expectedError && err != nil {
			t.Error("The error condition did not match the expected one")
		}
	}

	check(1, false)

	c, err = cp.acquire()
	check(0, false)

	cp.release(c)
	check(1, false)

	c, err = cp.acquire()
	check(0, false)

	cp.blacklist("127.0.0.1:9160")
	check(1, false)

	c, err = cp.acquire()
	check(1, true)

	time.Sleep(2e9)

	c, err = cp.acquire()
	check(0, false)
}

func TestRun(t *testing.T) {
	var err os.Error

	cpI, err := NewConnectionPool([]string{"127.0.0.1:9160"}, "TestGossie", PoolOptions{Size:1})
	if err != nil {
		t.Fatal("Error connecting to Cassandra:", err)
	}
	cp := cpI.(*connectionPool)
    var gotConnection bool

    check := func (_ire, _ue, _te, _err, expectedError bool) {
	    err := cp.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, os.Error) {
	    	gotConnection = c.client != nil
	        var ire *cassandra.InvalidRequestException
	        var ue *cassandra.UnavailableException
	        var te *cassandra.TimedOutException
	        var err os.Error
	        if _ire {
	        	ire = &cassandra.InvalidRequestException{}
	        }
	        if _ue {
	        	ue = &cassandra.UnavailableException{}
	        }
	        if _te {
	        	te = &cassandra.TimedOutException{}
	        }
	        if _err {
	        	err = os.NewError("uh")
	        }
	        return ire, ue, te, err
	    })

		if !gotConnection {
			t.Error("The transaction was passed a nil connection")
		}

		if !expectedError && err != nil {
			t.Error("The error condition did not match the expected one")
		}
    }

    check(false, false, false, false, false)
    check(true, false, false, false, true)
    check(false, true, false, false, true)
    check(false, false, true, false, true)
    check(false, false, false, true, true)
}


func BenchmarkGetOne(b *testing.B) {
	b.StopTimer()
	cp, _ := NewConnectionPool([]string{"127.0.0.1:9160", "127.0.0.1:9170", "127.0.0.1:9180"}, "TestGossie", PoolOptions{Size:50,Timeout:3000})
	b.StartTimer()
    for i := 0; i < b.N; i++ {
		cp.Query().Cf("AllTypes").Key([]byte("a")).GetOne()
    }
}
