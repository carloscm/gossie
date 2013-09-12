package gossie

import (
	"reflect"
	"testing"
)

/*
	todo:
	deeper tests, over more methods, and over all internal types
	compact mapping
*/

type everythingComp struct {
	Key      string `cf:"1" key:"Key" cols:"FBytes,FBool,FInt8,FInt16,FInt32,FInt,FInt64,FFloat32,FFloat64,FString,FUUID"`
	FBytes   []byte
	FBool    bool
	FInt8    int8
	FInt16   int16
	FInt32   int32
	FInt     int
	FInt64   int64
	FFloat32 float32
	FFloat64 float64
	FString  string
	FUUID    UUID
	Val      string
}

type tagsA struct {
	A int `cf:"1" key:"A"`
	B int
	C int
	D int
}

type tagsB struct {
	A int `cf:"1" key:"A" cols:"B"`
	B int `type:"AsciiType"`
	C int `skip:"true"`
	D int `name:"Z"`
}

type tagsC struct {
	A int `cf:"1" key:"A" cols:"B,C" value:"D" mapping:"compact"`
	B int
	C int
	D int
}

type tagsD struct {
	A int `cf:"1" key:"A" cols:"B,C" mapping:"compact"`
	B int
	C int
}

type structTestShell struct {
	mapping        Mapping
	expectedStruct interface{}
	resultStruct   interface{}
	expectedRow    *Row
	name           string
}

func (shell *structTestShell) checkMap(t *testing.T, m Mapping, expectedStruct interface{}, round int) {
	resultRow, err := m.Map(expectedStruct)
	if err != nil {
		t.Error("Error mapping struct: ", err)
	}
	if !reflect.DeepEqual(resultRow, shell.expectedRow) {
		t.Error("(Round ", round, ") Mapped struct ", shell.name, " does not match expected row ", shell.expectedRow, " actual ", resultRow)
	}
}

type testProvider struct {
	row   *Row
	pos   int
	limit int
}

func (t *testProvider) Key() ([]byte, error) {
	return t.row.Key, nil
}

func (t *testProvider) NextColumn() (*Column, error) {
	if t.pos >= len(t.row.Columns) {
		if t.pos >= t.limit {
			return nil, EndAtLimit
		} else {
			return nil, EndBeforeLimit
		}
	}
	c := t.row.Columns[t.pos]
	t.pos++
	return c, nil
}

func (t *testProvider) Rewind() {
	t.pos--
	if t.pos < 0 {
		t.pos = 0
	}
}

func (shell *structTestShell) checkUnmap(t *testing.T, m Mapping) interface{} {
	tp := &testProvider{shell.expectedRow, 0, 10000}
	err := m.Unmap(shell.resultStruct, tp)
	if err != nil {
		t.Error("Error unmapping struct: ", err)
	}
	if !reflect.DeepEqual(shell.resultStruct, shell.expectedStruct) {
		t.Error("Unmapped struct ", shell.name, " does not match expected instance ", shell.expectedStruct, " actual ", shell.resultStruct)
	}
	return shell.resultStruct
}

func (shell *structTestShell) checkFullMap(t *testing.T) {
	shell.checkMap(t, shell.mapping, shell.expectedStruct, 1)
	intermediateStruct := shell.checkUnmap(t, shell.mapping)
	shell.checkMap(t, shell.mapping, intermediateStruct, 2)
}

func TestMap(t *testing.T) {
	mA, _ := NewMapping(&tagsA{})
	mB, _ := NewMapping(&tagsB{})
	mC, _ := NewMapping(&tagsC{})
	mD, _ := NewMapping(&tagsD{})
	mE, _ := NewMapping(&everythingComp{})

	shells := []*structTestShell{
		&structTestShell{
			mapping:        mA,
			name:           "tagsA",
			expectedStruct: &tagsA{1, 2, 3, 4},
			resultStruct:   &tagsA{},
			expectedRow: &Row{
				Key: []byte{0, 0, 0, 0, 0, 0, 0, 1},
				Columns: []*Column{
					&Column{
						Name:  []byte{'B'},
						Value: []byte{0, 0, 0, 0, 0, 0, 0, 2},
					},
					&Column{
						Name:  []byte{'C'},
						Value: []byte{0, 0, 0, 0, 0, 0, 0, 3},
					},
					&Column{
						Name:  []byte{'D'},
						Value: []byte{0, 0, 0, 0, 0, 0, 0, 4},
					},
				},
			},
		},

		&structTestShell{
			mapping:        mB,
			name:           "tagsB",
			expectedStruct: &tagsB{1, 2, 0, 4},
			resultStruct:   &tagsB{},
			expectedRow: &Row{
				Key: []byte{0, 0, 0, 0, 0, 0, 0, 1},
				Columns: []*Column{
					&Column{
						Name:  []byte{0, 1, '2', 0, 0, 1, 'Z', 0},
						Value: []byte{0, 0, 0, 0, 0, 0, 0, 4},
					},
				},
			},
		},

		&structTestShell{
			mapping:        mC,
			name:           "tagsC",
			expectedStruct: &tagsC{1, 2, 3, 4},
			resultStruct:   &tagsC{},
			expectedRow: &Row{
				Key: []byte{0, 0, 0, 0, 0, 0, 0, 1},
				Columns: []*Column{
					&Column{
						Name:  []byte{0, 8, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 3, 0},
						Value: []byte{0, 0, 0, 0, 0, 0, 0, 4},
					},
				},
			},
		},

		&structTestShell{
			mapping:        mD,
			name:           "tagsD",
			expectedStruct: &tagsD{1, 2, 3},
			resultStruct:   &tagsD{},
			expectedRow: &Row{
				Key: []byte{0, 0, 0, 0, 0, 0, 0, 1},
				Columns: []*Column{
					&Column{
						Name:  []byte{0, 8, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 3, 0},
						Value: []byte{},
					},
				},
			},
		},

		&structTestShell{
			mapping: mE,
			name:    "everythingComp",
			expectedStruct: &everythingComp{"a", []byte{1, 2}, true, 3, 4, 5, 6, 7, 8.0, 9.0, "b",
				[16]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}, "c"},
			resultStruct: &everythingComp{},
			expectedRow: &Row{
				Key: []byte{97},
				Columns: []*Column{
					&Column{
						Name: []byte{0, 2, 1, 2, 0, 0, 1, 1, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 8, 0, 0, 0,
							0, 0, 0, 0, 5, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 6, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 4, 65, 0, 0, 0, 0, 0,
							8, 64, 34, 0, 0, 0, 0, 0, 0, 0, 0, 1, 98, 0, 0, 16, 0, 17, 34, 51, 68, 85, 102, 119, 136, 153, 170, 187, 204,
							221, 238, 255, 0, 0, 3, 86, 97, 108, 0},
						Value: []byte{99},
					},
				},
			},
		},
	}

	for _, shell := range shells {
		shell.checkFullMap(t)
	}
}

var testMaps = []map[string]interface{}{
	{
		"Id":       "Hello",
		"Age":      25,
		"Whatever": 0.2,
	},
	// Testing all gossie supported times.
	{
		"Id":         "34e8f7-c45wc-c45w45cdx",
		"AnInt":      25,
		"AnInt8":     int8(5),
		"AnInt16":    int16(5),
		"AnInt32":    int32(5),
		"AnInt64":    int64(5),
		"AFloat":     1.45,
		"AFloat32":   float32(1.5),
		"AFloat64":   float64(1.5),
		"AString":    "Heyy",
		"ABool":      false,
		"AByteSlice": []byte("Ahoi"),
	},
}

// We are using the serialized maps themselves as
// a scheme for unserializing.
func TestMapToRow(t *testing.T) {
	for _, m := range testMaps {
		id := m["Id"].(string)
		row, err := MapToRow("Id", m)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(row.Key, []byte(id)) {
			t.Fatal("Row key is incorrect", string(row.Key), m["Id"])
		}
		// For some reason cassandra returns empty columns too?!
		// TODO: investigate.
		//if len(row.Columns) != len(m)-1 {
		//	t.Fatal(len(row.Columns), len(m))
		//}
		m1, err := RowToMap("Id", m, row)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(m, m1) {
			t.Fatal(m, m1)
		}
	}
}
