package gossie

import (
	"reflect"
	"testing"
)

/*

todo:

    much more testing on struct mapping

    more test for name: and type:

*/

type errNoMeta struct {
	A int
}
type errNoMetaKeyColVal struct {
	A int `cf:"cfname"`
}
type errNoMetaColVal struct {
	A int `cf:"cfname" key:"A"`
}
type errNoMetaVal struct {
	A int `cf:"cfname" key:"A" col:"B"`
	B int
}
type errInvKey struct {
	A int `cf:"cfname" key:"Z" col:"B" val:"C"`
	B int
	C int
}
type errInvCol struct {
	A int `cf:"cfname" key:"A" col:"Z" val:"C"`
	B int
	C int
}
type errInvVal struct {
	A int `cf:"cfname" key:"A" col:"B" val:"Z"`
	B int
	C int
}
type errStarNameNotLast struct {
	A int `cf:"cfname" key:"A" col:"*name,B" val:"*value"`
	B int
	C int
}
type errSliceNotLast struct {
	A int `cf:"cfname" key:"A" col:"C,B" val:"D"`
	B int
	C []int
	D []int
}
type noErrA struct {
	A int `cf:"cfname" key:"A" col:"B" val:"C"`
	B int
	C int
}
type noErrB struct {
	A int `cf:"cfname" key:"A" col:"*name" val:"*value"`
	B int `name:"Z"`
	C int `type:"AsciiType"`
}
type noErrC struct {
	A int `cf:"cfname" key:"A" col:"B,*name" val:"*value"`
	B int
	C int
}
type noErrD struct {
	A int `cf:"cfname" key:"A" col:"B" val:"C"`
	B []int
	C []int
}
type noErrE struct {
	A int `cf:"cfname" key:"A" col:"B,C" val:"D"`
	B int
	C []int
	D []int
}
type everythingComp struct {
	Key      string `cf:"cfname" key:"Key" col:"FBytes,FBool,FInt8,FInt16,FInt32,FInt,FInt64,FFloat32,FFloat64,FString,FUUID,*name" val:"*value"`
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
	return newStructMapping(typ)
}

func structMapMustError(t *testing.T, instance interface{}) {
	_, err := buildMappingFromPtr(instance)
	if err == nil {
		t.Error("Expected error calling newStructMapping, got none")
	}
}

func checkMapping(t *testing.T, expected, actual interface{}, name string) {
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Mapping for struct sample", name, "does not match expected output")
	}
}

func TestStructMapping(t *testing.T) {
	structMapMustError(t, &errNoMeta{})
	structMapMustError(t, &errNoMetaKeyColVal{})
	structMapMustError(t, &errNoMetaColVal{})
	structMapMustError(t, &errNoMetaVal{})
	structMapMustError(t, &errInvKey{})
	structMapMustError(t, &errInvCol{})
	structMapMustError(t, &errInvVal{})
	structMapMustError(t, &errStarNameNotLast{})
	structMapMustError(t, &errSliceNotLast{})

	mapA, _ := buildMappingFromPtr(&noErrA{1, 2, 3})
	goodA := &structMapping{
		cf:  "cfname",
		key: &fieldMapping{fieldKind: baseTypeField, position: 0, name: "A", cassandraType: LongType, cassandraName: "A"},
		columns: []*fieldMapping{
			&fieldMapping{fieldKind: baseTypeField, position: 1, name: "B", cassandraType: LongType, cassandraName: "B"},
		},
		value:             &fieldMapping{fieldKind: baseTypeField, position: 2, name: "C", cassandraType: LongType, cassandraName: "C"},
		others:            make(map[string]*fieldMapping, 0),
		isCompositeColumn: false,
	}
	checkMapping(t, goodA, mapA, "mapA")

	mapB, _ := buildMappingFromPtr(&noErrB{1, 2, 3})
	goodB := &structMapping{
		cf:  "cfname",
		key: &fieldMapping{fieldKind: baseTypeField, position: 0, name: "A", cassandraType: LongType, cassandraName: "A"},
		columns: []*fieldMapping{
			&fieldMapping{fieldKind: starNameField, position: 0, name: "", cassandraType: 0, cassandraName: ""},
		},
		value: &fieldMapping{fieldKind: starValueField, position: 0, name: "", cassandraType: 0, cassandraName: ""},
		others: map[string]*fieldMapping{
			"Z": &fieldMapping{fieldKind: baseTypeField, position: 1, name: "B", cassandraType: LongType, cassandraName: "Z"},
			"C": &fieldMapping{fieldKind: baseTypeField, position: 2, name: "C", cassandraType: AsciiType, cassandraName: "C"},
		},
		isCompositeColumn: false,
	}
	checkMapping(t, goodB, mapB, "mapB")

	mapC, _ := buildMappingFromPtr(&noErrC{1, 2, 3})
	goodC := &structMapping{
		cf:  "cfname",
		key: &fieldMapping{fieldKind: baseTypeField, position: 0, name: "A", cassandraType: LongType, cassandraName: "A"},
		columns: []*fieldMapping{
			&fieldMapping{fieldKind: baseTypeField, position: 1, name: "B", cassandraType: LongType, cassandraName: "B"},
			&fieldMapping{fieldKind: starNameField, position: 0, name: "", cassandraType: 0, cassandraName: ""},
		},
		value: &fieldMapping{fieldKind: starValueField, position: 0, name: "", cassandraType: 0, cassandraName: ""},
		others: map[string]*fieldMapping{
			"C": &fieldMapping{fieldKind: baseTypeField, position: 2, name: "C", cassandraType: LongType, cassandraName: "C"},
		},
		isCompositeColumn: true,
	}
	checkMapping(t, goodC, mapC, "mapC")

	mapD, _ := buildMappingFromPtr(&noErrD{1, []int{2, 3}, []int{4, 5}})
	goodD := &structMapping{
		cf:  "cfname",
		key: &fieldMapping{fieldKind: baseTypeField, position: 0, name: "A", cassandraType: LongType, cassandraName: "A"},
		columns: []*fieldMapping{
			&fieldMapping{fieldKind: baseTypeSliceField, position: 1, name: "B", cassandraType: LongType, cassandraName: "B"},
		},
		value:             &fieldMapping{fieldKind: baseTypeSliceField, position: 2, name: "C", cassandraType: LongType, cassandraName: "C"},
		others:            make(map[string]*fieldMapping, 0),
		isCompositeColumn: false,
	}
	checkMapping(t, goodD, mapD, "mapD")

	mapE, _ := buildMappingFromPtr(&noErrE{1, 2, []int{3, 4}, []int{5, 6}})
	goodE := &structMapping{
		cf:  "cfname",
		key: &fieldMapping{fieldKind: baseTypeField, position: 0, name: "A", cassandraType: LongType, cassandraName: "A"},
		columns: []*fieldMapping{
			&fieldMapping{fieldKind: baseTypeField, position: 1, name: "B", cassandraType: LongType, cassandraName: "B"},
			&fieldMapping{fieldKind: baseTypeSliceField, position: 2, name: "C", cassandraType: LongType, cassandraName: "C"},
		},
		value:             &fieldMapping{fieldKind: baseTypeSliceField, position: 3, name: "D", cassandraType: LongType, cassandraName: "D"},
		others:            make(map[string]*fieldMapping, 0),
		isCompositeColumn: true,
	}
	checkMapping(t, goodE, mapE, "mapE")

	eComp, _ := buildMappingFromPtr(&everythingComp{"A", []byte{1, 2}, true, 3, 4, 5, 6, 7, 8.0, 9.0, "B",
		[16]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}, "C"})
	goodEComp := &structMapping{
		cf:  "cfname",
		key: &fieldMapping{fieldKind: baseTypeField, position: 0, name: "Key", cassandraType: UTF8Type, cassandraName: "Key"},
		columns: []*fieldMapping{
			&fieldMapping{fieldKind: baseTypeField, position: 1, name: "FBytes", cassandraType: BytesType, cassandraName: "FBytes"},
			&fieldMapping{fieldKind: baseTypeField, position: 2, name: "FBool", cassandraType: BooleanType, cassandraName: "FBool"},
			&fieldMapping{fieldKind: baseTypeField, position: 3, name: "FInt8", cassandraType: LongType, cassandraName: "FInt8"},
			&fieldMapping{fieldKind: baseTypeField, position: 4, name: "FInt16", cassandraType: LongType, cassandraName: "FInt16"},
			&fieldMapping{fieldKind: baseTypeField, position: 5, name: "FInt32", cassandraType: LongType, cassandraName: "FInt32"},
			&fieldMapping{fieldKind: baseTypeField, position: 6, name: "FInt", cassandraType: LongType, cassandraName: "FInt"},
			&fieldMapping{fieldKind: baseTypeField, position: 7, name: "FInt64", cassandraType: LongType, cassandraName: "FInt64"},
			&fieldMapping{fieldKind: baseTypeField, position: 8, name: "FFloat32", cassandraType: FloatType, cassandraName: "FFloat32"},
			&fieldMapping{fieldKind: baseTypeField, position: 9, name: "FFloat64", cassandraType: DoubleType, cassandraName: "FFloat64"},
			&fieldMapping{fieldKind: baseTypeField, position: 10, name: "FString", cassandraType: UTF8Type, cassandraName: "FString"},
			&fieldMapping{fieldKind: baseTypeField, position: 11, name: "FUUID", cassandraType: UUIDType, cassandraName: "FUUID"},
			&fieldMapping{fieldKind: starNameField, position: 0, name: "", cassandraType: 0, cassandraName: ""},
		},
		value: &fieldMapping{fieldKind: starValueField, position: 0, name: "", cassandraType: 0, cassandraName: ""},
		others: map[string]*fieldMapping{
			"Val": &fieldMapping{fieldKind: baseTypeField, position: 12, name: "Val", cassandraType: UTF8Type, cassandraName: "Val"},
		},
		isCompositeColumn: true,
	}
	checkMapping(t, goodEComp, eComp, "eComp")

}

type structTestShell struct {
	expectedStruct interface{}
	resultStruct   interface{}
	expectedRow    *Row
	name           string
}

func (shell *structTestShell) checkMap(t *testing.T, expectedStruct interface{}) {
	resultRow, err := Map(expectedStruct)
	if err != nil {
		t.Error("Error mapping struct: ", err)
	}
	if !reflect.DeepEqual(resultRow, shell.expectedRow) {
		t.Error("Mapped struct ", shell.name, " does not match expected row ", shell.expectedRow, " actual ", resultRow)
	}
}

func (shell *structTestShell) checkUnmap(t *testing.T) interface{} {
	err := Unmap(shell.expectedRow, shell.resultStruct)
	if err != nil {
		t.Error("Error umapping struct: ", err)
	}
	if !reflect.DeepEqual(shell.resultStruct, shell.expectedStruct) {
		t.Error("Unmapped struct ", shell.name, " does not match expected instance ", shell.expectedStruct, " actual ", shell.resultStruct)
	}
	return shell.resultStruct
}

func (shell *structTestShell) checkFullMap(t *testing.T) {
	shell.checkMap(t, shell.expectedStruct)
	intermediateStruct := shell.checkUnmap(t)
	shell.checkMap(t, intermediateStruct)
}

func TestMap(t *testing.T) {

	shells := []*structTestShell{
		&structTestShell{
			name:           "noErrA",
			expectedStruct: &noErrA{1, 2, 3},
			resultStruct:   &noErrA{},
			expectedRow: &Row{
				Key: []byte{0, 0, 0, 0, 0, 0, 0, 1},
				Columns: []*Column{
					&Column{
						Name:  []byte{0, 0, 0, 0, 0, 0, 0, 2},
						Value: []byte{0, 0, 0, 0, 0, 0, 0, 3},
					},
				},
			},
		},

		&structTestShell{
			name:           "noErrB",
			expectedStruct: &noErrB{1, 2, 3},
			resultStruct:   &noErrB{},
			expectedRow: &Row{
				Key: []byte{0, 0, 0, 0, 0, 0, 0, 1},
				Columns: []*Column{
					&Column{
						Name:  []byte{67},
						Value: []byte{51},
					},
					&Column{
						Name:  []byte{90},
						Value: []byte{0, 0, 0, 0, 0, 0, 0, 2},
					},
				},
			},
		},

		&structTestShell{
			name:           "noErrE",
			expectedStruct: &noErrE{1, 2, []int{5, 6}, []int{7, 8}},
			resultStruct:   &noErrE{},
			expectedRow: &Row{
				Key: []byte{0, 0, 0, 0, 0, 0, 0, 1},
				Columns: []*Column{
					&Column{
						Name:  []byte{0, 8, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 5, 0},
						Value: []byte{0, 0, 0, 0, 0, 0, 0, 7},
					},
					&Column{
						Name:  []byte{0, 8, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 6, 0},
						Value: []byte{0, 0, 0, 0, 0, 0, 0, 8},
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
