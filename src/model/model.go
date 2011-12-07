package model

import (
	enc "encoding/binary"
	//"strings"
)

type BaseValue interface {
	Bytes() []byte
	SetBytes([]byte)
}

type Pair interface {
	Name() BaseValue
	Value() BaseValue
	TTL() Long
	Timestamp() Long
}

type Row interface {
	Key() BaseValue
	Pairs() []Pair
}

/*type Entity struct {
	
}*/


//type Bytes string	// CQL blob
//type Ascii string	// CQL ascii
//type UUID string	// CQL uuid
//type Float float	// CQL float
//type Double double	// CQL double

/*
missing:
IntegerType	varint	Arbitrary-precision integer
UUIDType	uuid	Type 1 or type 4 UUID
DateType	timestamp	Date plus time, encoded as 8 bytes since epoch
BooleanType	boolean	true or false
DecimalType	decimal	Variable-precision decimal
CounterColumnType	counter	Distributed counter value (8-byte long)
*/

// Long

type Long int64	// CQL int, bigint

func (l Long) Bytes() []byte {
	b := make([]byte, 8)
	enc.BigEndian.PutUint64(b, uint64(l))
	return b
}

func (l *Long) SetBytes(b []byte)  {
	*l = Long(enc.BigEndian.Uint64(b))
}


// UTF8

type UTF8 string	// CQL text

func (u UTF8) Bytes() []byte {
	return []byte(string(u))
}

func (u *UTF8) SetBytes(b []byte)  {
	*u = UTF8(string(b[0:(len(b))]))
}
