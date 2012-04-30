package gossie

import (
//"reflect"
//"testing"
)

/*

todo:

	since most of the Cursor interface is still in flux the current tests are minimal

*/

/*
type ReasonableOne struct {
	Username string `cf:"Reasonable" key:"Username,TweetID"`
	TweetID  int64
	Lat      float32
	Lon      float32
	Body     string
	Cursor
}

func TestRead(t *testing.T) {
	cp, err := NewConnectionPool([]string{"127.0.0.1:9160"}, "TestGossie", PoolOptions{Size: 1, Timeout: 1000})
	if err != nil {
		t.Fatal("Error connecting to Cassandra:", err)
	}

	ro := &ReasonableOne{
		Username: "testuser",
		TweetID:  100000000000002,
		Lat:      1.00002,
		Lon:      -38.11,
		Body:     "hey this thing appears to work, nice!",
	}

	cursor := cp.Cursor()
	err = cursor.Write(ro)
	if err != nil {
		t.Error("Writing struct:", err)
	}

	ro2 := &ReasonableOne{
		Username: "testuser",
		TweetID:  100000000000003,
		Lat:      2.00002,
		Lon:      1.11,
		Body:     "more words",
	}
	err = cursor.Write(ro2)
	if err != nil {
		t.Error("Writing struct:", err)
	}

	ro3 := &ReasonableOne{
		Username: "testuser",
		TweetID:  100000000000002,
	}
	err = cursor.Read(ro3)
	if err != nil {
		t.Fatal("Reading struct:", err)
	}
	if !reflect.DeepEqual(ro, ro3) {
		t.Error("Read does not match Write")
	}

	ro3 = &ReasonableOne{
		Username: "testuser",
		TweetID:  100000000000003,
	}
	err = cursor.Read(ro3)
	if err != nil {
		t.Fatal("Reading struct:", err)
	}
	if !reflect.DeepEqual(ro2, ro3) {
		t.Error("Read does not match Write")
	}
}
*/
