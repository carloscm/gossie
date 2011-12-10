package model

import (
    "net"
    "os"
    "fmt"
    "thrift"
    "encoding/hex"
    Cassandra "cassandra"
)

const (
    CONSISTENCY_DEFAULT = 0
    CONSISTENCY_ZERO = 1
    CONSISTENCY_ONE = 2
    CONSISTENCY_QUORUM = 3
    CONSISTENCY_ALL = 4
    CONSISTENCY_LOCAL_QUORUM = 5
    CONSISTENCY_EACH_QUORUM = 6
)

type Error string

func (e *Error) String() string {
    return string(*e)
}

func newError(e string) *Error {
    return (*Error)(&e)
}

type Row interface {
}

type KeyOps interface {
    Key(BaseValue)
    Keys([]BaseValue)
}

type CommonOps interface {
    Cf(string)
    ConsistencyLevel(int, int)
}

type GetOps interface {
    CommonOps
    Column(BaseValue)
    Columns([]BaseValue)
    Slice(start, end BaseValue, limit int, reverse bool)
    Run() []Row
}

type RangeOps interface {
    Range(start, end BaseValue, limit int)
}

type IndexedOps interface {
    Eq(column, value BaseValue)
}

type GetSliceOps interface {
    //GetOps
    KeyOps

    Run() []Row
}

type GetIndexedOps interface {
    GetOps
    Eq(column, value BaseValue)
}

type GetRangeOps interface {
    GetOps
    RangeOps
}

type InsertOps interface {
    CommonOps
    Key(BaseValue)
    Column(column, value BaseValue)
    Run() os.Error
}

type Connection interface {
    //Get() GetSliceOps
    //GetIndexed() GetIndexedOps
    //GetRange() GetRangeOps
    Insert() InsertOps
    Close()

    //defer trans.Close()

}

type connection struct {
    transport *thrift.TFramedTransport
    client *Cassandra.CassandraClient
}

func NewConnection(hostPort, keyspace string) (Connection, os.Error) {
    c := &connection{}

    addr, err := net.ResolveTCPAddr("tcp", hostPort)
    if err != nil {
        return nil, err
    }

    fmt.Println("1")

    transport, err := thrift.NewTNonblockingSocketAddr(addr)
    if err != nil {
        return nil, err
    }

    fmt.Println("2")

    c.transport = thrift.NewTFramedTransport(transport)
    protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
    c.client = Cassandra.NewCassandraClientFactory(c.transport, protocolFactory)

    fmt.Println("3")

    err = c.transport.Open()
    if err != nil {
        return nil, err
    }

    fmt.Println("4")

    _, _, _, _, _, err = c.client.ExecuteCqlQuery([]byte(fmt.Sprint("USE ", keyspace)), Cassandra.NONE)
    if err != nil {
        return nil, err
    }

    fmt.Println("5")

    return c, nil
}

func (c *connection) Close() {
    c.transport.Close()
}

type common struct {
    conn *connection
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

func (c *connection) Insert() InsertOps {
    o := &insertOps{}
    o.conn = c
    return o
}

func (o *insertOps) Key(v BaseValue) {
    o.key = v.Bytes()
}

func (o *insertOps) Column(column, value BaseValue) {
    o.columns = append(o.columns, &columnValue{ column.Bytes(), value.Bytes() } )
}

func quoteBytes(v []byte) string {
    return fmt.Sprint("'", hex.EncodeToString(v), "'")
}

func (o *insertOps) Run() os.Error {

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

    _, _, _, _, _, err := o.conn.client.ExecuteCqlQuery([]byte(q), Cassandra.NONE)

    //o.columns = append(o.columns, &columnValue{ column.Bytes(), value.Bytes() } )
    return err
}
