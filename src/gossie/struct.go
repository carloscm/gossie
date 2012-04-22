package gossie

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

/*
todo:

	allow some form of "custom composite" to support common use cases like row keys in the form (read as ascii) 111:222:333.
	-> maybe support it via optional callables over the passed interface? OnMapKey/etc

*/

type mapping struct {
	cf          string
	key         *field
	composite   []*field
	values      []*field
	namedValues map[string]*field
}

type field struct {
	name          string
	index         []int
	cassandraName string
	cassandraType TypeDesc
}

func defaultType(t reflect.Type) TypeDesc {
	switch t.Kind() {
	case reflect.Bool:
		return BooleanType
	case reflect.String:
		return UTF8Type
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return LongType
	case reflect.Float32:
		return FloatType
	case reflect.Float64:
		return DoubleType
	case reflect.Array:
		if t.Name() == "UUID" && t.Size() == 16 {
			return UUIDType
		}
		return UnknownType
	case reflect.Slice:
		if et := t.Elem(); et.Kind() == reflect.Uint8 {
			return BytesType
		}
		return UnknownType
	}
	return UnknownType
}

func newField(sf reflect.StructField) (*field, error) {
	// ignore anon fields
	if sf.Anonymous || sf.Name == "" {
		return nil, nil
	}

	name := sf.Name
	index := sf.Index

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

	return &field{name, index, cassandraName, cassandraType}, nil
}

func (f *field) marshalName() ([]byte, error) {
	b, err := Marshal(f.cassandraName, UTF8Type)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Error marshaling field name for field ", f.name, ":", err))
	}
	return b, nil
}

func (f *field) marshalValue(structValue *reflect.Value) ([]byte, error) {
	v := structValue.FieldByIndex(f.index)
	b, err := Marshal(v.Interface(), f.cassandraType)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Error marshaling filed value for field ", f.name, ":", err))
	}
	return b, nil
}

func (f *field) unmarshalValue(b []byte, structValue *reflect.Value) error {
	v := structValue.FieldByIndex(f.index)
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

func newMapping(t reflect.Type) (*mapping, error) {
	m := &mapping{}
	n := t.NumField()

	cf := ""
	keyAndComposite := ""
	fields := make(map[string]*field, 0)

	// pass 1: gather field metadata
	for i := 0; i < n; i++ {
		sf := t.Field(i)
		if tagValue := sf.Tag.Get("cf"); tagValue != "" {
			cf = tagValue
		}
		if tagValue := sf.Tag.Get("key"); tagValue != "" {
			keyAndComposite = tagValue
		}
		f, err := newField(sf)
		if err != nil {
			return nil, errors.New(fmt.Sprint("Error in struct ", t.Name(), ": ", err))
		}
		if f != nil {
			fields[sf.Name] = f
		}
	}

	// pass 2: build mapping

	if cf == "" {
		return nil, errors.New(fmt.Sprint("No cf field in struct ", t.Name()))
	}
	m.cf = cf

	if keyAndComposite == "" {
		return nil, errors.New(fmt.Sprint("No key field in struct ", t.Name()))
	}
	keyAndCompositeCols := strings.Split(keyAndComposite, ",")
	if len(keyAndCompositeCols) < 1 {
		return nil, errors.New(fmt.Sprint("Not enough key/composite fields given in struct ", t.Name()))
	}
	key, compositeCols := keyAndCompositeCols[0], keyAndCompositeCols[1:]

	if keyField, found := fields[key]; !found {
		return nil, errors.New(fmt.Sprint("Referenced key field ", key, " does not exist in struct ", t.Name()))
	} else {
		m.key = keyField
		delete(fields, key)
	}

	m.composite = make([]*field, 0)
	for _, name := range compositeCols {
		f, found := fields[name]
		if !found {
			return nil, errors.New(fmt.Sprint("Referenced composite field ", name, " does not exist in struct ", t.Name()))
		}
		m.composite = append(m.composite, f)
		delete(fields, name)
	}

	m.values = make([]*field, 0)
	m.namedValues = make(map[string]*field, 0)
	for i := 0; i < n; i++ {
		sf := t.Field(i)
		if f, found := fields[sf.Name]; found {
			m.values = append(m.values, f)
			m.namedValues[f.cassandraName] = f
		}
	}

	return m, nil
}

var mapCache map[reflect.Type]*mapping
var mapCacheMutex *sync.Mutex = new(sync.Mutex)

func getMapping(v reflect.Value) (*mapping, error) {
	var m *mapping
	var err error
	found := false
	t := v.Type()
	mapCacheMutex.Lock()
	if mapCache == nil {
		mapCache = make(map[reflect.Type]*mapping)
	}
	if m, found = mapCache[t]; !found {
		m, err = newMapping(t)
		if err != nil {
			mapCache[t] = m
		}
	}
	mapCacheMutex.Unlock()
	return m, err
}

// mappedInstance stores the reflect.Value and mapping for a particular instance of a struct
type mappedInstance struct {
	source interface{}
	v      reflect.Value
	m      *mapping
}

func newMappedInstance(source interface{}) (*mappedInstance, error) {
	mi := &mappedInstance{source: source}
	var err error

	// always work with a pointer to struct
	vp := reflect.ValueOf(source)
	if vp.Kind() != reflect.Ptr {
		return nil, errors.New("Passed source is not a pointer to a struct")
	}
	if vp.IsNil() {
		return nil, errors.New("Passed source is not a pointer to a struct")
	}
	mi.v = reflect.Indirect(vp)
	if mi.v.Kind() != reflect.Struct {
		return nil, errors.New("Passed source is not a pointer to a struct")
	}
	if !mi.v.CanSet() {
		return nil, errors.New("Cannot modify the passed struct instance")
	}

	mi.m, err = getMapping(mi.v)
	if err != nil {
		return nil, err
	}

	return mi, nil
}

func internalMap(source interface{}) (*Row, *mappedInstance, error) {

	// deconstruct the source struct into a reflect.Value and a (cached) struct mapping
	mi, err := newMappedInstance(source)
	if err != nil {
		return nil, nil, err
	}

	// allocate new row to return
	row := &Row{}

	// marshal the key field
	b, err := mi.m.key.marshalValue(&mi.v)
	if err != nil {
		return nil, nil, err
	}
	row.Key = b

	// prepare composite, if needed
	composite := make([]byte, 0)
	for _, f := range mi.m.composite {
		b, err := f.marshalValue(&mi.v)
		if err != nil {
			return nil, nil, err
		}
		composite = append(composite, packComposite(b, eocEquals)...)
	}

	// add columns
	for _, f := range mi.m.values {
		columnName, err := f.marshalName()
		if err != nil {
			return nil, nil, err
		}
		if len(composite) > 0 {
			columnName = append(composite, packComposite(columnName, eocEquals)...)
		}
		columnValue, err := f.marshalValue(&mi.v)
		if err != nil {
			return nil, nil, err
		}
		row.Columns = append(row.Columns, &Column{Name: columnName, Value: columnValue})
	}

	return row, mi, nil
}

func Map(source interface{}) (*Row, error) {
	row, _, err := internalMap(source)
	return row, err
}

func Unmap(row *Row, destination interface{}) error {

	// deconstruct the source struct into a reflect.Value and a (cached) struct mapping
	mi, err := newMappedInstance(destination)
	if err != nil {
		return err
	}

	// unmarshal key field
	err = mi.m.key.unmarshalValue(row.Key, &mi.v)
	if err != nil {
		return err
	}

	compositeFieldsAreSet := false

	// unmarshal col/values
	for _, column := range row.Columns {
		var components [][]byte
		if len(mi.m.composite) > 0 {
			components = unpackComposite(column.Name)
		} else {
			components = [][]byte{column.Name}
		}
		if len(components) != (len(mi.m.composite) + 1) {
			return errors.New(fmt.Sprint("Returned number of components in composite column name does not match struct key: composite in struct ", mi.v.Type().Name()))
		}

		// FIXME: it is possible for a row to contain multiple composite values instead of an uniform one. assume
		// that is not the case for now!
		// iterate over composite components, just once, to set the composite fields
		if !compositeFieldsAreSet {
			for i, f := range mi.m.composite {
				b := components[i]
				err := f.unmarshalValue(b, &mi.v)
				if err != nil {
					return errors.New(fmt.Sprint("Error unmarshaling composite field: ", err))
				}
			}
			compositeFieldsAreSet = true
		}

		// lookup field by name
		var name string
		err = Unmarshal(components[len(components)-1], UTF8Type, &name)
		if err != nil {
			return errors.New(fmt.Sprint("Error unmarshaling composite field as UTF8Type for field name in struct ", mi.v.Type().Name(), ", error: ", err))
		}
		if f, found := mi.m.namedValues[name]; found {
			err := f.unmarshalValue(column.Value, &mi.v)
			if err != nil {
				return errors.New(fmt.Sprint("Error unmarshaling column value: ", err))
			}
		}
	}

	return nil
}
