package gossie

import (
	"testing"
)

func TestConnection(t *testing.T) {

	c, err := NewConnection("127.0.0.1:9999", "NotExists", 3000)
	if err == nil {
		t.Fatal("Invalid connection parameters did not return error")
	}

	c, err = NewConnection("127.0.0.1:9160", "Twitter", 3000)
	if err != nil {
		t.Fatal("Error connecting to Cassandra:", err)
	}

	c.Close()
}
