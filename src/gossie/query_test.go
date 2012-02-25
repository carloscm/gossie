package gossie

import (
	"testing"
	"fmt"
)

func TestQuery(t *testing.T) {

	cp, err := NewConnectionPool([]string{"127.0.0.1:9160"}, "TestGossie", PoolOptions{Size:1,Timeout:1000})
	if err != nil {
		t.Fatal("Error connecting to Cassandra:", err)
	}

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

	cp.Close()
}

func buildAllTypesRow(key string) *Row {
	var row Row
	row.Key, _ = Marshal(key, BytesType)
	v1, _ := Marshal([]byte{1,2,3}, BytesType)
	v2, _ := Marshal("hi!", AsciiType)
	v3, _ := Marshal("leña al fuego", UTF8Type)
	v4, _ := Marshal(int64(1e15), LongType)
	u, _ := NewUUID("00112233-4455-6677-8899-aabbccddeeff")
	v5, _ := Marshal(u, UUIDType)
	v6, _ := Marshal(true, BooleanType)
	v7, _ := Marshal(float32(-1.1), FloatType)
	v8, _ := Marshal(float64(-1.00000000000000000000001), DoubleType)
	row.Columns = []*Column{
		&Column{Name:[]byte("colBytesType"), Value:v1},
		&Column{Name:[]byte("colAsciiType"), Value:v2},
		&Column{Name:[]byte("colUTF8Type"), Value:v3},
		&Column{Name:[]byte("colLongType"), Value:v4},
		&Column{Name:[]byte("colUUIDType"), Value:v5},
		&Column{Name:[]byte("colBooleanType"), Value:v6},
		&Column{Name:[]byte("colFloatType"), Value:v7},
		&Column{Name:[]byte("colDoubleType"), Value:v8},
		//&Column{Name:[]byte("colDateType"), Value:
	}
	return &row
}

func TestMutation(t *testing.T) {

	cp, err := NewConnectionPool([]string{"127.0.0.1:9160"}, "TestGossie", PoolOptions{Size:1,Timeout:1000})
	if err != nil {
		t.Fatal("Error connecting to Cassandra:", err)
	}

	m := cp.Mutation()
	for i := 0; i < 3; i++ {
		m.Insert("AllTypes", buildAllTypesRow(fmt.Sprint("row", i)))
	}
	err = m.Run()
	t.Log(err)

	//t.Fatal("wut")

	cp.Close()
}

func BenchmarkGetOne(b *testing.B) {
	b.StopTimer()
	cp, _ := NewConnectionPool([]string{"127.0.0.1:9160"}, "TestGossie", PoolOptions{Size:1,Timeout:1000})
	b.StartTimer()
    for i := 0; i < b.N; i++ {
		cp.Query().Cf("AllTypes").Key([]byte("a")).GetOne()
    }
}

func BenchmarkInsert(b *testing.B) {
	b.StopTimer()
	cp, _ := NewConnectionPool([]string{"127.0.0.1:9160"}, "TestGossie", PoolOptions{Size:1,Timeout:1000})
	row := buildAllTypesRow("row")
	b.StartTimer()
    for i := 0; i < b.N; i++ {
    	row.Key, _ = Marshal(fmt.Sprint("row", i), BytesType)
		cp.Mutation().Insert("AllTypes", row).Run()
    }
}
