package gossie

import (
    //"fmt"
    //"cassandra"
    enc "encoding/binary"
    //"strings"
    "strconv"
    "os"
)

const (
    _ = iota
    BytesType
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
    CounterColumnType
)

var (
    ErrorUnsupportedMarshaling = os.NewError("Cannot marshal value")
    ErrorUnsupportedUnmarshaling = os.NewError("Cannot unmarshal value")
)

type TypeDesc int

/*
    to do:

    FloatType
    DoubleType
    IntegerType
    DecimalType
    UUIDType

    float32
    float64
    all uints

    maybe some more (un)marshalings, maybe something better for DateType?

    more error checking, pass along strconv errors
*/

func Marshal(value interface{}, typeDesc TypeDesc) ([]byte, os.Error) {
    // dereference in case we got a pointer
    var dvalue interface{}
    switch v := value.(type) {
        case *[]byte:   dvalue = *v
        case *bool:     dvalue = *v
        case *int8:     dvalue = *v
        case *int16:    dvalue = *v
        case *int:      dvalue = *v
        case *int32:    dvalue = *v
        case *int64:    dvalue = *v
        case *string:   dvalue = *v
        default:        dvalue = v
    }

    switch v := dvalue.(type) {
        case []byte:    return v, nil
        case bool:      return marshalBool(v, typeDesc)
        case int8:      return marshalInt(int64(v), 1, typeDesc)
        case int16:     return marshalInt(int64(v), 2, typeDesc)
        case int:       return marshalInt(int64(v), 4, typeDesc)
        case int32:     return marshalInt(int64(v), 4, typeDesc)
        case int64:     return marshalInt(v, 8, typeDesc)
        case string:    return marshalString(v, typeDesc)
    }
    return nil, ErrorUnsupportedMarshaling
}

func marshalBool(value bool, typeDesc TypeDesc) ([]byte, os.Error) {
    switch typeDesc {
        case BytesType, BooleanType:
            b := make([]byte, 1)
            if value {
                b[0] = 1
            }
            return b, nil

        case AsciiType, UTF8Type:
            b := make([]byte, 1)
            if value {
                b[0] = '1'
            } else {
                b[0] = '0'
            }
            return b, nil

        case LongType:
            b := make([]byte, 8)
            if value {
                b[7] = 1
            }
            return b, nil
    }
    return nil, ErrorUnsupportedMarshaling
}

func marshalInt(value int64, size int, typeDesc TypeDesc) ([]byte, os.Error) {
    switch typeDesc {

        case LongType:
            b := make([]byte, 8)
            enc.BigEndian.PutUint64(b, uint64(value))
            return b, nil

        case BytesType:
            b := make([]byte, 8)
            enc.BigEndian.PutUint64(b, uint64(value))
            return b[len(b)-size:], nil

        case DateType:
            if size != 8 {
                return nil, ErrorUnsupportedMarshaling
            }
            b := make([]byte, 8)
            enc.BigEndian.PutUint64(b, uint64(value))
            return b, nil

        case AsciiType, UTF8Type:
            return marshalString(strconv.Itoa64(value), UTF8Type)
    }
    return nil, ErrorUnsupportedMarshaling
}

func marshalString(value string, typeDesc TypeDesc) ([]byte, os.Error) {
    // let cassandra check the ascii-ness of the []byte
    switch typeDesc {
        case BytesType, AsciiType, UTF8Type:
            return []byte(value), nil

        case LongType:
            i, err := strconv.Atoi64(value)
            if err != nil {
                return nil, err
            }
            return marshalInt(i, 8, LongType)

/* fix this!
        case UUIDType:
            if len(value) != 36 {
                return nil, ErrorUnsupportedMarshaling
            }
            ints := strings.Split(value, "-")
            if len(ints) != 5 {
                return nil, ErrorUnsupportedMarshaling
            }
            b := marshalInt(strconv.Btoi64(ints[0], 16), 4, BytesType)
            b = append(b, marshalInt(strconv.Btoi64(ints[1], 16), 2, BytesType))
            b = append(b, marshalInt(strconv.Btoi64(ints[2], 16), 2, BytesType))
            b = append(b, marshalInt(strconv.Btoi64(ints[3], 16), 2, BytesType))
            b = append(b, marshalInt(strconv.Btoi64(ints[4], 16), 6, BytesType))
            return b, nil
*/

    }
    return nil, ErrorUnsupportedMarshaling
}

func Unmarshal(b []byte, typeDesc TypeDesc, value interface{}) os.Error {
    switch v := value.(type) {
        case *[]byte:    *v = b; return nil
        case *bool:      return unmarshalBool(b, typeDesc, v)
        //case *int8:      return unmarshalInt(b, v, 1, typeDesc)
        //case *int16:     return unmarshalInt(b, v, 2, typeDesc)
        //case *int:       return unmarshalInt(b, v, 4, typeDesc)
        //case *int32:     return unmarshalInt(b, v, 4, typeDesc)
        //case *int64:     return unmarshalInt(b, v, 8, typeDesc)
        //case *string:    return unmarshalString(b, v, typeDesc)
    }
    return ErrorUnsupportedUnmarshaling
}

func unmarshalBool(b []byte, typeDesc TypeDesc, value *bool) os.Error {
    switch typeDesc {
        case BytesType, BooleanType:
            if len(b) < 1 {
                return ErrorUnsupportedUnmarshaling
            }
            if b[0] == 0 {
                *value = false
            } else {
                *value = true
            }
            return nil

        case AsciiType, UTF8Type:
            if len(b) < 1 {
                return ErrorUnsupportedUnmarshaling
            }
            if b[0] == '0' {
                *value = false
            } else {
                *value = true
            }
            return nil

        case LongType:
            if len(b) != 8 {
                return ErrorUnsupportedUnmarshaling
            }
            if b[7] == 0 {
                *value = false
            } else {
                *value = true
            }
            return nil
    }
    return ErrorUnsupportedUnmarshaling
}

type Value interface {
    Bytes() []byte
    SetBytes([]byte)
    //TypeDesc() TypeDesc
}

func makeTypeDesc(cassType string) TypeDesc {

    // not a simple class type, check for composite and parse it
    /* disable composite support for now...
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
    */

    // simple types
    switch cassType {
        case "org.apache.cassandra.db.marshal.BytesType":      return BytesType
        case "org.apache.cassandra.db.marshal.AsciiType":      return AsciiType
        case "org.apache.cassandra.db.marshal.UTF8Type":       return UTF8Type
        case "org.apache.cassandra.db.marshal.LongType":       return LongType
        case "org.apache.cassandra.db.marshal.IntegerType":    return IntegerType
        case "org.apache.cassandra.db.marshal.DecimalType":    return DecimalType
        case "org.apache.cassandra.db.marshal.UUIDType":       return UUIDType
        case "org.apache.cassandra.db.marshal.BooleanType":    return BooleanType
        case "org.apache.cassandra.db.marshal.FloatType":      return FloatType
        case "org.apache.cassandra.db.marshal.DoubleType":     return DoubleType
        case "org.apache.cassandra.db.marshal.DateType":       return DateType
    }

    // not a recognized type
    return BytesType
}
