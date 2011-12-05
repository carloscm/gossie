package model

// column types must implement io.Reader and io.Writer

//type Bytes string	// CQL blob
//type Ascii string	// CQL ascii
type UTF8 string	// CQL text
//type Long int		// CQL int, bigint
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

func (s *UTF8) Read(p []byte) (n int, err os.Error) {
	
}

type Entity interface {
	begin(id string)
	pair(key string, value ColumnType)
	end()
}
