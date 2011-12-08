package model

import (
	enc "encoding/binary"
)

type BaseValue interface {
	Bytes() []byte
	SetBytes([]byte)
}

/*
type Row interface {
	Key() BaseValue
	SetKey(BaseValue)
	Pairs() []Pair
}

type example struct {
	id Bytes "KEY"
	name Bytes
	address Bytes
	email Bytes
	cookie Bytes
}
*/


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

func (l *Long) Bytes() []byte {
	b := make([]byte, 8)
	enc.BigEndian.PutUint64(b, uint64(*l))
	return b
}

func (l *Long) SetBytes(b []byte)  {
	*l = Long(enc.BigEndian.Uint64(b))
}


// "strings" CQL blob/ascii/text

type Bytes string

func (u *Bytes) Bytes() []byte {
	return []byte(string(*u))
}

func (u *Bytes) SetBytes(b []byte)  {
	*u = Bytes(string(b[0:(len(b))]))
}
