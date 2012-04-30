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

// NewStructMapping returns a mapping valid for all instances of the type of the passed struct
func NewStructMapping(source interface{}) (Mapping, error) {
	v, err := validStruct(source)
	if err != nil {
		return nil, err
	}
	m, err := getStructMapping(v)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *structMapping) Cf() string {
	return m.cf
}

func (m *structMapping) MinColumns() int {
	return len(m.values)
}

type structMapping struct {
	rtype       reflect.Type
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

func internalNewStructMapping(t reflect.Type) (*structMapping, error) {
	m := &structMapping{rtype: t}
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

	// pass 2: build structMapping

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

var structMappingCache map[reflect.Type]*structMapping
var structMappingCacheMutex *sync.Mutex = new(sync.Mutex)

func getStructMapping(v *reflect.Value) (*structMapping, error) {
	var m *structMapping
	var err error
	found := false
	t := v.Type()
	structMappingCacheMutex.Lock()
	if structMappingCache == nil {
		structMappingCache = make(map[reflect.Type]*structMapping)
	}
	if m, found = structMappingCache[t]; !found {
		m, err = internalNewStructMapping(t)
		if err != nil {
			structMappingCache[t] = m
		}
	}
	structMappingCacheMutex.Unlock()
	return m, err
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

func (m *structMapping) validate(source interface{}) (*reflect.Value, error) {
	v, err := validStruct(source)
	if err != nil {
		return nil, err
	}
	t := v.Type()
	if t != m.rtype {
		return nil, errors.New("The passed struct does not have the same type has the mapping")
	}
	return v, nil
}

func (m *structMapping) Map(source interface{}) (*Row, error) {
	v, err := m.validate(source)
	if err != nil {
		return nil, err
	}

	// allocate new row to return
	row := &Row{}

	// marshal the key field
	b, err := m.key.marshalValue(v)
	if err != nil {
		return nil, err
	}
	row.Key = b

	// prepare composite, if needed
	composite := make([]byte, 0)
	for _, f := range m.composite {
		b, err := f.marshalValue(v)
		if err != nil {
			return nil, err
		}
		composite = append(composite, packComposite(b, eocEquals)...)
	}

	// add columns
	for _, f := range m.values {
		columnName, err := f.marshalName()
		if err != nil {
			return nil, err
		}
		if len(composite) > 0 {
			columnName = append(composite, packComposite(columnName, eocEquals)...)
		}
		columnValue, err := f.marshalValue(v)
		if err != nil {
			return nil, err
		}
		row.Columns = append(row.Columns, &Column{Name: columnName, Value: columnValue})
	}

	return row, nil
}

func (m *structMapping) Unmap(destination interface{}, offset int, row *Row) (int, error) {
	readColumns := 0

	v, err := m.validate(destination)
	if err != nil {
		return readColumns, err
	}

	// unmarshal key field
	err = m.key.unmarshalValue(row.Key, v)
	if err != nil {
		return readColumns, err
	}

	compositeFieldsAreSet := false

	// unmarshal col/values
	columns := row.Columns[offset : offset+m.MinColumns()]
	for _, column := range columns {
		readColumns++
		var components [][]byte
		if len(m.composite) > 0 {
			components = unpackComposite(column.Name)
		} else {
			components = [][]byte{column.Name}
		}
		if len(components) != (len(m.composite) + 1) {
			return readColumns, errors.New(fmt.Sprint("Returned number of components in composite column name does not match struct key: composite in struct ", v.Type().Name()))
		}

		// FIXME: it is possible for a row to contain multiple composite values instead of an uniform one. assume
		// that is not the case for now!
		// iterate over composite components, just once, to set the composite fields
		if !compositeFieldsAreSet {
			for i, f := range m.composite {
				b := components[i]
				err := f.unmarshalValue(b, v)
				if err != nil {
					return readColumns, errors.New(fmt.Sprint("Error unmarshaling composite field: ", err))
				}
			}
			compositeFieldsAreSet = true
		}

		// lookup field by name
		var name string
		err = Unmarshal(components[len(components)-1], UTF8Type, &name)
		if err != nil {
			return readColumns, errors.New(fmt.Sprint("Error unmarshaling composite field as UTF8Type for field name in struct ", v.Type().Name(), ", error: ", err))
		}
		if f, found := m.namedValues[name]; found {
			err := f.unmarshalValue(column.Value, v)
			if err != nil {
				return readColumns, errors.New(fmt.Sprint("Error unmarshaling column value: ", err))
			}
		}
	}

	return readColumns, nil
}
