package gossie

import (
	"testing"
)

func TestSchema(t *testing.T) {

	c, err := NewConnection("127.0.0.1:9160", "TestGossie", 3000)
	if err != nil {
		t.Fatal("Error connecting to Cassandra:", err)
	}

	schema := newSchema(c.(*connection))

/*
	if len(schema.ColumnFamilies) != 3 {
		t.Error("Test schema must have 3 CFs")
	}

	if schema.ColumnFamilies["AllTypes"] == nil {
		t.Error("Test CF AllTypes is nil")
	} else {
		cf := schema.ColumnFamilies["AllTypes"]
    	if cf.DefaultComparator
    	cf.DefaultValidator
    	cf.KeyValidator

    	cf.NamedColumns map[string]TypeDesc
*/	
}

/*
	op := c.Insert()
	op.Cf("Users")

	k := Bytes("user3")
	op.Key(&k)

	ct := Bytes("name")
	vt := Bytes("hehehe")
	op.Column(&ct, &vt)
	op.Run()

	c.Close()
	*/
}

