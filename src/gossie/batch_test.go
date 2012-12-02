package gossie

import (
	"reflect"
	"testing"
)

func TestBatch(t *testing.T) {
	cp, err := NewConnectionPool(localEndpointPool, keyspace, PoolOptions{Size: 1, Timeout: shortTimeout})
	if err != nil {
		t.Fatal("Error connecting to Cassandra:", err)
	}

	m0, err := NewMapping(&ReasonableZero{})
	if err != nil {
		t.Fatal("Error building mapping:", err)
	}
	m2, err := NewMapping(&ReasonableTwo{})
	if err != nil {
		t.Fatal("Error building mapping:", err)
	}

	r := &ReasonableZero{"batchtest", 1.00002, -38.11, "hey this thing appears to work, nice!"}
	rs := &ReasonableZero{}

	r100 := &ReasonableTwo{
		Username: "batchuser1",
		TweetID:  int64(100000000000010),
		Version:  int64(10010),
		Lat:      1.00002,
		Lon:      -38.11,
		Body:     "hey this thing appears to work, nice!",
	}
	r101 := &ReasonableTwo{
		Username: "batchuser1",
		TweetID:  int64(100000000000010),
		Version:  int64(10011),
		Lat:      1.00002,
		Lon:      -38.11,
		Body:     "hey this thing appears to work, nice!",
	}
	r110 := &ReasonableTwo{
		Username: "batchuser1",
		TweetID:  int64(100000000000011),
		Version:  int64(10010),
		Lat:      1.00002,
		Lon:      -38.11,
		Body:     "hey this thing appears to work, nice!",
	}
	r1s := &ReasonableTwo{}

	b := cp.Batch()
	err = b.DeleteAll(m0, r).Run()
	if err != nil {
		t.Fatal("Error writing:", err)
	}

	result, err := cp.Query(m0).Get("batchtest")
	err = result.Next(rs)
	if err != Done {
		t.Fatal("Expected 0 results (Done) but got:", err)
	}

	b = cp.Batch()
	err = b.Insert(m0, r).Run()
	if err != nil {
		t.Fatal("Error writing:", err)
	}

	result, err = cp.Query(m0).Get("batchtest")
	err = result.Next(rs)
	if err != nil {
		t.Fatal("Expected 1 result but got:", err)
	}
	if !reflect.DeepEqual(r, rs) {
		t.Fatal("Returned row does not match")
	}

	b = cp.Batch()
	err = b.Delete(m0, r).Run()
	if err != nil {
		t.Fatal("Error writing:", err)
	}

	result, err = cp.Query(m0).Get("batchtest")
	err = result.Next(rs)
	if err != Done {
		t.Fatal("Expected 0 results (Done) but got:", err)
	}

	// ---

	b = cp.Batch()
	err = b.DeleteAll(m2, r100).Run()
	if err != nil {
		t.Fatal("Error writing:", err)
	}

	result, err = cp.Query(m2).Get("batchuser1")
	err = result.Next(r1s)
	if err != Done {
		t.Fatal("Expected 0 results (Done) but got:", err)
	}

	b = cp.Batch()
	err = b.Insert(m2, r100).
		Insert(m2, r101).
		Insert(m2, r110).
		Run()
	if err != nil {
		t.Fatal("Error writing:", err)
	}

	result, err = cp.Query(m2).Get("batchuser1")
	err = result.Next(r1s)
	if err != nil {
		t.Fatal("Expected 1 of 3 results but got:", err)
	}
	if !reflect.DeepEqual(r100, r1s) {
		t.Fatal("Returned row does not match")
	}
	err = result.Next(r1s)
	if err != nil {
		t.Fatal("Expected 2 of 3 results but got:", err)
	}
	if !reflect.DeepEqual(r101, r1s) {
		t.Fatal("Returned row does not match")
	}
	err = result.Next(r1s)
	if err != nil {
		t.Fatal("Expected 3 of 3 results but got:", err)
	}
	if !reflect.DeepEqual(r110, r1s) {
		t.Fatal("Returned row does not match")
	}

	b = cp.Batch()
	err = b.Delete(m2, r101).Run()
	if err != nil {
		t.Fatal("Error writing:", err)
	}

	result, err = cp.Query(m2).Get("batchuser1")
	err = result.Next(r1s)
	if err != nil {
		t.Fatal("Expected 1 of 2 results but got:", err)
	}
	if !reflect.DeepEqual(r100, r1s) {
		t.Fatal("Returned row does not match")
	}
	err = result.Next(r1s)
	if err != nil {
		t.Fatal("Expected 2 of 2 results but got:", err)
	}
	if !reflect.DeepEqual(r110, r1s) {
		t.Fatal("Returned row does not match")
	}

}
