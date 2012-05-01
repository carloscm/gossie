package gossie

import (
	"errors"
	"fmt"
)

/*
	todo:
	compact struct mapping

	ideas:
	map to Go maps
	map to Go slices (N slices?)
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
	v, err := validStruct(source)
	if err != nil {
		return -1
	}
	si, err := inspectStruct(v)
	if err != nil {
		return -1
	}
	// struct fields minus the components fields minus one field for the key
	return len(si.goFields) - len(m.components) - 1
}

func (m *sparseMapping) Map(source interface{}) (*Row, error) {
	v, err := validStruct(source)
	if err != nil {
		return nil, err
	}
	si, err := inspectStruct(v)
	if err != nil {
		return nil, err
	}

	row := &Row{}

	// marshal the key field
	if f, found := si.goFields[m.key]; found {
		b, err := f.marshalValue(v)
		if err != nil {
			return nil, err
		}
		row.Key = b
	} else {
		return nil, errors.New(fmt.Sprint("Mapping key field ", m.key, " not found in passed struct of type ", v.Type().Name()))
	}

	// prepare composite, if needed
	composite := make([]byte, 0)
	for _, c := range m.components {
		if f, found := si.goFields[c]; found {
			b, err := f.marshalValue(v)
			if err != nil {
				return nil, err
			}
			composite = append(composite, packComposite(b, eocEquals)...)
		} else {
			return nil, errors.New(fmt.Sprint("Mapping component field ", c, " not found in passed struct of type ", v.Type().Name()))
		}
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

func (m *sparseMapping) Unmap(destination interface{}, offset int, row *Row) (int, error) {
	readColumns := 0

	v, err := validStruct(destination)
	if err != nil {
		return readColumns, err
	}
	si, err := inspectStruct(v)
	if err != nil {
		return readColumns, err
	}

	// unmarshal key field
	if f, found := si.goFields[m.key]; found {
		err = f.unmarshalValue(row.Key, v)
		if err != nil {
			return readColumns, err
		}
	} else {
		return readColumns, errors.New(fmt.Sprint("Mapping key field ", m.key, " not found in passed struct of type ", v.Type().Name()))
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
		var components [][]byte
		if len(m.components) > 0 {
			components = unpackComposite(column.Name)
		} else {
			components = [][]byte{column.Name}
		}
		if len(components) != (len(m.components) + 1) {
			return readColumns, errors.New(fmt.Sprint("Returned number of components in composite column name does not match struct mapping in struct ", v.Type().Name()))
		}

		// FIXME: it is possible for a row to contain multiple composite
		// values instead of an uniform one, indicating that a new "object"
		// started. assume that is not the case for now!
		// iterate over composite components, just once, to set the composite
		// fields
		if !compositeFieldsAreSet {
			for i, c := range m.components {
				if f, found := si.goFields[c]; found {
					b := components[i]
					err := f.unmarshalValue(b, v)
					if err != nil {
						return readColumns, err
					}
				} else {
					return readColumns, errors.New(fmt.Sprint("Mapping component field ", c, " not found in passed struct of type ", v.Type().Name()))
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
		if f, found := si.cassandraFields[name]; found {
			err := f.unmarshalValue(column.Value, v)
			if err != nil {
				return readColumns, errors.New(fmt.Sprint("Error unmarshaling column: ", name, " value: ", err))
			}
		}
	}

	return readColumns, nil
}

// Compact returns a mapping for Go structs that represents a Cassandra row
// column as a full Go struct. The field named by the value is mapped to the
// column value. Each passed component field name is mapped, in order, to the
// column composite values.
/*
func NewCompact(cf string, keyField string, valueField string, ...componentFields string) Mapping {
	return &compactMapping {
		cf: cf,
		key: keyField,
		components: componentFields
		value: valueField
	}
}

type compactMapping struct {
	cf string
	key string
	components []string
	value string
}
*/
