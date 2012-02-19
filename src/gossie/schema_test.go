package gossie

import (
	"testing"
)

func TestSchema(t *testing.T) {

	c, err := newConnection("127.0.0.1:9160", "TestGossie", 3000)
	if err != nil {
		t.Fatal("Error connecting to Cassandra:", err)
	}

	schema := newSchema(c)
	defer c.close()

	if len(schema.ColumnFamilies) != 3 {
		t.Error("Test schema must have 3 CFs")
	}

	if schema.ColumnFamilies["AllTypes"] == nil {
		t.Error("Test CF AllTypes is nil")
	} else {
		cf := schema.ColumnFamilies["AllTypes"]

    	if cf.DefaultComparator.Desc != AsciiType {
    		t.Error("Test CF AllTypes DefaultComparator is not AsciiType")
    	}
    	if cf.DefaultValidator.Desc != UTF8Type {
    		t.Error("Test CF AllTypes DefaultValidator is not UTF8Type")
    	}
    	if cf.KeyValidator.Desc != BytesType {
    		t.Error("Test CF AllTypes KeyValidator is not BytesType")
    	}

    	var check = map[string]TypeDesc {
			"colBytesType" : BytesType,
			"colAsciiType" : AsciiType,
			"colUTF8Type" : UTF8Type,
			"colLongType" : LongType,
			"colIntegerType" : IntegerType,
			"colDecimalType" : DecimalType,
			"colUUIDType" : UUIDType,
			"colBooleanType" : BooleanType,
			"colFloatType" : FloatType,
			"colDoubleType" : DoubleType,
			"colDateType" : DateType,
    	}

		if len(cf.NamedColumns) != len(check) {
			t.Error("Test CF AllTypes contains an unexpected amount of named columns")
		}

    	for name, desc := range check {
    		if cf.NamedColumns[name].Desc != desc {
    			t.Error("Test CF AllTypes column ", name, " is not the expected type")
    		}
    	}

    }

    if schema.ColumnFamilies["Counters"] == nil {
		t.Error("Test CF Counters is nil")
	} else {
		cf := schema.ColumnFamilies["Counters"]

    	if cf.DefaultComparator.Desc != AsciiType {
    		t.Error("Test CF Counters DefaultComparator is not AsciiType")
    	}
    	if cf.DefaultValidator.Desc != CounterColumnType {
    		t.Error("Test CF Counters DefaultValidator is not CounterColumnType")
    	}
    	if cf.KeyValidator.Desc != BytesType {
    		t.Error("Test CF Counters KeyValidator is not BytesType")
    	}

		if len(cf.NamedColumns) != 0 {
			t.Error("Test CF Counters has named columns")
		}
    }

    if schema.ColumnFamilies["Composite"] == nil {
		t.Error("Test CF Composite is nil")
	} else {
		cf := schema.ColumnFamilies["Composite"]

    	if cf.DefaultComparator.Desc != CompositeType {
    		t.Error("Test CF Composite DefaultComparator is not CompositeType")
    	}
    	if cf.DefaultValidator.Desc != BytesType {
    		t.Error("Test CF Composite DefaultValidator is not BytesType")
    	}
    	if cf.KeyValidator.Desc != BytesType {
    		t.Error("Test CF Composite KeyValidator is not BytesType")
    	}

		if len(cf.NamedColumns) != 0 {
			t.Error("Test CF Composite has named columns")
		}

		var check = []TypeDesc {BytesType, AsciiType, UTF8Type, LongType, IntegerType, DecimalType, UUIDType, BooleanType, FloatType, DoubleType, DateType}

		if len(cf.DefaultComparator.Components) != len(check) {
			t.Error("Test CF Composite has incorrect number of components")
		}

    	for i, desc := range check {
    		if cf.DefaultComparator.Components[i].Desc != desc {
    			t.Error("Test CF Composite comparator has incorrect comparator in position ", i)
    		}
    	}

    }

}

