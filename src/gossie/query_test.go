package gossie

import (
	"reflect"
	"testing"
	"time"
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

type Timeseries struct {
	Username string `cf:"Timeseries" key:"Username" cols:"TimeUUID"`
	TimeUUID UUID
	Seq      int
}

type ReasonableTwo struct {
	Username string `cf:"ReasonableTwo" key:"Username" cols:"TweetID,Version"`
	TweetID  int64
	Version  int64
	Lat      float32
	Lon      float32
	Body     string
}

type CompositeFull struct {
	Username    string `mapping:"compact" cf:"CompositeFull" key:"Username" cols:"Bytes,Ascii,UTF8,Long,Int32,UUIDf,LexicalUUID,TimeUUID,Boolean,Float,Double,Date" value:"Value"`
	Bytes       []byte
	Ascii       string
	UTF8        string
	Long        int64
	Int32       int `type:"Int32Type"`
	UUIDf       UUID
	LexicalUUID UUID
	TimeUUID    UUID
	Boolean     bool
	Float       float32
	Double      float64
	Date        time.Time
	Value       int
}

func createTimeseries(t *testing.T, cp ConnectionPool) int {
	cp.Writer().Delete("Timeseries", []byte("testuser")).Run()

	mT, err := NewMapping(&Timeseries{})

	// pseudo random number based on the sub-ms part of the current unix nano time
	seqBase := int(time.Now().UnixNano() % 1e6)

	w := cp.Writer()
	for i := 0; i < 100; i++ {
		u, err := NewTimeUUID()
		if err != nil {
			t.Fatal("Error generating TimeUUID:", err)
		}
		r := &Timeseries{
			Username: "testuser",
			TimeUUID: u,
			Seq:      seqBase + i,
		}
		row, err := mT.Map(r)
		if err != nil {
			t.Fatal("Error mapping:", err)
		}
		w.InsertTtl("Timeseries", row, 60)
	}
	err = w.Run()
	if err != nil {
		t.Fatal("Error writing:", err)
	}
	return seqBase
}

func createCompositeFull(t *testing.T, cp ConnectionPool) int {
	cp.Writer().Delete("CompositeFull", []byte("testuser")).Run()

	mC, err := NewMapping(&CompositeFull{})

	// pseudo random number based on the sub-ms part of the current unix nano time
	seqBase := int(time.Now().UnixNano() % 1e6)

	w := cp.Writer()

	uR1, err := NewRandomUUID()
	if err != nil {
		t.Fatal("Error generating random UUID:", err)
	}
	uR2, err := NewRandomUUID()
	if err != nil {
		t.Fatal("Error generating random UUID:", err)
	}
	for i := 0; i < 100; i++ {
		uT, err := NewTimeUUID()
		if err != nil {
			t.Fatal("Error generating TimeUUID:", err)
		}
		r := &CompositeFull{
			Username:    "testuser",
			Bytes:       []byte("yep"),
			Ascii:       "good",
			UTF8:        "works",
			Long:        int64(i / 10),
			Int32:       -42,
			UUIDf:       uR1,
			LexicalUUID: uR2,
			TimeUUID:    uT,
			Boolean:     true,
			Float:       42.42,
			Double:      4e42,
			Date:        time.Now(),
			Value:       seqBase + i,
		}
		row, err := mC.Map(r)
		if err != nil {
			t.Fatal("Error mapping:", err)
		}
		w.InsertTtl("CompositeFull", row, 60)
	}
	err = w.Run()
	if err != nil {
		t.Fatal("Error writing:", err)
	}
	return seqBase
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
	mT, err := NewMapping(&Timeseries{})
	if err != nil {
		t.Fatal("Error building mapping:", err)
	}
	mC, err := NewMapping(&CompositeFull{})
	if err != nil {
		t.Fatal("Error building mapping:", err)
	}

	w := cp.Writer()
	r := &ReasonableZero{"testuser", 1.00002, -38.11, "hey this thing appears to work, nice!"}
	row, err := m0.Map(r)
	if err != nil {
		t.Fatal("Error mapping:", err)
	}
	w.InsertTtl("ReasonableZero", row, 60)
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
		w.InsertTtl("ReasonableOne", row, 60)
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
			w.InsertTtl("ReasonableTwo", row, 60)
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

	seqBase := createTimeseries(t, cp)
	qT := cp.Query(mT)
	rT := &Timeseries{}

	res, err = qT.Get("testuser")
	if err != nil {
		t.Fatal("Query get error:", err)
	}
	for i := 99; i >= 0; i-- {
		err = res.Next(rT)
		if err != nil {
			t.Fatal("Result next error:", err)
		}
		if rT.Seq != seqBase+i {
			t.Log(rT)
			t.Error("Read does not match Write")
		}
	}

	seqBase = createTimeseries(t, cp)

	qT.Reversed(true)

	res, err = qT.Get("testuser")
	if err != nil {
		t.Fatal("Query get error:", err)
	}
	for i := 0; i < 100; i++ {
		err = res.Next(rT)
		if err != nil {
			t.Fatal("Result next error:", err)
		}
		if rT.Seq != seqBase+i {
			t.Log(seqBase + i)
			t.Log(rT)
			t.Error("Read does not match Write")
		}
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

	/////

	seqBase = createCompositeFull(t, cp)

	qC := cp.Query(mC)
	rC := &CompositeFull{}

	res, err = qC.Get("testuser")
	if err != nil {
		t.Fatal("Query get error:", err)
	}
	for i := 0; i < 100; i++ {
		err = res.Next(rC)
		if err != nil {
			t.Fatal("Result next error:", err)
		}
		if rC.Value != seqBase+i {
			t.Log(seqBase + i)
			t.Log(rT)
			t.Error("Read does not match Write")
		}
	}

}
