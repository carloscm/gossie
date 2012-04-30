package gossie

import (
	"reflect"
	"testing"
)

/*

todo:

    more testing on struct mapping, more cases, more automation

    more tests for name: and type:

*/

type errNoMeta struct {
	A int
}
type errNoMetaKey struct {
	A int `cf:"cfname"`
}
type errInvKey struct {
	A int `cf:"cfname" key:"Z"`
	B int
	C int
}
type noErrA struct {
	A int `cf:"cfname" key:"A"`
	B int
	C int
}
type noErrB struct {
	A int `cf:"cfname" key:"A,B"`
	B int
	C int `type:"AsciiType"`
	D int `name:"Z"`
}
type everythingComp struct {
	Key      string `cf:"cfname" key:"Key,FBytes,FBool,FInt8,FInt16,FInt32,FInt,FInt64,FFloat32,FFloat64,FString,FUUID"`
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

func buildMappingFromPtr(instance interface{}) (*structMapping, error) {
	valuePtr := reflect.ValueOf(instance)
	value := reflect.Indirect(valuePtr)
	typ := value.Type()
	return internalNewStructMapping(typ)
}

func structMapMustError(t *testing.T, instance interface{}) {
	_, err := buildMappingFromPtr(instance)
	if err == nil {
		t.Error("Expected error calling newMapping, got none")
	}
}

func checkMapping(t *testing.T, expected, actual interface{}, name string) {
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Mapping for struct sample", name, "does not match expected output")
	}
}

func TestStructMapping(t *testing.T) {
	structMapMustError(t, &errNoMeta{})
	structMapMustError(t, &errNoMetaKey{})
	structMapMustError(t, &errInvKey{})

	mapA, err := buildMappingFromPtr(&noErrA{1, 2, 3})
	valuePtr := reflect.ValueOf(&noErrA{})
	value := reflect.Indirect(valuePtr)
	typ := value.Type()
	goodA := &structMapping{
		rtype:     typ,
		cf:        "cfname",
		key:       &field{index: []int{0}, name: "A", cassandraType: LongType, cassandraName: "A"},
		composite: []*field{},
		values: []*field{
			&field{index: []int{1}, name: "B", cassandraType: LongType, cassandraName: "B"},
			&field{index: []int{2}, name: "C", cassandraType: LongType, cassandraName: "C"},
		},
		namedValues: map[string]*field{
			"B": &field{index: []int{1}, name: "B", cassandraType: LongType, cassandraName: "B"},
			"C": &field{index: []int{2}, name: "C", cassandraType: LongType, cassandraName: "C"},
		},
	}
	if err != nil {
		t.Fatal("Unexpected error calling mapA newMapping:", err)
	}
	checkMapping(t, goodA, mapA, "mapA")

	mapB, err := buildMappingFromPtr(&noErrB{1, 2, 3, 4})
	valuePtr = reflect.ValueOf(&noErrB{})
	value = reflect.Indirect(valuePtr)
	typ = value.Type()
	goodB := &structMapping{
		rtype: typ,
		cf:    "cfname",
		key:   &field{index: []int{0}, name: "A", cassandraType: LongType, cassandraName: "A"},
		composite: []*field{
			&field{index: []int{1}, name: "B", cassandraType: LongType, cassandraName: "B"},
		},
		values: []*field{
			&field{index: []int{2}, name: "C", cassandraType: AsciiType, cassandraName: "C"},
			&field{index: []int{3}, name: "D", cassandraType: LongType, cassandraName: "Z"},
		},
		namedValues: map[string]*field{
			"C": &field{index: []int{2}, name: "C", cassandraType: AsciiType, cassandraName: "C"},
			"Z": &field{index: []int{3}, name: "D", cassandraType: LongType, cassandraName: "Z"},
		},
	}
	if err != nil {
		t.Fatal("Unexpected error calling mapB newMapping:", err)
	}
	checkMapping(t, goodB, mapB, "mapB")

	eComp, err := buildMappingFromPtr(&everythingComp{"A", []byte{1, 2}, true, 3, 4, 5, 6, 7, 8.0, 9.0, "B",
		[16]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}, "C"})
	valuePtr = reflect.ValueOf(&everythingComp{})
	value = reflect.Indirect(valuePtr)
	typ = value.Type()
	goodEComp := &structMapping{
		rtype: typ,
		cf:    "cfname",
		key:   &field{index: []int{0}, name: "Key", cassandraType: UTF8Type, cassandraName: "Key"},
		composite: []*field{
			&field{index: []int{1}, name: "FBytes", cassandraType: BytesType, cassandraName: "FBytes"},
			&field{index: []int{2}, name: "FBool", cassandraType: BooleanType, cassandraName: "FBool"},
			&field{index: []int{3}, name: "FInt8", cassandraType: LongType, cassandraName: "FInt8"},
			&field{index: []int{4}, name: "FInt16", cassandraType: LongType, cassandraName: "FInt16"},
			&field{index: []int{5}, name: "FInt32", cassandraType: LongType, cassandraName: "FInt32"},
			&field{index: []int{6}, name: "FInt", cassandraType: LongType, cassandraName: "FInt"},
			&field{index: []int{7}, name: "FInt64", cassandraType: LongType, cassandraName: "FInt64"},
			&field{index: []int{8}, name: "FFloat32", cassandraType: FloatType, cassandraName: "FFloat32"},
			&field{index: []int{9}, name: "FFloat64", cassandraType: DoubleType, cassandraName: "FFloat64"},
			&field{index: []int{10}, name: "FString", cassandraType: UTF8Type, cassandraName: "FString"},
			&field{index: []int{11}, name: "FUUID", cassandraType: UUIDType, cassandraName: "FUUID"},
		},
		values: []*field{
			&field{index: []int{12}, name: "Val", cassandraType: UTF8Type, cassandraName: "Val"},
		},
		namedValues: map[string]*field{
			"Val": &field{index: []int{12}, name: "Val", cassandraType: UTF8Type, cassandraName: "Val"},
		},
	}
	if err != nil {
		t.Fatal("Unexpected error calling eComp newMapping:", err)
	}
	checkMapping(t, goodEComp, eComp, "eComp")

}

type structTestShell struct {
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

func (shell *structTestShell) checkUnmap(t *testing.T, m Mapping) interface{} {
	n, err := m.Unmap(shell.resultStruct, 0, shell.expectedRow)
	if err != nil {
		t.Error("Error umapping struct: ", err)
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
	m, err := NewStructMapping(shell.expectedStruct)
	if err != nil {
		t.Error("Unexpected error creating a struct mapping: ", err)
	}
	shell.checkMap(t, m, shell.expectedStruct, 1)
	intermediateStruct := shell.checkUnmap(t, m)
	shell.checkMap(t, m, intermediateStruct, 2)
}

func TestMap(t *testing.T) {

	m, _ := NewStructMapping(&noErrA{})
	_, err := m.Map(&noErrB{})
	if err == nil {
		t.Error("Error expected when trying to use a struct Mapping.Map with an instance of a different struct type")
	}

	shells := []*structTestShell{
		&structTestShell{
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
			name: "everythingComp",
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
