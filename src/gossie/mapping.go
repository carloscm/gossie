package gossie

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	. "github.com/wadey/gossie/src/cassandra"
)

/*
	ideas:
	mapping for Go maps
	mapping for Go slices (N slices?)
*/

// Mapping maps the type of a Go object to/from a Cassandra row.
type Mapping interface {

	// Cf returns the column family name
	Cf() string

	// Returns true for compact mapping
	Compact() bool

	Components() []string

	// MarshalKey marshals the passed key value into a []byte
	MarshalKey(key interface{}) ([]byte, error)

	// MarshalField marshals the single field value into a []byte
	MarshalField(field string, value interface{}) ([]byte, error)

	// Unmarshal single field into value pointer
	UnmarshalField(field string, data []byte, valuep interface{}) error

	// MarshalComponent marshals the passed component value at the position into a []byte
	MarshalComponent(component interface{}, position int) ([]byte, error)

	// Map converts a Go object compatible with this Mapping into a Row
	Map(source interface{}) (*Row, error)

	// Ummap fills the passed Go object with data from a row
	Unmap(destination interface{}, provider RowProvider) error
}

var (
	EndBeforeLimit = errors.New("No more results found before reaching the limit")
	EndAtLimit     = errors.New("No more results found but reached the limit")
)

// RowProvider abstracts the details of reading a series of columns from a Cassandra row
type RowProvider interface {

	// Key returns the row key
	Key() ([]byte, error)

	// NextColumn returns the next column in the row, and advances the column pointer
	NextColumn() (*Column, error)

	// Rewind moves back the column pointer one position
	Rewind()
}

var (
	noMoreComponents = errors.New("No more components allowed")
)

// NewMapping looks up the field tag 'mapping' in the passed struct type
// to decide which mapping it is using, then builds a mapping using the 'cf',
// 'key', 'cols' and 'value' field tags.
func NewMapping(source interface{}) (Mapping, error) {
	_, si, err := validateAndInspectStruct(source)
	if err != nil {
		return nil, err
	}

	cf, found := si.globalTags["cf"]
	if !found {
		return nil, errors.New(fmt.Sprint("Mandatory struct tag 'cf' not found in passed struct of type ", si.rtype.Name()))
	}

	key, found := si.globalTags["key"]
	if !found {
		return nil, errors.New(fmt.Sprint("Mandatory struct tag 'key' not found in passed struct of type ", si.rtype.Name()))
	}
	_, found = si.goFields[key]
	if !found {
		return nil, errors.New(fmt.Sprint("Key field ", key, " not found in passed struct of type ", si.rtype.Name()))
	}

	colsS := []string{}
	cols, found := si.globalTags["cols"]
	if found {
		colsS = strings.Split(cols, ",")
	}
	for _, c := range colsS {
		_, found := si.goFields[c]
		if !found {
			return nil, errors.New(fmt.Sprint("Composite field ", c, " not found in passed struct of type ", si.rtype.Name()))
		}
	}

	value, found := si.globalTags["value"]
	if found {
		_, found := si.goFields[value]
		if !found {
			return nil, errors.New(fmt.Sprint("Value field ", value, " not found in passed struct of type ", si.rtype.Name()))
		}
	}

	mapping, found := si.globalTags["mapping"]
	if !found {
		mapping = "sparse"
	}

	switch mapping {
	case "sparse":
		return newSparseMapping(si, cf, key, colsS...), nil
	case "compact":
		return newCompactMapping(si, cf, key, value, colsS...), nil
	}

	return nil, errors.New(fmt.Sprint("Unrecognized mapping type ", mapping, " in passed struct of type ", si.rtype.Name()))
}

func MustNewMapping(source interface{}) Mapping {
	ret, err := NewMapping(source)
	if err != nil {
		panic(err)
	}
	return ret
}

func newSparseMapping(si *structInspection, cf string, keyField string, componentFields ...string) Mapping {
	cm := make(map[string]bool, 0)
	for _, f := range componentFields {
		cm[f] = true
	}
	return &sparseMapping{
		si:            si,
		cf:            cf,
		key:           keyField,
		components:    componentFields,
		componentsMap: cm,
	}
}

type sparseMapping struct {
	si            *structInspection
	cf            string
	key           string
	components    []string
	componentsMap map[string]bool
}

func (m *sparseMapping) Cf() string {
	return m.cf
}

func (m *sparseMapping) Compact() bool {
	return false
}

func (m *sparseMapping) Components() []string {
	return m.components
}

func (m *sparseMapping) MarshalField(field string, value interface{}) ([]byte, error) {
	f, ok := m.si.goFields[field]
	if !ok {
		return nil, fmt.Errorf("No such field %s", field)
	}
	b, err := Marshal(value, f.cassandraType)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Error marshaling passed value for field ", f.name, ":", err))
	}
	return b, nil
}

func (m *sparseMapping) UnmarshalField(field string, b []byte, valuep interface{}) error {
	f, ok := m.si.goFields[field]
	if !ok {
		return fmt.Errorf("No such field %s", field)
	}
	err := Unmarshal(b, f.cassandraType, valuep)
	if err != nil {
		return errors.New(fmt.Sprint("Error unmarshaling passed value for field ", f.name, ":", err))
	}
	return nil
}

func (m *sparseMapping) MarshalKey(key interface{}) ([]byte, error) {
	f := m.si.goFields[m.key]
	b, err := Marshal(key, f.cassandraType)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Error marshaling passed value for the key in field ", f.name, ":", err))
	}
	return b, nil
}

func (m *sparseMapping) MarshalComponent(component interface{}, position int) ([]byte, error) {
	if position >= len(m.components) {
		return nil, errors.New(fmt.Sprint("The mapping has a component length of ", len(m.components), " and the passed position is ", position))
	}
	f := m.si.goFields[m.components[position]]
	b, err := Marshal(component, f.cassandraType)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Error marshaling passed value for a composite component in field ", f.name, ":", err))
	}
	return b, nil
}

func (m *sparseMapping) startMap(source interface{}, compact bool) (*Row, *reflect.Value, *structInspection, []byte, error) {
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
	var composite []byte
	if compact && len(m.components) == 1 {
		c := m.components[0]
		if f, found := si.goFields[c]; found {
			composite, err = f.marshalValue(v)
			if err != nil {
				return nil, nil, nil, nil, err
			}
		} else {
			return nil, nil, nil, nil, errors.New(fmt.Sprint("Mapping component field ", c, " not found in passed struct of type ", v.Type().Name()))
		}
	} else {
		composite = make([]byte, 0)
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
	}

	return row, v, si, composite, nil
}

func (m *sparseMapping) Map(source interface{}) (*Row, error) {
	row, v, si, composite, err := m.startMap(source, false)
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
		if f.skipEmpty && f.isEmpty(v) {
			continue
		}
		columnValue, err := f.marshalValue(v)
		if err != nil {
			return nil, err
		}
		row.Columns = append(row.Columns, &Column{Name: columnName, Value: columnValue})
	}

	return row, nil
}

func (m *sparseMapping) startUnmap(destination interface{}, provider RowProvider) (*reflect.Value, *structInspection, error) {
	v, si, err := validateAndInspectStruct(destination)
	if err != nil {
		return nil, nil, err
	}

	// unmarshal key field
	if f, found := si.goFields[m.key]; found {
		key, err := provider.Key()
		if err != nil {
			return nil, nil, err
		}
		err = f.unmarshalValue(key, v)
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

func (m *sparseMapping) extractComponents(column *Column, v *reflect.Value, bias int) ([][]byte, error) {
	var components [][]byte
	if len(m.components) > 0 {
		components = unpackComposite(column.Name)
	} else {
		components = [][]byte{column.Name}
	}
	if len(components) != (len(m.components) + bias) {
		return components, errors.New(fmt.Sprint("Returned number of components in composite column name does not match struct mapping in struct ", v.Type().Name()))
	}
	return components, nil
}

// TODO: speed this up
func (m *sparseMapping) isNewComponents(prev, next [][]byte, bias int) bool {
	if len(prev) != len(next) {
		return true
	}
	for i := 0; i < len(prev)-bias; i++ {
		p := prev[i]
		n := next[i]
		if len(p) != len(n) {
			return true
		}
		for j := 0; j < len(p); j++ {
			if p[j] != n[j] {
				return true
			}
		}
	}
	return false
}

func (m *sparseMapping) Unmap(destination interface{}, provider RowProvider) error {
	v, si, err := m.startUnmap(destination, provider)
	if err != nil {
		return err
	}

	compositeFieldsAreSet := false
	var previousComponents [][]byte

	for {
		column, err := provider.NextColumn()
		if err == Done {
			return Done
		} else if err == EndBeforeLimit {
			if compositeFieldsAreSet {
				break
			} else {
				return Done
			}
		} else if err == EndAtLimit {
			return Done
		} else if err != nil {
			return err
		}

		components, err := m.extractComponents(column, v, 1)
		if err != nil {
			return err
		}
		if !compositeFieldsAreSet {
			// first column
			if err := m.unmapComponents(v, si, components); err != nil {
				return err
			}
			compositeFieldsAreSet = true
		} else {
			if m.isNewComponents(previousComponents, components, 1) {
				provider.Rewind()
				break
			}
		}

		// lookup field by name
		var name string
		err = Unmarshal(components[len(components)-1], UTF8Type, &name)
		if err != nil {
			return errors.New(fmt.Sprint("Error unmarshaling composite field as UTF8Type for field name in struct ", v.Type().Name(), ", error: ", err))
		}
		if f, found := si.cassandraFields[name]; found {
			if column.Value != nil {
				err := f.unmarshalValue(column.Value, v)
				if err != nil {
					return errors.New(fmt.Sprint("Error unmarshaling column: ", name, " value: ", err))
				}
			}
		}

		previousComponents = components
	}

	return nil
}

func newCompactMapping(si *structInspection, cf string, keyField string, valueField string, componentFields ...string) Mapping {
	return &compactMapping{
		sparseMapping: *(newSparseMapping(si, cf, keyField, componentFields...).(*sparseMapping)),
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

func (m *compactMapping) Compact() bool {
	return true
}

func (m *compactMapping) Map(source interface{}) (*Row, error) {
	row, v, si, composite, err := m.startMap(source, true)
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
		row.Columns = append(row.Columns, &Column{Name: composite, Value: make([]byte, 0, 0)})
	}
	return row, nil
}

func (m *compactMapping) Unmap(destination interface{}, provider RowProvider) error {
	v, si, err := m.startUnmap(destination, provider)
	if err != nil {
		return err
	}

	column, err := provider.NextColumn()
	if err == Done {
		return Done
	} else if err == EndBeforeLimit {
		return Done
	} else if err == EndAtLimit {
		return Done
	} else if err != nil {
		return err
	}

	if len(m.components) == 1 {
		c := m.components[0]
		if f, found := si.goFields[c]; found {
			err := f.unmarshalValue(column.Name, v)
			if err != nil {
				return err
			}
		} else {
			return errors.New(fmt.Sprint("Mapping component field ", c, " not found in passed struct of type ", v.Type().Name()))
		}
	} else {
		components, err := m.extractComponents(column, v, 0)
		if err != nil {
			return err
		}
		if err := m.unmapComponents(v, si, components); err != nil {
			return err
		}
	}
	if f, found := si.goFields[m.value]; found {
		if column.Value != nil {
			err := f.unmarshalValue(column.Value, v)
			if err != nil {
				return errors.New(fmt.Sprint("Error unmarshaling column for compact value: ", err))
			}
		}
	}

	return nil
}
