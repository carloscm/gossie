package gossie

import (
    "os"
    "fmt"
    //"thrift"
    "encoding/hex"
    //Cassandra "cassandra"
)


/////////////////////////////////////
// Queries

type Query interface {
    Run() 
}

type Row interface {
}

// all the get queries support column sets and slices
type BaseGetter interface {
    ConsistencyLevel(int)
    Cf(string)
    Columns([]Value)
    Slice(start, end Value, limit int, reverse bool)
    Run() ([]Row, os.Error)
}

// use when you know the key(s) of your rows
type GetSlicer interface {
    BaseGetter
    Key(Value)
    Keys([]Value)
}

// use when you know the values for a secondary index and want to query by it
/*
type GetIndexer interface {
    BaseGetter
    Eq(column, value Value)
}
*/

// use when you know a range of keys you want to iterate (when using random partitioner
// this is only useful for iterating over an entire CF)
type GetRanger interface {
    BaseGetter
    Range(start, end Value, limit int)
}

type BatchMutator interface {
    ConsistencyLevel(int)
    Insert(cf string, key Value, row Row)
    Delete(cf string, key Value)
    DeleteSlice(cf string, key Value, start, end Value, limit int)
    DeleteColumns(cf string, key Value, row Row)
    Run() os.Error
}

type common struct {
    readConsistency int
    writeConsistency int
    cf string
}

func (c *common) ConsistencyLevel(read, write int) {
    c.readConsistency = read
    c.writeConsistency = write
}

func (c *common) Cf(cf string) {
    c.cf = cf
}

type columnValue struct {
    column []byte
    value []byte
}

type insertOps struct {
    common
    key []byte
    columns []*columnValue
}

/*
func (c *connection) Insert() InsertOps {
    o := &insertOps{}
    o.conn = c
    return o
}
*/

func (o *insertOps) Key(v Value) {
    o.key = v.Bytes()
}

func (o *insertOps) Column(column, value Value) {
    o.columns = append(o.columns, &columnValue{ column.Bytes(), value.Bytes() } )
}

func quoteBytes(v []byte) string {
    return fmt.Sprint("'", hex.EncodeToString(v), "'")
}

func (o *insertOps) Run() os.Error {


/*
recordar que timestamp es obligatorio!

    BatchMutate (map<'key', map<'cf', list<Mutation>>, ...)
        Mutation:
            ColumnOrSuperColumn -> Column -> name, value, timestamp, ttl
            OR Deletion -> timestamp, ..., SlicePredicate

*/

    if o.key == nil {
        return newError("Missing key in insert op")
    }

    if len(o.columns) == 0 {
        return newError("Must pass at least one column in insert op")
    }

    cols, vals := "", ""
    for _, cv := range o.columns {
        cols = fmt.Sprint(cols, ", ", quoteBytes(cv.column))
        vals = fmt.Sprint(vals, ", ", quoteBytes(cv.value))
    }

    //INSERT INTO users (KEY, password) VALUES ('jsmith', 'ch@ngem3a') USING TTL 86400;

    

    q := fmt.Sprint("INSERT INTO ", o.cf, " (KEY", cols, ") VALUES (", quoteBytes(o.key), vals, ")")

    fmt.Println(q)

    //_, _, _, _, _, err := o.conn.client.ExecuteCqlQuery([]byte(q), Cassandra.NONE)

    //o.columns = append(o.columns, &columnValue{ column.Bytes(), value.Bytes() } )
    return nil //err
}
