package gossie

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
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
	name           string
	index          int
	cassandraName  string
	gossieType     GossieType
	gossieTypeArgs *string
	cassandraType  TypeDesc
	skipEmpty      bool
}

var recognizedGlobalTags []string = []string{"mapping", "cf", "key", "cols", "value", "marshal"}

type GossieType interface {
	// Marshaler should wrap the struct value in a gossie.Marshaler
	// tagArgs are the arguments provided in the tag, such as `marshal:"json,{}"`
	// the `{}` would be the tagArgs
	Marshaler(value interface{}, tagArgs *string) Marshaler
	// Unmarshaler should wrap the struct value in a gossie.Unmarshaler
	// tagArgs are the arguments provided in the tag, such as `marshal:"json,{}"`
	// the `{}` would be the tagArgs
	Unmarshaler(value interface{}, tagArgs *string) Unmarshaler
}

// Allows you to specify `marshal:"json"` for example to use the jsonType GossieType
// This makes it easier to use common custom encodings
var gossieTypes = map[string]GossieType{
	"json":       &jsonType{},
	"boolstring": &boolStringType{},
}

// Register a custom GossieType to be used with the given "marshal" struct tag
func RegisterGossieType(name string, gossieType GossieType) {
	gossieTypes[name] = gossieType
}

func newField(index int, sf reflect.StructField) (*field, error) {
	// ignore anon fields
	if sf.Anonymous || sf.Name == "" {
		return nil, nil
	}

	if tagType := sf.Tag.Get("skip"); tagType == "true" {
		return nil, nil
	}

	name := sf.Name

	// Check if a specific GossieType has been requested
	var gossieType GossieType
	var gossieTypeArgs *string
	if tagMarshal := sf.Tag.Get("marshal"); tagMarshal != "" {
		parts := strings.SplitN(tagMarshal, ",", 2)
		gossieType = gossieTypes[parts[0]]
		if len(parts) > 1 {
			gossieTypeArgs = &parts[1]
		}
		if gossieType == nil {
			return nil, fmt.Errorf("Unregistered marshal type: %v", tagMarshal)
		}
	}

	var cassandraType TypeDesc
	if tagType := sf.Tag.Get("type"); tagType != "" {
		cassandraType = parseTypeDesc(tagType)
	} else if gossieType != nil {
		cassandraType = BytesType
	} else {
		cassandraType = defaultType(sf.Type)
	}

	if cassandraType == UnknownType {
		return nil, errors.New(fmt.Sprint("Field ", name, " has unsupported type"))
	}

	cassandraName := name
	if tagName := sf.Tag.Get("name"); tagName != "" {
		cassandraName = tagName
	}

	skipEmpty := false
	if tagSkipEmpty := sf.Tag.Get("skipempty"); tagSkipEmpty == "true" {
		skipEmpty = true
	}

	return &field{name, index, cassandraName, gossieType, gossieTypeArgs, cassandraType, skipEmpty}, nil
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
	vi := v.Interface()
	if f.gossieType != nil {
		vi = f.gossieType.Marshaler(vi, f.gossieTypeArgs)
	}
	b, err := Marshal(vi, f.cassandraType)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Error marshaling field value for field ", f.name, ":", err))
	}
	return b, nil
}

func (f *field) isEmpty(structValue *reflect.Value) bool {
	v := structValue.Field(f.index)
	switch v.Kind() {
	case reflect.Slice: // for []byte
		return v.IsNil() || v.Len() == 0
	default:
		return v.Interface() == reflect.Zero(v.Type()).Interface()
	}
}

func (f *field) unmarshalValue(b []byte, structValue *reflect.Value) error {
	v := structValue.Field(f.index)
	if !v.CanAddr() {
		return errors.New(fmt.Sprint("Cannot obtain pointer to field ", f.name))
	}
	vp := v.Addr()
	vpi := vp.Interface()
	if f.gossieType != nil {
		vpi = f.gossieType.Unmarshaler(vpi, f.gossieTypeArgs)
	}
	err := Unmarshal(b, f.cassandraType, vpi)
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
