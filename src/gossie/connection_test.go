package gossie

import (
	"errors"
	"github.com/apesternikov/gossie/src/cassandra"
	"github.com/stretchrcom/testify/assert"
	"testing"
	"time"
)

func TestConnection(t *testing.T) {

	/* kind of pointless
	   c, err := newConnection(invalidEndpoint, "NotExists", standardTimeout)
	   if err == nil {
	       t.Fatal("Invalid connection parameters did not return error")
	   }
	*/
	n := &node{node: localEndpoint}
	c, err := newConnection(n, "NotExists", shortTimeout, map[string]string{})
	if err == nil {
		t.Fatal("Invalid keyspace did not return error")
	}

	c, err = newConnection(n, keyspace, shortTimeout, map[string]string{})
	if err != nil {
		t.Fatal("Error connecting to Cassandra:", err)
	}

	assert.Equal(t, c.node, n)

	c.close()
}

func TestNewConnectionPool(t *testing.T) {

	/* kind of pointless
	   cp, err := NewConnectionPool([]string{invalidEndpoint}, "NotExists", poolOptions)
	   if err == nil {
	       t.Fatal("Invalid connection parameters did not return error")
	   }
	*/

	cp, err := NewConnectionPool(localEndpointPool, "NotExists", poolOptions)
	if err == nil {
		t.Fatal("Invalid keyspace did not return error")
	}

	/* kind of pointless
	   cp, err = NewConnectionPool(localEndpointsPool, keyspace, poolOptions)
	   if err != nil {
	       t.Fatal("Error connecting to Cassandra:", err)
	   }
	*/

	cp, err = NewConnectionPool(localEndpointPool, keyspace, poolOptions)
	if err != nil {
		t.Fatal("Error connecting to Cassandra:", err)
	}

	if cp.Keyspace() != keyspace {
		t.Fatal("Invalid keyspace")
	}

	cp.Close()
}

// possible test for users of SimpleAuthenticator
/*
func TestNewConnectionPoolWithAuth(t *testing.T) {
	poolOptionsAuth := poolOptions
	poolOptionsAuth.Authentication = map[string]string{
		"keyspace": "invalid",
		"username": "invalid",
		"password": "invalid",
	}
	cp, err := NewConnectionPool(localEndpointPool, "TestGossie", poolOptionsAuth)
	if err == nil {
		t.Fatal("Invalid keyspace did not return error")
	}

	poolOptionsAuth.Authentication["keyspace"] = "TestGossie"
	cp, err = NewConnectionPool(localEndpointPool, "TestGossie", poolOptionsAuth)
	if err == nil {
		t.Fatal("Invalid username did not return error")
	}

	poolOptionsAuth.Authentication["username"] = "test"
	cp, err = NewConnectionPool(localEndpointPool, "TestGossie", poolOptionsAuth)
	if err == nil {
		t.Fatal("Invalid password did not return error")
	}

	poolOptionsAuth.Authentication["password"] = "testpw"
	cp, err = NewConnectionPool(localEndpointPool, "TestGossie", poolOptionsAuth)
	if err != nil {
		t.Fatal("Correct credetinals did not allow login")
	}

	if cp.Keyspace() != keyspace {
		t.Fatal("Invalid keyspace")
	}

	cp.Close()
}
*/

func TestAcquireRelease(t *testing.T) {
	var err error
	var c *connection

	cpI, err := NewConnectionPool(localEndpointPool, keyspace, PoolOptions{Size: 1, Grace: 1})
	if err != nil {
		t.Fatal("Error connecting to Cassandra:", err)
	}
	cp := cpI.(*connectionPool)
	n := cp.nodes[0]

	assert.Equal(t, len(n.available.l), 1)
	assert.NoError(t, err)

	c, err = cp.acquire()
	assert.Equal(t, len(n.available.l), 0)
	assert.NoError(t, err)

	cp.release(c)
	assert.Equal(t, len(n.available.l), 1)
	assert.NoError(t, err)

	c, err = cp.acquire()
	assert.Equal(t, len(n.available.l), 0)
	assert.NoError(t, err)

	n.blacklist()
	assert.Equal(t, len(n.available.l), 0)
	assert.NoError(t, err)

	c, err = cp.acquire()
	assert.Equal(t, len(n.available.l), 0)
	assert.Error(t, err)
	assert.Nil(t, c)

	//TODO: replace with fake clock for testability
	time.Sleep(2e9)

	c, err = cp.acquire()
	assert.Equal(t, len(n.available.l), 0)
	assert.NoError(t, err)
}

func TestRun(t *testing.T) {
	var err error

	cpI, err := NewConnectionPool(localEndpointPool, keyspace, PoolOptions{Size: 1})
	if err != nil {
		t.Fatal("Error connecting to Cassandra:", err)
	}
	cp := cpI.(*connectionPool)
	var gotConnection bool

	check := func(_ire, _ue, _te, _err, expectedError bool) {
		err := cp.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, error) {
			gotConnection = c.client != nil
			var ire *cassandra.InvalidRequestException
			var ue *cassandra.UnavailableException
			var te *cassandra.TimedOutException
			var err error
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
				err = errors.New("uh")
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
