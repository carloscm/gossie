package gossie

import (
	"reflect"
	"testing"
)

type everythingComp struct {
	Key      string
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
		t.Log(shell.expectedRow.Columns[0])
		t.Log(shell.expectedRow.Columns[1])
		t.Log(resultRow.Columns[0])
		t.Log(resultRow.Columns[1])
		t.Error("(Round ", round, ") Mapped struct ", shell.name, " does not match expected row ", shell.expectedRow, " actual ", resultRow)
	}
}

func (shell *structTestShell) checkUnmap(t *testing.T, m Mapping) interface{} {
	n, err := m.Unmap(shell.resultStruct, 0, shell.expectedRow)
	if err != nil {
		t.Error("Error unmapping struct: ", err)
	}
	if n != len(shell.expectedRow.Columns) {
		t.Error("Wrong number of columns consumed when unmapping struct")
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

	shells := []*structTestShell{
		&structTestShell{
			mapping:        NewSparse("cfname", "A"),
			name:           "noErrA",
			expectedStruct: &noErrA{1, 2, 3},
			resultStruct:   &noErrA{},
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
				},
			},
		},

		&structTestShell{
			mapping:        NewSparse("cfname", "A", "B"),
			name:           "noErrB",
			expectedStruct: &noErrB{1, 2, 3, 4},
			resultStruct:   &noErrB{},
			expectedRow: &Row{
				Key: []byte{0, 0, 0, 0, 0, 0, 0, 1},
				Columns: []*Column{
					&Column{
						Name:  []byte{0, 8, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 1, 'C', 0},
						Value: []byte{'3'},
					},
					&Column{
						Name:  []byte{0, 8, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 1, 'Z', 0},
						Value: []byte{0, 0, 0, 0, 0, 0, 0, 4},
					},
				},
			},
		},

		&structTestShell{
			mapping: NewSparse("cfname", "Key", "FBytes", "FBool", "FInt8", "FInt16", "FInt32", "FInt", "FInt64", "FFloat32", "FFloat64", "FString", "FUUID"),
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