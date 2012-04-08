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

	allow to reuse the same fields as both key and col component

	OR: introduce special key field tags to use during sharding, since the "row key" concept losses meaning in this case,
	since it is managed by the lib

    go maps support for things like
    type s struct {
        a    int `cf:"cfname" key:"a" col:"atts" val:"atts"`
        atts map[string]string
    }
    type s2 struct {
        a    int `cf:"cfname" key:"a" col:"b,atts" val:"atts"`
        b    UUID
        atts map[string]string
    }
    --> then think about slicing/pagging this, oops

    support composite key and composite values, not just composite column names (are those actually in use by anybody???)
*/

const (
	_ = iota
	baseTypeField
	baseTypeSliceField
	starNameField
	starValueField
)

// mapping stores how to map from/to a struct
type fieldMapping struct {
	fieldKind     int
	position      int
	name          string
	cassandraName string
	cassandraType TypeDesc
}
type structMapping struct {
	cf                string
	key               *fieldMapping
	columns           []*fieldMapping
	value             *fieldMapping
	others            map[string]*fieldMapping
	isCompositeColumn bool
	isSliceColumn     bool
	isStarNameColumn  bool
}

func defaultCassandraType(t reflect.Type) (TypeDesc, int) {
	switch t.Kind() {
	case reflect.Bool:
		return BooleanType, baseTypeField
	case reflect.String:
		return UTF8Type, baseTypeField
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return LongType, baseTypeField
	case reflect.Float32:
		return FloatType, baseTypeField
	case reflect.Float64:
		return DoubleType, baseTypeField
	case reflect.Array:
		if t.Name() == "UUID" && t.Size() == 16 {
			return UUIDType, baseTypeField
		}
		return UnknownType, baseTypeField
	case reflect.Slice:
		if et := t.Elem(); et.Kind() == reflect.Uint8 {
			return BytesType, baseTypeField
		} else {
			if subTD, subKind := defaultCassandraType(et); subTD != UnknownType && subKind == baseTypeField {
				return subTD, baseTypeSliceField
			}
			return UnknownType, baseTypeField
		}
		return UnknownType, baseTypeField
	}
	return UnknownType, baseTypeField
}

func newFieldMapping(pos int, sf reflect.StructField, overrideName, overrideType string) *fieldMapping {
	fm := &fieldMapping{}
	fm.cassandraType, fm.fieldKind = defaultCassandraType(sf.Type)
	// signal an invalid Go base type
	if fm.cassandraType == UnknownType {
		return nil
	}
	if overrideType != "" {
		fm.cassandraType = parseTypeDesc(overrideType)
	}
	fm.position = pos
	if overrideName != "" {
		fm.cassandraName = overrideName
	} else {
		fm.cassandraName = sf.Name
	}
	fm.name = sf.Name
	return fm
}

func newStructMapping(t reflect.Type) (*structMapping, error) {
	sm := &structMapping{}
	n := t.NumField()
	found := false
	// globally recognized meta fields in the struct tag
	meta := map[string]string{
		"cf":  "",
		"key": "",
		"col": "",
		"val": "",
	}
	// hold a field mapping for every candidate field
	fields := make(map[string]*fieldMapping)

	// pass 1: gather field metadata
	for i := 0; i < n; i++ {
		sf := t.Field(i)
		// find the field tags
		for key, _ := range meta {
			if tagValue := sf.Tag.Get(key); tagValue != "" {
				meta[key] = tagValue
			}
		}
		// build a field mapping for all non-anon named fields with a suitable Go type
		if sf.Name != "" && !sf.Anonymous {
			if fm := newFieldMapping(i, sf, sf.Tag.Get("name"), sf.Tag.Get("type")); fm == nil {
				continue
			} else {
				fields[sf.Name] = fm
			}
		}
	}

	// pass 2: struct data for each meta field
	if name := meta["cf"]; name != "" {
		sm.cf = meta["cf"]
	} else {
		return nil, errors.New(fmt.Sprint("No cf field in struct ", t.Name()))
	}

	if name := meta["key"]; name != "" {
		if sm.key, found = fields[name]; !found {
			return nil, errors.New(fmt.Sprint("Referenced key field ", name, " does not exist in struct ", t.Name()))
		}
		if sm.key.fieldKind != baseTypeField {
			return nil, errors.New(fmt.Sprint("Referenced key field ", name, " in struct ", t.Name(), " has invalid type"))
		}
		delete(fields, name)
	} else {
		return nil, errors.New(fmt.Sprint("No key field in struct ", t.Name()))
	}

	if name := meta["val"]; (name != "") || (name == "*value") {
		if name == "*value" {
			sm.value = &fieldMapping{fieldKind: starValueField}
		} else if sm.value, found = fields[name]; !found {
			return nil, errors.New(fmt.Sprint("Referenced value field ", name, " does not exist in struct ", t.Name()))
		}
		delete(fields, name)
	} else {
		return nil, errors.New(fmt.Sprint("No val field in struct ", t.Name()))
	}

	if meta["col"] != "" {
		colNames := strings.Split(meta["col"], ",")
		for i, name := range colNames {
			isLast := i == (len(colNames) - 1)
			var fm *fieldMapping
			if name == "*name" {
				if !isLast {
					return nil, errors.New(fmt.Sprint("*name can only be used in the last position of a composite, error in struct ", t.Name()))
				} else {
					sm.isStarNameColumn = true
					fm = &fieldMapping{fieldKind: starNameField}
				}
			} else if fm, found = fields[name]; !found {
				return nil, errors.New(fmt.Sprint("Referenced column field ", name, " does not exist in struct ", t.Name()))
			}
			if fm.fieldKind == baseTypeSliceField {
				sm.isSliceColumn = true
				if !isLast {
					return nil, errors.New(fmt.Sprint("Slice struct fields can only be used in the last position of a composite, error in struct ", t.Name()))
				}
			}
			delete(fields, name)
			sm.columns = append(sm.columns, fm)
		}
		sm.isCompositeColumn = len(sm.columns) > 1
		sm.others = make(map[string]*fieldMapping)
		for _, fm := range fields {
			sm.others[fm.cassandraName] = fm
		}
	} else {
		return nil, errors.New(fmt.Sprint("No col field in struct ", t.Name()))
	}

	return sm, nil
}

var mapCache map[reflect.Type]*structMapping
var mapCacheMutex *sync.Mutex = new(sync.Mutex)

func getMapping(v reflect.Value) (*structMapping, error) {
	var sm *structMapping
	var err error
	found := false
	t := v.Type()
	mapCacheMutex.Lock()
	if mapCache == nil {
		mapCache = make(map[reflect.Type]*structMapping)
	}
	if sm, found = mapCache[t]; !found {
		sm, err = newStructMapping(t)
		if err != nil {
			mapCache[t] = sm
		}
	}
	mapCacheMutex.Unlock()
	return sm, err
}

type mappedStruct struct {
	source interface{}
	v      reflect.Value
	sm     *structMapping
}

func newMappedStruct(source interface{}) (*mappedStruct, error) {
	ms := &mappedStruct{source: source}
	var err error

	// always work with a pointer to struct
	vp := reflect.ValueOf(source)
	if vp.Kind() != reflect.Ptr {
		return nil, errors.New("Passed source is not a pointer to a struct")
	}
	if vp.IsNil() {
		return nil, errors.New("Passed source is not a pointer to a struct")
	}
	ms.v = reflect.Indirect(vp)
	if ms.v.Kind() != reflect.Struct {
		return nil, errors.New("Passed source is not a pointer to a struct")
	}
	if !ms.v.CanSet() {
		return nil, errors.New("Cannot modify the passed struct instance")
	}

	ms.sm, err = getMapping(ms.v)
	if err != nil {
		return nil, err
	}

	return ms, nil
}

func (ms *mappedStruct) marshalKey() ([]byte, error) {
	vk := ms.v.Field(ms.sm.key.position)
	b, err := Marshal(vk.Interface(), ms.sm.key.cassandraType)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Error marshaling key field ", ms.sm.key.name, ":", err))
	}
	return b, nil
}

func internalMap(source interface{}) (*Row, *mappedStruct, error) {
	// deconstruct the source struct into a reflect.Value and a (cached) struct mapping
	ms, err := newMappedStruct(source)
	if err != nil {
		return nil, nil, err
	}

	// allocate new row to return
	row := &Row{}

	// marshal the key field
	b, err := ms.marshalKey()
	if err != nil {
		return nil, nil, err
	}
	row.Key = b

	// marshal columns and values
	return row, ms, ms.mapField(row, 0, make([]byte, 0), make([]byte, 0), 0)
}

func Map(source interface{}) (*Row, error) {
	row, _, err := internalMap(source)
	return row, err
}

func (ms *mappedStruct) mapField(row *Row, component int, composite []byte, value []byte, valueIndex int) error {

	// check if there are components left
	if component < len(ms.sm.columns) {

		fm := ms.sm.columns[component]

		// switch type of field named by component
		switch fm.fieldKind {

		// base type
		case baseTypeField:
			// set value of the current composite field to the field value
			v := ms.v.Field(fm.position)
			b, err := Marshal(v.Interface(), fm.cassandraType)
			if err != nil {
				return errors.New(fmt.Sprint("Error marshaling field ", fm.name, ":", err))
			}
			if ms.sm.isCompositeColumn {
				composite = packComposite(composite, b, false, false, false)
			} else {
				composite = b
			}
			return ms.mapField(row, component+1, composite, value, valueIndex)

		// slice of base type
		case baseTypeSliceField:
			// iterate slice and map more columns
			v := ms.v.Field(fm.position)
			n := v.Len()
			for i := 0; i < n; i++ {
				// set value of the current composite field to the field value
				vi := v.Index(i)
				b, err := Marshal(vi.Interface(), fm.cassandraType)
				if err != nil {
					return errors.New(fmt.Sprint("Error marshaling field ", fm.name, ":", err))
				}
				var subComposite []byte
				if ms.sm.isCompositeColumn {
					subComposite = packComposite(composite, b, false, false, false)
				} else {
					subComposite = b
				}
				err = ms.mapField(row, component+1, subComposite, value, i)
				if err != nil {
					return err
				}
			}

		// *name
		case starNameField:
			// iterate over non-key/col/val-referenced struct fields and map more columns
			for _, fm := range ms.sm.others {
				// set value of the current composite field to the field name (possibly overriden by name:)
				b, err := Marshal(fm.cassandraName, UTF8Type)
				if err != nil {
					return errors.New(fmt.Sprint("Error marshaling field ", fm.name, ":", err))
				}
				var subComposite []byte
				if ms.sm.isCompositeColumn {
					subComposite = packComposite(composite, b, false, false, false)
				} else {
					subComposite = b
				}
				// marshal field value and pass it to next field mapper in case it is *value
				v := ms.v.Field(fm.position)
				b, err = Marshal(v.Interface(), fm.cassandraType)
				if err != nil {
					return errors.New(fmt.Sprint("Error marshaling field ", fm.name, ":", err))
				}

				err = ms.mapField(row, component+1, subComposite, b, valueIndex)
				if err != nil {
					return err
				}
			}
		}

	} else {
		// no components left, emit column

		fm := ms.sm.value

		// switch type of value field
		switch fm.fieldKind {

		case starValueField:
			// use passed value
			row.Columns = append(row.Columns, &Column{Name: composite, Value: value})

		case baseTypeSliceField:
			// set value to the passed value index in this slice
			vs := ms.v.Field(fm.position)
			v := vs.Index(valueIndex)
			b, err := Marshal(v.Interface(), fm.cassandraType)
			if err != nil {
				return errors.New(fmt.Sprint("Error marshaling field ", fm.name, ":", err))
			}
			row.Columns = append(row.Columns, &Column{Name: composite, Value: b})

		case baseTypeField:
			// set value to the field value
			v := ms.v.Field(fm.position)
			b, err := Marshal(v.Interface(), fm.cassandraType)
			if err != nil {
				return errors.New(fmt.Sprint("Error marshaling field ", fm.name, ":", err))
			}
			row.Columns = append(row.Columns, &Column{Name: composite, Value: b})

			// support literal case?

			// support zero, non-set case?
		}
	}

	return nil
}

func Unmap(row *Row, destination interface{}) error {

	// deconstruct the source struct into a reflect.Value and a (cached) struct mapping
	ms, err := newMappedStruct(destination)
	if err != nil {
		return err
	}

	// unmarshal key
	vk := ms.v.Field(ms.sm.key.position)
	if !vk.CanAddr() {
		return errors.New("Cannot obtain pointer to key field")
	}
	vkp := vk.Addr()
	err = Unmarshal(row.Key, ms.sm.key.cassandraType, vkp.Interface())
	if err != nil {
		return errors.New(fmt.Sprint("Error unmarshaling key field ", ms.sm.key.name, ":", err))
	}

	// unmarshal col/values

	setField := func(fm *fieldMapping, b []byte, index int) error {
		vfield := ms.v.Field(fm.position)
		if index >= 0 {
			vfield = vfield.Index(index)
		}
		if !vfield.CanAddr() {
			return errors.New(fmt.Sprint("Cannot obtain pointer to field ", vfield.Type().Name(), " in struct ", ms.v.Type().Name()))
		}
		vfieldp := vfield.Addr()

		err = Unmarshal(b, fm.cassandraType, vfieldp.Interface())
		if err != nil {
			return errors.New(fmt.Sprint("Error unmarshaling composite field ", vfield.Type().Name(), " in struct ", ms.v.Type().Name(), ", error: ", err))
		}
		return nil
	}

	prepareSlice := func(fm *fieldMapping, n int) {
		vfield := ms.v.Field(fm.position)
		t := vfield.Type()
		s := reflect.MakeSlice(t, n, n)
		vfield.Set(s)
	}

	rowLength := len(row.Columns)

	// prepare slice components and value
	for _, fm := range ms.sm.columns {
		if fm.fieldKind == baseTypeSliceField {
			prepareSlice(fm, rowLength)
		}
	}
	if ms.sm.value.fieldKind == baseTypeSliceField {
		prepareSlice(ms.sm.value, rowLength)
	}

	for i, column := range row.Columns {
		var components [][]byte
		if ms.sm.isCompositeColumn {
			components = unpackComposite(column.Name)
		} else {
			components = [][]byte{column.Name}
		}
		if len(components) != len(ms.sm.columns) {
			return errors.New(fmt.Sprint("Returned number of components in composite column name does not match struct col: component in struct ", ms.v.Type().Name()))
		}

		// iterate over column name components and set them, plus values
		for j, b := range components {
			fm := ms.sm.columns[j]
			switch fm.fieldKind {

			case baseTypeField:
				if err = setField(fm, b, -1); err != nil {
					return err
				}

			case starNameField:
				var name string
				err = Unmarshal(b, UTF8Type, &name)
				if err != nil {
					return errors.New(fmt.Sprint("Error unmarshaling composite field as UTF8Type for *name in struct ", ms.v.Type().Name(), ", error: ", err))
				}
				if valueFM, found := ms.sm.others[name]; found {
					if err = setField(valueFM, column.Value, -1); err != nil {
						return err
					}
				}

			case baseTypeSliceField:
				if err = setField(fm, b, i); err != nil {
					return err
				}
			}
		}

		// set value field for the non-*name cases
		if ms.sm.value != nil {
			switch ms.sm.value.fieldKind {

			case baseTypeField:
				if err = setField(ms.sm.value, column.Value, -1); err != nil {
					return err
				}

			case baseTypeSliceField:
				if err = setField(ms.sm.value, column.Value, i); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
