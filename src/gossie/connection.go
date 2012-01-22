package gossie

import (
    "net"
    "os"
    //"fmt"
    "thrift"
    //"encoding/hex"
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

const (
    DEFAULT_SIZE = 10
    DEFAULT_READ_CONSISTENCY = CONSISTENCY_QUORUM
    DEFAULT_WRITE_CONSISTENCY = CONSISTENCY_QUORUM
    DEFAULT_TIMEOUT = 3000
    //DEFAULT_RECYCLE = 60
)

/////////////////////////////////////
// Error

type Error string

func (e *Error) String() string {
    return string(*e)
}

func newError(e string) *Error {
    return (*Error)(&e)
}

/////////////////////////////////////
// ConnectionPool

type ConnectionPool interface {
    Acquire() Connection
    Release(Connection)
}

type PoolOptions struct {
    Size int                // keep up to Size connections open and ready
    ReadConsistency int     // default read consistency
    WriteConsistency int    // default write consistency
    Timeout int             // socket timeout in ms
    //Recycle int             // close the connection after Recycle seconds no matter what (it it's unused)
}

func (o *PoolOptions) defaults() {
    if o.Size == 0 { o.Size = DEFAULT_SIZE }
    if o.ReadConsistency == 0 { o.ReadConsistency = DEFAULT_READ_CONSISTENCY }
    if o.WriteConsistency == 0 { o.WriteConsistency = DEFAULT_WRITE_CONSISTENCY }
    if o.Timeout == 0 { o.Timeout = DEFAULT_TIMEOUT }
    //if o.Recycle == 0 { o.Recycle = DEFAULT_RECYCLE }
}

type connectionPool struct {
    keyspace string
    hosts []string
    options PoolOptions

    //available []Connection
    //inUse []Connection
}

func NewConnectionPool(keyspace string, hosts []string, options PoolOptions) (ConnectionPool, os.Error) {

    options.defaults()

    cp := &connectionPool {
        keyspace: keyspace,
        hosts: hosts,
        options: options,
        //available: make([]*connection, 0, size),
        //inUse: make([]*connection, 0, size)
    }

    return cp, nil
}


func (*connectionPool) Acquire() Connection {
    return nil
}

func (*connectionPool) Release(c Connection) {
    
}


/////////////////////////////////////
// Connection

type Connection interface {
    //Execute(Query) Result
    Close()
    Keyspace() string
    Client() *Cassandra.CassandraClient
}

type connection struct {
    socket *thrift.TNonblockingSocket
    transport *thrift.TFramedTransport
    client *Cassandra.CassandraClient
    keyspace string
    lastUsage int
}

func NewConnection(hostPort, keyspace string, timeout int) (Connection, os.Error) {

    addr, err := net.ResolveTCPAddr("tcp", hostPort)
    if err != nil {
        return nil, err
    }

    c := &connection{}

    c.socket, err = thrift.NewTNonblockingSocketAddr(addr)
    if err != nil {
        return nil, err
    }

    // socket not open yet, so no error expected
    c.socket.SetTimeout(int64(timeout) * 1000000000)

    c.transport = thrift.NewTFramedTransport(c.socket)
    protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
    c.client = Cassandra.NewCassandraClientFactory(c.transport, protocolFactory)

    err = c.transport.Open()
    if err != nil {
        return nil, err
    }

    ire, err := c.client.SetKeyspace(keyspace)

    if err != nil {
        return nil, err
    }

    if ire != nil {
        return nil, newError("Cannot set the keyspace")
    }

    c.keyspace = keyspace

    return c, nil
}

func (c *connection) Keyspace() string {
    return c.keyspace
}

func (c *connection) Close() {
    c.transport.Close()
}

func (c *connection) Client() *Cassandra.CassandraClient {
    return c.client
}
