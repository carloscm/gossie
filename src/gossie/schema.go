package gossie

import (
    "fmt"
    //"thrift"
    //"encoding/hex"
    "cassandra"
)

type Schema struct {
    ColumnFamilies map[string]*ColumnFamily
}

type ColumnFamily struct {
    DefaultComparator TypeDesc
    DefaultValidator TypeDesc
    KeyValidator TypeDesc
    NamedColumns map[string]TypeDesc
}

func newSchema(c *connection) *Schema {

    ksDef, nfe, ire, err := c.client.DescribeKeyspace(c.keyspace)

    if ksDef == nil || nfe != nil || ire != nil || err != nil {
        return nil
    }

    cfDefs := ksDef.CfDefs


    schema := &Schema{ColumnFamilies:make(map[string]*ColumnFamily)}

    for cfDefT := range cfDefs.Iter() {

        // FIXME: this is weird, but happens a lot. thrift4go problem?
        if cfDefT == nil {
            continue
        }

        // FIXME: add support for counter CFs
        cfDef, _ := cfDefT.(*cassandra.CfDef)

        if cfDef.ColumnType != "Standard" {
            continue
        }

        cf := &ColumnFamily{}

        cf.DefaultComparator = makeTypeDesc(cfDef.ComparatorType)
        cf.DefaultValidator = makeTypeDesc(cfDef.DefaultValidationClass)
        cf.KeyValidator = makeTypeDesc(cfDef.KeyValidationClass)

        cf.NamedColumns = make(map[string]TypeDesc)

        for colDefT := range cfDef.ColumnMetadata.Iter() {
            // FIXME: this is weird, but happens a lot. thrift4go problem?
            if colDefT == nil {
                continue
            }
            colDef, _ := colDefT.(*cassandra.ColumnDef)
            name := string(colDef.Name[0:(len(colDef.Name))])
            cf.NamedColumns[name] = makeTypeDesc(colDef.ValidationClass)
        }

        schema.ColumnFamilies[cfDef.Name] = cf
    }

    //fmt.Println(schema)

    return nil
}

