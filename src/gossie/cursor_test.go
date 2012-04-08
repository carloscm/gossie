package gossie

import (
	"reflect"
	"testing"
)

type ReasonableOne struct {
	Username string `cf:"Reasonable" key:"Username" col:"TweetID,*name" val:"*value"`
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

	cursor := cp.Cursor(ro)
	err = cursor.Write()
	if err != nil {
		t.Error("Writing struct:", err)
	}

	ro2 := &ReasonableOne{
		Username: "testuser",
		TweetID:  100000000000002,
	}

	cursor2 := cp.Cursor(ro2)
	err = cursor2.Read(1)
	if err != nil {
		t.Error("Reading struct:", err)
	}

	if !reflect.DeepEqual(ro, ro2) {
		t.Error("Read does not match Write")
	}

}
