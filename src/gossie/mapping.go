package gossie

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

/*
	ideas:
	mapping for Go maps
	mapping for Go slices (N slices?)
*/

// Mapping maps the type of a Go object to/from (a slice of) a Cassandra row.
type Mapping interface {

	// Cf returns the column family name
	Cf() string

	// MinColumns returns the minimal number of columns required by the mapped Go
	// object
	MinColumns(source interface{}) int

	// Map converts a Go object compatible with this Mapping into a Row
	Map(source interface{}) (*Row, error)

	// Ummap fills the passed Go object with data from the row, starting at the
	// offset column. It returns the count of consumed columns, or -1 if there
	// wasn't enough columns to fill the Go object
	Unmap(destination interface{}, offset int, row *Row) (int, error)
}

// MappingFromTags looks up the field tag 'mapping' in the passed struct type
// to decide which mapping it is using, then builds a mapping using the 'cf',
// 'key', 'cols' and 'value' field tags.
func MappingFromTags(source interface{}) (Mapping, error) {
	_, si, err := validateAndInspectStruct(source)
	if err != nil {
		return nil, err
	}
	// mandatory tags for all the provided mappings
	for _, t := range []string{"cf", "key"} {
		_, found := si.globalTags[t]
		if !found {
			return nil, errors.New(fmt.Sprint("Mandatory struct tag ", t, " not found in passed struct of type ", si.rtype.Name()))
		}
	}
	// optional tags
	colsS := []string{}
	cols, found := si.globalTags["cols"]
	if found {
		colsS = strings.Split(cols, ",")
	}
	mapping, found := si.globalTags["mapping"]
	if !found {
		mapping = "sparse"
	}
	value := si.globalTags["value"]

	switch mapping {
	case "sparse":
		return NewSparse(si.globalTags["cf"], si.globalTags["key"], colsS...), nil
	case "compact":
		if value == "" {
			return nil, errors.New(fmt.Sprint("Mandatory struct tag value for compact mapping not found in passed struct of type ", si.rtype.Name()))
		}
		return NewCompact(si.globalTags["cf"], si.globalTags["key"], si.globalTags["value"], colsS...), nil
	}

	return nil, errors.New(fmt.Sprint("Unrecognized mapping type ", mapping, " in passed struct of type ", si.rtype.Name()))
}

// Sparse returns a mapping for Go structs that represents a Cassandra row key
// as a struct field, zero or more composite column names as zero or more
// struct fields, and the rest of the struct fields as extra columns with the
// name being the last composite column name, and the value the column value.
func NewSparse(cf string, keyField string, componentFields ...string) Mapping {
	cm := make(map[string]bool, 0)
	for _, f := range componentFields {
		cm[f] = true
	}
	return &sparseMapping{
		cf:            cf,
		key:           keyField,
		components:    componentFields,
		componentsMap: cm,
	}
}

type sparseMapping struct {
	cf            string
	key           string
	components    []string
	componentsMap map[string]bool
}

func (m *sparseMapping) Cf() string {
	return m.cf
}

func (m *sparseMapping) MinColumns(source interface{}) int {
	_, si, err := validateAndInspectStruct(source)
	if err != nil {
		return -1
	}
	// struct fields minus the components fields minus one field for the key
	return len(si.goFields) - len(m.components) - 1
}

func (m *sparseMapping) startMap(source interface{}) (*Row, *reflect.Value, *structInspection, []byte, error) {
	v, si, err := validateAndInspectStruct(source)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	row := &Row{}

	// marshal the key field
	if f, found := si.goFields[m.key]; found {
		b, err := f.marshalValue(v)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		row.Key = b
	} else {
		return nil, nil, nil, nil, errors.New(fmt.Sprint("Mapping key field ", m.key, " not found in passed struct of type ", v.Type().Name()))
	}

	// prepare composite, if needed
	composite := make([]byte, 0)
	for _, c := range m.components {
		if f, found := si.goFields[c]; found {
			b, err := f.marshalValue(v)
			if err != nil {
				return nil, nil, nil, nil, err
			}
			composite = append(composite, packComposite(b, eocEquals)...)
		} else {
			return nil, nil, nil, nil, errors.New(fmt.Sprint("Mapping component field ", c, " not found in passed struct of type ", v.Type().Name()))
		}
	}

	return row, v, si, composite, nil
}

func (m *sparseMapping) Map(source interface{}) (*Row, error) {
	row, v, si, composite, err := m.startMap(source)
	if err != nil {
		return nil, err
	}

	// add columns
	for _, f := range si.orderedFields {
		if f.name == m.key {
			continue
		}
		if _, found := m.componentsMap[f.name]; found {
			continue
		}
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

func (m *sparseMapping) startUnmap(destination interface{}, row *Row) (*reflect.Value, *structInspection, error) {
	v, si, err := validateAndInspectStruct(destination)
	if err != nil {
		return nil, nil, err
	}

	// unmarshal key field
	if f, found := si.goFields[m.key]; found {
		err = f.unmarshalValue(row.Key, v)
		if err != nil {
			return nil, nil, err
		}
	} else {
		return nil, nil, errors.New(fmt.Sprint("Mapping key field ", m.key, " not found in passed struct of type ", v.Type().Name()))
	}

	return v, si, nil
}

func (m *sparseMapping) unmapComponents(v *reflect.Value, si *structInspection, components [][]byte) error {
	for i, c := range m.components {
		if f, found := si.goFields[c]; found {
			b := components[i]
			err := f.unmarshalValue(b, v)
			if err != nil {
				return err
			}
		} else {
			return errors.New(fmt.Sprint("Mapping component field ", c, " not found in passed struct of type ", v.Type().Name()))
		}
	}
	return nil
}

func (m *sparseMapping) extractComponents(column *Column, v *reflect.Value, biasN int) ([][]byte, error) {
	var components [][]byte
	if len(m.components) > 0 {
		components = unpackComposite(column.Name)
	} else {
		components = [][]byte{column.Name}
	}
	if len(components) != (len(m.components) + biasN) {
		return components, errors.New(fmt.Sprint("Returned number of components in composite column name does not match struct mapping in struct ", v.Type().Name()))
	}
	return components, nil
}

func (m *sparseMapping) Unmap(destination interface{}, offset int, row *Row) (int, error) {
	readColumns := 0

	v, si, err := m.startUnmap(destination, row)
	if err != nil {
		return readColumns, err
	}

	compositeFieldsAreSet := false

	// unmarshal col/values
	min := m.MinColumns(destination)
	if min > len(row.Columns) {
		return -1, nil
	}
	columns := row.Columns[offset : offset+min]
	for _, column := range columns {
		readColumns++
		components, err := m.extractComponents(column, v, 1)
		if err != nil {
			return readColumns, err
		}

		// FIXME: it is possible for a row to contain multiple composite
		// values instead of an uniform one, indicating that a new "object"
		// started. assume that is not the case for now!
		// iterate over composite components, just once, to set the composite
		// fields
		if !compositeFieldsAreSet {
			if err := m.unmapComponents(v, si, components); err != nil {
				return readColumns, err
			}
			compositeFieldsAreSet = true
		}

		// lookup field by name
		var name string
		err = Unmarshal(components[len(components)-1], UTF8Type, &name)
		if err != nil {
			return readColumns, errors.New(fmt.Sprint("Error unmarshaling composite field as UTF8Type for field name in struct ", v.Type().Name(), ", error: ", err))
		}
		if f, found := si.cassandraFields[name]; found {
			err := f.unmarshalValue(column.Value, v)
			if err != nil {
				return readColumns, errors.New(fmt.Sprint("Error unmarshaling column: ", name, " value: ", err))
			}
		}
	}

	return readColumns, nil
}

// NewCompact returns a mapping for Go structs that represents a Cassandra row
// column as a full Go struct. The field named by the value is mapped to the
// column value. Each passed component field name is mapped, in order, to the
// column composite values.
func NewCompact(cf string, keyField string, valueField string, componentFields ...string) Mapping {
	return &compactMapping{
		sparseMapping: *(NewSparse(cf, keyField, componentFields...).(*sparseMapping)),
		value:         valueField,
	}
}

type compactMapping struct {
	sparseMapping
	value string
}

func (m *compactMapping) Cf() string {
	return m.cf
}

func (m *compactMapping) MinColumns(source interface{}) int {
	return 1
}

func (m *compactMapping) Map(source interface{}) (*Row, error) {
	row, v, si, composite, err := m.startMap(source)
	if err != nil {
		return nil, err
	}
	if f, found := si.goFields[m.value]; found {
		columnValue, err := f.marshalValue(v)
		if err != nil {
			return nil, err
		}
		row.Columns = append(row.Columns, &Column{Name: composite, Value: columnValue})
	} else {
		return nil, errors.New(fmt.Sprint("Mapping value field ", m.value, " not found in passed struct of type ", v.Type().Name()))
	}
	return row, nil
}

func (m *compactMapping) Unmap(destination interface{}, offset int, row *Row) (int, error) {
	v, si, err := m.startUnmap(destination, row)
	if err != nil {
		return 0, err
	}
	if len(row.Columns) <= 0 {
		return -1, nil
	}
	column := row.Columns[offset]
	components, err := m.extractComponents(column, v, 0)
	if err != nil {
		return 1, err
	}
	if err := m.unmapComponents(v, si, components); err != nil {
		return 1, err
	}
	if f, found := si.goFields[m.value]; found {
		err := f.unmarshalValue(column.Value, v)
		if err != nil {
			return 1, errors.New(fmt.Sprint("Error unmarshaling column for compact value: ", err))
		}
	} else {
		return 1, errors.New(fmt.Sprint("Mapping value field ", m.value, " not found in passed struct of type ", v.Type().Name()))
	}
	return 1, nil
}
