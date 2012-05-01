package gossie

import (
	"reflect"
	"testing"
)

type noErrA struct {
	A int
	B int
	C int
}
type noErrB struct {
	A int
	B int
	C int `type:"AsciiType"`
	D int `name:"Z"`
}

func buildInspectionFromPtr(instance interface{}) (*structInspection, error) {
	valuePtr := reflect.ValueOf(instance)
	value := reflect.Indirect(valuePtr)
	return inspectStruct(&value)
}

func structMapMustError(t *testing.T, instance interface{}) {
	_, err := buildInspectionFromPtr(instance)
	if err == nil {
		t.Error("Expected error calling newInspection, got none")
	}
}

func checkInspection(t *testing.T, expected, actual interface{}, name string) {
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Inspection for struct sample", name, "does not match expected output")
	}
}

func TestStructInspection(t *testing.T) {

	mapA, err := buildInspectionFromPtr(&noErrA{1, 2, 3})
	valuePtr := reflect.ValueOf(&noErrA{})
	value := reflect.Indirect(valuePtr)
	typ := value.Type()
	goodA := &structInspection{
		rtype: typ,
		orderedFields: []*field{
			&field{index: 0, name: "A", cassandraType: LongType, cassandraName: "A"},
			&field{index: 1, name: "B", cassandraType: LongType, cassandraName: "B"},
			&field{index: 2, name: "C", cassandraType: LongType, cassandraName: "C"},
		},
		goFields: map[string]*field{
			"A": &field{index: 0, name: "A", cassandraType: LongType, cassandraName: "A"},
			"B": &field{index: 1, name: "B", cassandraType: LongType, cassandraName: "B"},
			"C": &field{index: 2, name: "C", cassandraType: LongType, cassandraName: "C"},
		},
		cassandraFields: map[string]*field{
			"A": &field{index: 0, name: "A", cassandraType: LongType, cassandraName: "A"},
			"B": &field{index: 1, name: "B", cassandraType: LongType, cassandraName: "B"},
			"C": &field{index: 2, name: "C", cassandraType: LongType, cassandraName: "C"},
		},
	}
	if err != nil {
		t.Fatal("Unexpected error calling mapA newInspection:", err)
	}
	checkInspection(t, goodA, mapA, "mapA")

	mapB, err := buildInspectionFromPtr(&noErrB{1, 2, 3, 4})
	valuePtr = reflect.ValueOf(&noErrB{})
	value = reflect.Indirect(valuePtr)
	typ = value.Type()
	goodB := &structInspection{
		rtype: typ,
		orderedFields: []*field{
			&field{index: 0, name: "A", cassandraType: LongType, cassandraName: "A"},
			&field{index: 1, name: "B", cassandraType: LongType, cassandraName: "B"},
			&field{index: 2, name: "C", cassandraType: AsciiType, cassandraName: "C"},
			&field{index: 3, name: "D", cassandraType: LongType, cassandraName: "Z"},
		},
		goFields: map[string]*field{
			"A": &field{index: 0, name: "A", cassandraType: LongType, cassandraName: "A"},
			"B": &field{index: 1, name: "B", cassandraType: LongType, cassandraName: "B"},
			"C": &field{index: 2, name: "C", cassandraType: AsciiType, cassandraName: "C"},
			"D": &field{index: 3, name: "D", cassandraType: LongType, cassandraName: "Z"},
		},
		cassandraFields: map[string]*field{
			"A": &field{index: 0, name: "A", cassandraType: LongType, cassandraName: "A"},
			"B": &field{index: 1, name: "B", cassandraType: LongType, cassandraName: "B"},
			"C": &field{index: 2, name: "C", cassandraType: AsciiType, cassandraName: "C"},
			"Z": &field{index: 3, name: "D", cassandraType: LongType, cassandraName: "Z"},
		},
	}
	if err != nil {
		t.Fatal("Unexpected error calling mapB newInspection:", err)
	}
	checkInspection(t, goodB, mapB, "mapB")
}
