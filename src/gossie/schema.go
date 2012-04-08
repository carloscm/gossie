package gossie

import (
	//"fmt"
	//"thrift"
	//"encoding/hex"
	"cassandra"
)

/*
to do:
    generate CQL schema from tagged Go structs
    validate tagged Go structs against schemas
    handle ReversedType
    handle type options
	handle composited column names in the schema (is this in use/allowed?)
*/

type Schema struct {
	ColumnFamilies map[string]*ColumnFamily
}

type ColumnFamily struct {
	DefaultComparator TypeClass
	DefaultValidator  TypeClass
	KeyValidator      TypeClass
	NamedColumns      map[string]TypeClass
}

func newSchema(ksDef *cassandra.KsDef) *Schema {
	cfDefs := ksDef.CfDefs
	schema := &Schema{ColumnFamilies: make(map[string]*ColumnFamily)}

	for cfDefT := range cfDefs.Iter() {

		// FIXME: this is weird, but happens a lot. thrift4go problem?
		if cfDefT == nil {
			continue
		}

		cfDef, _ := cfDefT.(*cassandra.CfDef)

		if cfDef.ColumnType != "Standard" {
			continue
		}

		cf := &ColumnFamily{}

		cf.DefaultComparator = parseTypeClass(cfDef.ComparatorType)
		cf.DefaultValidator = parseTypeClass(cfDef.DefaultValidationClass)
		cf.KeyValidator = parseTypeClass(cfDef.KeyValidationClass)

		cf.NamedColumns = make(map[string]TypeClass)

		for colDefT := range cfDef.ColumnMetadata.Iter() {
			// FIXME: this is weird, but happens a lot. thrift4go problem?
			if colDefT == nil {
				continue
			}
			colDef, _ := colDefT.(*cassandra.ColumnDef)
			name := string(colDef.Name[0:(len(colDef.Name))])
			cf.NamedColumns[name] = parseTypeClass(colDef.ValidationClass)
		}

		schema.ColumnFamilies[cfDef.Name] = cf
	}

	//fmt.Println(schema)

	return schema
}
