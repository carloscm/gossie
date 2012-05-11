package gossie

import (
	"reflect"
	"testing"
)

/*
todo:
	real tests
*/

type ReasonableOne struct {
	Username string `cf:"Reasonable" key:"Username" cols:"TweetID"`
	TweetID  int64
	Lat      float32
	Lon      float32
	Body     string
}

func TestRead(t *testing.T) {
	cp, err := NewConnectionPool([]string{"127.0.0.1:9160"}, "TestGossie", PoolOptions{Size: 1, Timeout: 1000})
	if err != nil {
		t.Fatal("Error connecting to Cassandra:", err)
	}

	m, _ := NewMapping(&ReasonableOne{})

	ro := &ReasonableOne{
		Username: "testuser",
		TweetID:  100000000000002,
		Lat:      1.00002,
		Lon:      -38.11,
		Body:     "hey this thing appears to work, nice!",
	}
	r, _ := m.Map(ro)
	cp.Writer().Insert("Reasonable", r).Run()

	ro2 := &ReasonableOne{
		Username: "testuser",
		TweetID:  100000000000003,
		Lat:      2.00002,
		Lon:      1.11,
		Body:     "more words",
	}
	r, _ = m.Map(ro2)
	cp.Writer().Insert("Reasonable", r).Run()

	q := cp.Query(m)

	res, err := q.Get("testuser", int64(100000000000002))
	if err != nil {
		t.Fatal("Query get error:", err)
	}
	ro3 := &ReasonableOne{}
	err = res.Next(ro3)
	if err != nil {
		t.Fatal("Result next error:", err)
	}
	if !reflect.DeepEqual(ro, ro3) {
		t.Error("Read does not match Write")
	}
	ro3 = &ReasonableOne{}
	err = res.Next(ro3)
	if err != Done {
		t.Log(ro3)
		t.Fatal("Result Next is not Done:", err)
	}

	res2, err := q.Get("testuser")
	if err != nil {
		t.Fatal("Query get error:", err)
	}
	ro4 := &ReasonableOne{}
	err = res2.Next(ro4)
	if err != nil {
		t.Fatal("Result next error:", err)
	}
	if !reflect.DeepEqual(ro, ro4) {
		t.Error("Read does not match Write")
	}
	err = res2.Next(ro4)
	if err != nil {
		t.Fatal("Result next error:", err)
	}
	if !reflect.DeepEqual(ro2, ro4) {
		t.Error("Read does not match Write")
	}
	err = res.Next(ro3)
	if err != Done {
		t.Fatal("Result Next is not Done:", err)
	}

}
