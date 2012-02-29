package gossie

import (
	"testing"
	"fmt"
	"reflect"
)

type testColumn struct {
	name interface{}
	nameType TypeDesc
	value interface{}
	valueType TypeDesc
}

type testRow struct {
	key interface{}
	keyType TypeDesc
	columns []testColumn
}

func addColumn(row *Row, c *testColumn) {
	n, _ := Marshal(c.name, c.nameType)
	v, _ := Marshal(c.value, c.valueType)
	row.Columns = append(row.Columns, &Column{Name:n,Value:v})
}

func buildRow(r *testRow) *Row {
	var row Row
	row.Key, _ = Marshal(r.key, r.keyType)
	for _, c := range r.columns {
		addColumn(&row, &c)
	}
	return &row
}

func buildCounterTestRow(key string) *testRow {
	return &testRow{
		key:key, keyType:BytesType, columns:[]testColumn{
			testColumn{"fortytwo",	AsciiType,		int64(-42),		LongType},
			testColumn{"one",		AsciiType,		int64(1),		LongType},
			testColumn{"wtf",		AsciiType,		int64(1e15),	LongType},
		},
	}
}

func buildCounterRow(key string) *Row {
	return buildRow(buildCounterTestRow(key))
}

func buildAllTypesTestRow(key string) *testRow {
	u, _ := NewUUID("00112233-4455-6677-8899-aabbccddeeff")
	return &testRow{
		key:key, keyType:BytesType, columns:[]testColumn{
			testColumn{"colAsciiType",	AsciiType,	"hi!",			AsciiType},
			testColumn{"colBooleanType",AsciiType,	true,			BooleanType},
			testColumn{"colBytesType",	AsciiType,	[]byte{1,2,3},	BytesType},
			testColumn{"colDoubleType",	AsciiType,	float64(-1.1),	DoubleType},
			testColumn{"colFloatType",	AsciiType,	float32(1.1),	FloatType},
			testColumn{"colLongType",	AsciiType,	int64(1e15),	LongType},
			testColumn{"colUTF8Type",	AsciiType,	"leña al fuego",UTF8Type},
			testColumn{"colUUIDType",	AsciiType,	u, UUIDType},
		},
	}
}

func buildAllTypesAfterDeletesTestRow(key string) *testRow {
	u, _ := NewUUID("00112233-4455-6677-8899-aabbccddeeff")
	return &testRow{
		key:key, keyType:BytesType, columns:[]testColumn{
			testColumn{"colBytesType",	AsciiType,	[]byte{1,2,3},	BytesType},
			testColumn{"colDoubleType",	AsciiType,	float64(-1.1),	DoubleType},
			testColumn{"colFloatType",	AsciiType,	float32(1.1),	FloatType},
			testColumn{"colLongType",	AsciiType,	int64(1e15),	LongType},
			testColumn{"colUTF8Type",	AsciiType,	"leña al fuego",UTF8Type},
			testColumn{"colUUIDType",	AsciiType,	u, UUIDType},
		},
	}
}

func buildAllTypesRow(key string) *Row {
	return buildRow(buildAllTypesTestRow(key))
}

func checkRow(t *testing.T, expected *testRow, actual *Row) {
	exKey, _ := Marshal(expected.key, expected.keyType)
	if !reflect.DeepEqual(exKey, actual.Key) {
		t.Error("Keys differ: ", expected.key, " vs ", actual.Key)
	}
	if len(expected.columns) != len(actual.Columns) {
		t.Fatal("Number of columns differ: ", len(expected.columns), " vs ", len(actual.Columns))
	}
	for i, c := range expected.columns {
		exName, _ := Marshal(c.name, c.nameType)
		exValue, _ := Marshal(c.value, c.valueType)
		if !reflect.DeepEqual(exName, actual.Columns[i].Name) {
			t.Error("Column index ", i, ", named ", c.name, ", is not named the same in actual row: ", actual.Columns[i].Name)
		}
		if !reflect.DeepEqual(exValue, actual.Columns[i].Value) {
			t.Error("Column index ", i, ", named ", c.name, ", has not the expected value: ", exValue, " vs ", actual.Columns[i].Value)
		}
	}
}

func buildIntSliceFromRow(row *Row) []int64 {
	var r []int64
	for _, c := range row.Columns {
		var v int64
		Unmarshal(c.Value, LongType, &v)
		r = append(r, v)
	}
	return r
}

func TestMutationAndQuery(t *testing.T) {
	cp, err := NewConnectionPool([]string{"127.0.0.1:9160"}, "TestGossie", PoolOptions{Size:1,Timeout:1000})
	if err != nil {
		t.Fatal("Error connecting to Cassandra:", err)
	}

	m := cp.Mutation()
	for i := 0; i < 3; i++ {
		key := fmt.Sprint("row", i)
		m.Insert("AllTypes", buildAllTypesRow(key))
		m.DeltaCounters("Counters", buildCounterRow(key))
	}
	err = m.Run()
	if err != nil {
		t.Error("Error running mutation: ", err)
	}

	for i := 0; i < 3; i++ {
		key := fmt.Sprint("row", i)
		row, err := cp.Query().Cf("AllTypes").Get([]byte(key))
		if err != nil {
			t.Error("Error running query: ", err)
		}
		checkRow(t, buildAllTypesTestRow(key), row)
	}

	row, err := cp.Query().Cf("Counters").Get([]byte("row0"))
	if err != nil {
		t.Error("Error running query: ", err)
	}
	if row == nil {
		t.Fatal("An expected row was not returned")
	}
	before := buildIntSliceFromRow(row)

	err = cp.Mutation().DeltaCounters("Counters", buildCounterRow("row0")).Run()
	if err != nil {
		t.Error("Error running query: ", err)
	}
	row, err = cp.Query().Cf("Counters").Get([]byte("row0"))
	if err != nil {
		t.Error("Error running query: ", err)
	}
	if row == nil {
		t.Fatal("An expected row was not returned")
	}
	after := buildIntSliceFromRow(row)

	if  after[0] - before[0] != int64(-42) ||
		after[1] - before[1] != int64(1) ||
		after[2] - before[2] != int64(1e15) {
		t.Error("Counter row was not updated as expected: ", before, " to ", after)
	}

	m = cp.Mutation()
	m.Delete("AllTypes", []byte("row0"))
	m.Delete("Counters", []byte("row0"))
	m.DeleteColumns("AllTypes", []byte("row1"), [][]byte{[]byte("colBooleanType"),[]byte("colAsciiType")})
	m.DeleteColumns("Counters", []byte("row1"), [][]byte{[]byte("one"),[]byte("wtf")})
	err = m.Run()
	if err != nil {
		t.Error("Error running mutation: ", err)
	}

	row, err = cp.Query().Cf("AllTypes").Get([]byte("rowNo"))
	if err != nil {
		t.Error("Error running query: ", err)
	}
	if row != nil {
		t.Error("An expected unexisting row was returned: ", row)
	}

	row, err = cp.Query().Cf("AllTypes").Get([]byte("row0"))
	if err != nil {
		t.Error("Error running query: ", err)
	}
	if row != nil {
		t.Error("An expected deleted row was returned: ", row)
	}

	count, err := cp.Query().Cf("AllTypes").Count([]byte("row0"))
	if err != nil {
		t.Error("Error running query: ", err)
	}
	if count > 0 {
		t.Error("An expected deleted row had higher than 0 columns counted: ", count)
	}

	row, err = cp.Query().Cf("Counters").Get([]byte("row0"))
	if err != nil {
		t.Error("Error running query: ", err)
	}
	if row != nil {
		t.Error("An expected deleted row was returned: ", row)
	}

	row, err = cp.Query().Cf("AllTypes").Get([]byte("row1"))
	if err != nil {
		t.Error("Error running query: ", err)
	}
	checkRow(t, buildAllTypesAfterDeletesTestRow("row1"), row)

	count, err = cp.Query().Cf("AllTypes").Count([]byte("row1"))
	if err != nil {
		t.Error("Error running query: ", err)
	}
	if count != 6 {
		t.Error("A row had an unexpected column count ", count)
	}

	row, err = cp.Query().Cf("AllTypes").Get([]byte("row2"))
	if err != nil {
		t.Error("Error running query: ", err)
	}
	checkRow(t, buildAllTypesTestRow("row2"), row)

	count, err = cp.Query().Cf("AllTypes").Count([]byte("row2"))
	if err != nil {
		t.Error("Error running query: ", err)
	}
	if count != 8 {
		t.Error("A row had an unexpected column count ", count)
	}

	rows, err := cp.Query().Cf("AllTypes").MultiGet([][]byte{[]byte("rowNo1"),[]byte("rowNo2"),[]byte("rowNo3")})
	if err != nil {
		t.Error("Error running query: ", err)
	}
	if rows == nil {
		t.Error("Expected a result in MultiGet call, even with 0 expected results: ", rows)
	}
	if len(rows) != 0 {
		t.Error("Expected 0 rows in MultiGet call, got ", len(rows))
	}

	rows, err = cp.Query().Cf("AllTypes").MultiGet([][]byte{[]byte("row0"),[]byte("row1"),[]byte("row2")})
	if err != nil {
		t.Error("Error running query: ", err)
	}
	if rows == nil {
		t.Fatal("Expected a result in MultiGet call")
	}
	if len(rows) != 2 {
		t.Error("Expected 2 rows in MultiGet call, got ", len(rows))
	}
	for _, row := range rows {
		k := string(row.Key)
		if k == "row2" {
			checkRow(t, buildAllTypesTestRow("row2"), row)
		} else if k == "row1" {
			checkRow(t, buildAllTypesAfterDeletesTestRow("row1"), row)
		} else {
			t.Error("Unexpected row returned in MultiGet call: ", k)
		}
	}

	rows, err = cp.Query().Cf("AllTypes").RangeGet(&Range{Count:1000})
	if err != nil {
		t.Error("Error running query: ", err)
	}
	if rows == nil {
		t.Fatal("Expected a result in RangeGet call")
	}
	if len(rows) != 2 {
		t.Error("Expected 2 rows in RangeGet call, got ", len(rows))
	}
	for _, row := range rows {
		k := string(row.Key)
		if k == "row2" {
			checkRow(t, buildAllTypesTestRow("row2"), row)
		} else if k == "row1" {
			checkRow(t, buildAllTypesAfterDeletesTestRow("row1"), row)
		} else {
			t.Error("Unexpected row returned in RangeGet call: ", k)
		}
	}

	cp.Close()
}

func BenchmarkGet(b *testing.B) {
	b.StopTimer()
	cp, _ := NewConnectionPool([]string{"127.0.0.1:9160"}, "TestGossie", PoolOptions{Size:1,Timeout:1000})
	b.StartTimer()
    for i := 0; i < b.N; i++ {
		cp.Query().Cf("AllTypes").Get([]byte("a"))
    }
}

func BenchmarkInsert(b *testing.B) {
	b.StopTimer()
	cp, _ := NewConnectionPool([]string{"127.0.0.1:9160"}, "TestGossie", PoolOptions{Size:1,Timeout:1000})
	row := buildAllTypesRow("row")
	b.StartTimer()
    for i := 0; i < b.N; i++ {
		m := cp.Mutation()
		for j := 0; j < 10; j++ {
    		row.Key, _ = Marshal(fmt.Sprint("row", i+j), BytesType)
			m.Insert("AllTypes", row)
		}
    	m.Run()
    }
}
