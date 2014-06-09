package gossie

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type structInspection struct {
	rtype           reflect.Type
	orderedFields   []*field
	goFields        map[string]*field
	cassandraFields map[string]*field
	globalTags      map[string]string
}

type field struct {
	name          string
	index         int
	cassandraName string
	cassandraType TypeDesc
	skipEmpty     bool
}

var recognizedGlobalTags []string = []string{"mapping", "cf", "key", "cols", "value"}

func newField(index int, sf reflect.StructField) (*field, error) {
	// ignore anon fields
	if sf.Anonymous || sf.Name == "" {
		return nil, nil
	}

	if tagType := sf.Tag.Get("skip"); tagType == "true" {
		return nil, nil
	}

	name := sf.Name

	cassandraType := defaultType(sf.Type)
	if cassandraType == UnknownType {
		return nil, errors.New(fmt.Sprint("Field ", name, " has unsupported type"))
	}
	if tagType := sf.Tag.Get("type"); tagType != "" {
		cassandraType = parseTypeDesc(tagType)
	}

	cassandraName := name
	if tagName := sf.Tag.Get("name"); tagName != "" {
		cassandraName = tagName
	}

	skipEmpty := false
	if tagSkipEmpty := sf.Tag.Get("skipempty"); tagSkipEmpty == "true" {
		skipEmpty = true
	}

	return &field{name, index, cassandraName, cassandraType, skipEmpty}, nil
}

func (f *field) marshalName() ([]byte, error) {
	b, err := Marshal(f.cassandraName, UTF8Type)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Error marshaling field name for field ", f.name, ":", err))
	}
	return b, nil
}

func (f *field) marshalValue(structValue *reflect.Value) ([]byte, error) {
	v := structValue.Field(f.index)
	b, err := Marshal(v.Interface(), f.cassandraType)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Error marshaling field value for field ", f.name, ":", err))
	}
	return b, nil
}

func (f *field) isEmpty(structValue *reflect.Value) bool {
	v := structValue.Field(f.index)
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}

func (f *field) unmarshalValue(b []byte, structValue *reflect.Value) error {
	v := structValue.Field(f.index)
	if !v.CanAddr() {
		return errors.New(fmt.Sprint("Cannot obtain pointer to field ", f.name))
	}
	vp := v.Addr()
	err := Unmarshal(b, f.cassandraType, vp.Interface())
	if err != nil {
		return errors.New(fmt.Sprint("Error unmarshaling field ", f.name, ":", err))
	}
	return nil
}

func newStructInspection(t reflect.Type) (*structInspection, error) {
	si := &structInspection{
		rtype:           t,
		orderedFields:   make([]*field, 0),
		goFields:        make(map[string]*field, 0),
		cassandraFields: make(map[string]*field, 0),
		globalTags:      make(map[string]string),
	}
	n := t.NumField()
	for i := 0; i < n; i++ {
		sf := t.Field(i)
		for _, t := range recognizedGlobalTags {
			if v := sf.Tag.Get(t); v != "" {
				si.globalTags[t] = v
			}
		}
		f, err := newField(i, sf)
		if err != nil {
			return nil, errors.New(fmt.Sprint("Error in struct ", t.Name(), ": ", err))
		}
		if f != nil {
			si.orderedFields = append(si.orderedFields, f)
			si.goFields[f.name] = f
			si.cassandraFields[f.cassandraName] = f
		}
	}
	return si, nil
}

var structInspectionCache map[reflect.Type]*structInspection
var structInspectionCacheMutex *sync.Mutex = new(sync.Mutex)

func inspectStruct(v *reflect.Value) (*structInspection, error) {
	var si *structInspection
	var err error
	found := false
	t := v.Type()
	structInspectionCacheMutex.Lock()
	if structInspectionCache == nil {
		structInspectionCache = make(map[reflect.Type]*structInspection)
	}
	if si, found = structInspectionCache[t]; !found {
		si, err = newStructInspection(t)
		if err != nil {
			structInspectionCache[t] = si
		}
	}
	structInspectionCacheMutex.Unlock()
	return si, err
}

func validStruct(source interface{}) (*reflect.Value, error) {
	// always work with a pointer to struct
	vp := reflect.ValueOf(source)
	if vp.Kind() != reflect.Ptr {
		return nil, errors.New("Passed source is not a pointer to a struct")
	}
	if vp.IsNil() {
		return nil, errors.New("Passed source is not a pointer to a struct")
	}
	v := reflect.Indirect(vp)
	if v.Kind() != reflect.Struct {
		return nil, errors.New("Passed source is not a pointer to a struct")
	}
	if !v.CanSet() {
		return nil, errors.New("Cannot modify the passed struct instance")
	}
	return &v, nil
}

func validateAndInspectStruct(source interface{}) (*reflect.Value, *structInspection, error) {
	v, err := validStruct(source)
	if err != nil {
		return nil, nil, err
	}
	si, err := inspectStruct(v)
	if err != nil {
		return nil, nil, err
	}
	return v, si, nil
}
