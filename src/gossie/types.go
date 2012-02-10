package gossie

import (
    //"fmt"
    //"cassandra"
    enc "encoding/binary"
    "strings"
)

const (
    BytesType = 1
    AsciiType
    UTF8Type
    LongType
    IntegerType
    DecimalType
    UUIDType
    BooleanType
    FloatType
    DoubleType
    DateType
)

func Marshal(value interface{}, typeDesc int) []byte {

    switch v := value.(type) {
        default:
            return nil

        case bool:
            return marshalBool(t, typeDesc)

        case int:
            return marshalInt(t, typeDesc)

        case string:
            return marshalString(t, typeDesc)

    }
}

/*
        case BytesType:
        case AsciiType:
        case UTF8Type:
        case LongType:
        case IntegerType:
        case DecimalType:
        case UUIDType:
        case BooleanType:
        case FloatType:
        case DoubleType:
        case DateType:
*/

func marshalBool(value bool, typeDesc int) []byte {
    switch typeDesc {
        case BytesType, BooleanType:
            b = make([]byte, 1)
            if value {
                b[0] = 1
            } else {
                b[0] = 0
            }
            return b

        case AsciiType, UTF8Type:
            b = make([]byte, 1)
            if value {
                b[0] = '1'
            } else {
                b[0] = '0'
            }
            return b

        case LongType:
            b = make([]byte, 8)
            if value {
                b[7] = 1
            } else {
                b[7] = 0
            }
            return b

        /* unsuported marshalling:
        case IntegerType:
        case DecimalType:
        case UUIDType:
        case FloatType:
        case DoubleType:
        case DateType:
        */
    }
    return nil
}


func marshalInt(value int, typeDesc int) []byte {
    switch typeDesc {
        case BytesType, BooleanType:
    }
    return nil
}


type Value interface {
    Bytes() []byte
    SetBytes([]byte)
}

type TypeDesc interface {
    Validate(v Value) bool
}

type Bytes string
func (u *Bytes) Bytes() []byte {
    return []byte(string(*u))
}
func (u *Bytes) SetBytes(b []byte)  {
    *u = Bytes(string(b[0:(len(b))]))
}

type bytesTypeDesc struct {}
func (u *bytesTypeDesc) Validate(v Value) bool {
    _, ok := v.(*Bytes)
    return ok
}

type Long int64
func (l *Long) Bytes() []byte {
    b := make([]byte, 8)
    enc.BigEndian.PutUint64(b, uint64(*l))
    return b
}
func (l *Long) SetBytes(b []byte)  {
    *l = Long(enc.BigEndian.Uint64(b))
}

type longTypeDesc struct {}
func (u *longTypeDesc) Validate(v Value) bool {
    _, ok := v.(*Long)
    return ok
}

type compositeTypeDesc struct {
    components []TypeDesc
}
func (u *compositeTypeDesc) Validate(v Value) bool {
    return false
}

func makeTypeDesc(cassType string) TypeDesc {

    // not a simple class type, check for composite and parse it
    if (strings.HasPrefix(cassType, "org.apache.cassandra.db.marshal.CompositeType(")) {
        composite := &compositeTypeDesc{}
        componentsString := cassType[strings.Index(cassType, "(")+1:len(cassType)-1]
        componentsSlice := strings.Split(componentsString, ",")
        components := make([]TypeDesc, 0)
        for _, component := range componentsSlice {
            components = append(components, makeTypeDesc(component))
        }
        composite.components = components
        return composite
    }

    // simple types
    switch cassType {
        case "org.apache.cassandra.db.marshal.LongType":
            return &longTypeDesc{}
    }

    // not a recognized type
    return &bytesTypeDesc{}
}

//BytesType
//AsciiType
//UTF8Type
//LongType
//IntegerType  Arbitrary-precision integer
//DecimalType decimal Variable-precision decimal
//UUIDType
//CounterColumnType
//BooleanType
//FloatType
//DoubleType
//DateType
