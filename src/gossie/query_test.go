package gossie

import (
	"reflect"
	"testing"
)

/*
todo:
	refactor tests into something resembling maintenable code
*/

type ReasonableZero struct {
	Username string `cf:"ReasonableZero" key:"Username"`
	Lat      float32
	Lon      float32
	Body     string
}

type ReasonableOne struct {
	Username string `cf:"ReasonableOne" key:"Username" cols:"TweetID"`
	TweetID  int64
	Lat      float32
	Lon      float32
	Body     string
}

type ReasonableTwo struct {
	Username string `cf:"ReasonableTwo" key:"Username" cols:"TweetID,Version"`
	TweetID  int64
	Version  int64
	Lat      float32
	Lon      float32
	Body     string
}

func TestQueryGet(t *testing.T) {
	cp, err := NewConnectionPool([]string{"127.0.0.1:9160"}, "TestGossie", PoolOptions{Size: 1, Timeout: 1000})
	if err != nil {
		t.Fatal("Error connecting to Cassandra:", err)
	}

	m0, err := NewMapping(&ReasonableZero{})
	if err != nil {
		t.Fatal("Error building mapping:", err)
	}
	m1, err := NewMapping(&ReasonableOne{})
	if err != nil {
		t.Fatal("Error building mapping:", err)
	}
	m2, err := NewMapping(&ReasonableTwo{})
	if err != nil {
		t.Fatal("Error building mapping:", err)
	}

	w := cp.Writer()
	r := &ReasonableZero{"testuser", 1.00002, -38.11, "hey this thing appears to work, nice!"}
	row, err := m0.Map(r)
	if err != nil {
		t.Fatal("Error mapping:", err)
	}
	w.Insert("ReasonableZero", row)
	err = w.Run()
	if err != nil {
		t.Fatal("Error writing:", err)
	}

	w = cp.Writer()
	for i := 0; i < 100; i++ {
		r := &ReasonableOne{
			Username: "testuser",
			TweetID:  int64(100000000000000) + int64(i),
			Lat:      1.00002,
			Lon:      -38.11,
			Body:     "hey this thing appears to work, nice!",
		}
		row, err := m1.Map(r)
		if err != nil {
			t.Fatal("Error mapping:", err)
		}
		w.Insert("ReasonableOne", row)
	}
	err = w.Run()
	if err != nil {
		t.Fatal("Error writing:", err)
	}

	w = cp.Writer()
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			r := &ReasonableTwo{
				Username: "testuser",
				TweetID:  int64(100000000000000) + int64(i),
				Version:  int64(i*10000 + j),
				Lat:      1.00002,
				Lon:      -38.11,
				Body:     "hey this thing appears to work, nice!",
			}
			row, err := m2.Map(r)
			if err != nil {
				t.Fatal("Error mapping:", err)
			}
			w.Insert("ReasonableTwo", row)
		}
	}
	err = w.Run()
	if err != nil {
		t.Fatal("Error writing:", err)
	}

	/////

	q0 := cp.Query(m0)
	r0 := &ReasonableZero{}

	res, err := q0.Get("nope")
	if err != nil {
		t.Fatal("Query get error:", err)
	}
	err = res.Next(r0)
	if err != Done {
		t.Fatal("Result Next is not Done:", err)
	}

	res, err = q0.Get("testuser")
	if err != nil {
		t.Fatal("Query get error:", err)
	}
	err = res.Next(r0)
	if err != nil {
		t.Fatal("Result next error:", err)
	}
	if !reflect.DeepEqual(r0, &ReasonableZero{"testuser", 1.00002, -38.11, "hey this thing appears to work, nice!"}) {
		t.Error("Read does not match Write")
	}
	err = res.Next(r0)
	if err != Done {
		t.Fatal("Result Next is not Done:", err)
	}

	/////

	q1 := cp.Query(m1)
	r1 := &ReasonableOne{}

	res, err = q1.Get("nope")
	if err != nil {
		t.Fatal("Query get error:", err)
	}
	err = res.Next(r1)
	if err != Done {
		t.Fatal("Result Next is not Done:", err)
	}

	res, err = q1.Get("testuser")
	if err != nil {
		t.Fatal("Query get error:", err)
	}
	for i := 0; i < 100; i++ {
		r := &ReasonableOne{
			Username: "testuser",
			TweetID:  int64(100000000000000) + int64(i),
			Lat:      1.00002,
			Lon:      -38.11,
			Body:     "hey this thing appears to work, nice!",
		}
		err = res.Next(r1)
		if err != nil {
			t.Fatal("Result next error:", err)
		}
		if !reflect.DeepEqual(r1, r) {
			t.Error("Read does not match Write")
		}
	}
	err = res.Next(r1)
	if err != Done {
		t.Fatal("Result Next is not Done:", err)
	}

	res, err = q1.GetBetween("testuser", int64(100000000000050), int64(100000000000070))
	if err != nil {
		t.Fatal("Query get error:", err)
	}
	for i := 50; i < 70; i++ {
		r := &ReasonableOne{
			Username: "testuser",
			TweetID:  int64(100000000000000) + int64(i),
			Lat:      1.00002,
			Lon:      -38.11,
			Body:     "hey this thing appears to work, nice!",
		}
		err = res.Next(r1)
		if err != nil {
			t.Fatal("Result next error:", err)
		}
		if !reflect.DeepEqual(r1, r) {
			t.Log(r1)
			t.Error("Read does not match Write")
		}
	}
	err = res.Next(r1)
	if err != Done {
		t.Fatal("Result Next is not Done:", err)
	}

	res, err = q1.Get("testuser", int64(100000000000010))
	if err != nil {
		t.Fatal("Query get error:", err)
	}
	err = res.Next(r1)
	if err != nil {
		t.Fatal("Result next error:", err)
	}
	if !reflect.DeepEqual(r1, &ReasonableOne{"testuser", int64(100000000000010), 1.00002, -38.11, "hey this thing appears to work, nice!"}) {
		t.Error("Read does not match Write")
	}
	err = res.Next(r1)
	if err != Done {
		t.Fatal("Result Next is not Done:", err)
	}

	/////

	q2 := cp.Query(m2)
	r2 := &ReasonableTwo{}

	res, err = q2.Get("nope")
	if err != nil {
		t.Fatal("Query get error:", err)
	}
	err = res.Next(r2)
	if err != Done {
		t.Fatal("Result Next is not Done:", err)
	}

	res, err = q2.Get("testuser")
	if err != nil {
		t.Fatal("Query get error:", err)
	}
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			r := &ReasonableTwo{
				Username: "testuser",
				TweetID:  int64(100000000000000) + int64(i),
				Version:  int64(i*10000 + j),
				Lat:      1.00002,
				Lon:      -38.11,
				Body:     "hey this thing appears to work, nice!",
			}
			err = res.Next(r2)
			if err != nil {
				t.Fatal("Result next error:", err)
			}
			if !reflect.DeepEqual(r2, r) {
				t.Error("Read does not match Write")
			}
		}
	}
	err = res.Next(r2)
	if err != Done {
		t.Fatal("Result Next is not Done:", err)
	}

	res, err = q2.GetBetween("testuser", int64(100000000000002), int64(100000000000007))
	if err != nil {
		t.Fatal("Query get error:", err)
	}
	for i := 2; i < 7; i++ {
		for j := 0; j < 10; j++ {
			r := &ReasonableTwo{
				Username: "testuser",
				TweetID:  int64(100000000000000) + int64(i),
				Version:  int64(i*10000 + j),
				Lat:      1.00002,
				Lon:      -38.11,
				Body:     "hey this thing appears to work, nice!",
			}
			err = res.Next(r2)
			if err != nil {
				t.Fatal("Result next error:", err)
			}
			if !reflect.DeepEqual(r2, r) {
				t.Error("Read does not match Write")
			}
		}
	}
	err = res.Next(r2)
	if err != Done {
		t.Fatal("Result Next is not Done:", err)
	}

	res, err = q2.Get("testuser", int64(100000000000003))
	if err != nil {
		t.Fatal("Query get error:", err)
	}
	for j := 0; j < 10; j++ {
		r := &ReasonableTwo{
			Username: "testuser",
			TweetID:  int64(100000000000003),
			Version:  int64(30000 + j),
			Lat:      1.00002,
			Lon:      -38.11,
			Body:     "hey this thing appears to work, nice!",
		}
		err = res.Next(r2)
		if err != nil {
			t.Fatal("Result next error:", err)
		}
		if !reflect.DeepEqual(r2, r) {
			t.Error("Read does not match Write")
		}
	}
	err = res.Next(r2)
	if err != Done {
		t.Fatal("Result Next is not Done:", err)
	}

	res, err = q2.GetBetween("testuser", int64(100000000000003), int64(30006), int64(30009))
	if err != nil {
		t.Fatal("Query get error:", err)
	}
	for j := 6; j < 9; j++ {
		r := &ReasonableTwo{
			Username: "testuser",
			TweetID:  int64(100000000000003),
			Version:  int64(30000 + j),
			Lat:      1.00002,
			Lon:      -38.11,
			Body:     "hey this thing appears to work, nice!",
		}
		err = res.Next(r2)
		if err != nil {
			t.Fatal("Result next error:", err)
		}
		if !reflect.DeepEqual(r2, r) {
			t.Error("Read does not match Write")
		}
	}
	err = res.Next(r2)
	if err != Done {
		t.Fatal("Result Next is not Done:", err)
	}

	res, err = q2.Get("testuser", int64(100000000000005), int64(50003))
	if err != nil {
		t.Fatal("Query get error:", err)
	}
	err = res.Next(r2)
	if err != nil {
		t.Fatal("Result next error:", err)
	}
	if !reflect.DeepEqual(r2, &ReasonableTwo{"testuser", 100000000000005, 50003, 1.00002, -38.11, "hey this thing appears to work, nice!"}) {
		t.Error("Read does not match Write")
	}
	err = res.Next(r2)
	if err != Done {
		t.Fatal("Result Next is not Done:", err)
	}

}
