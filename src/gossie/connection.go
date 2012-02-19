package gossie

import (
    "net"
    "os"
    "thrift"
    "cassandra"
)

/*

    to do:

    failover and retry!

    auth
    more pooling options

*/

const (
    CONSISTENCY_DEFAULT = 0
    CONSISTENCY_ONE = 1
    CONSISTENCY_QUORUM = 2
    CONSISTENCY_LOCAL_QUORUM = 3
    CONSISTENCY_EACH_QUORUM = 4
    CONSISTENCY_ALL = 5
    CONSISTENCY_ANY = 6
    CONSISTENCY_TWO = 7
    CONSISTENCY_THREE = 8
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
    Keyspace() string
    Schema() *Schema
    Query() Query
    //Mutation() Mutation
    Close()
}

type PoolOptions struct {
    Size int                // keep up to Size connections open and ready
    ReadConsistency int     // default read consistency
    WriteConsistency int    // default write consistency
    Timeout int             // socket timeout in ms
    //Recycle int             // close the connection after Recycle seconds no matter what (if it's unused)
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
    schema *Schema

    hack *connection

    //available []Connection
    //inUse []Connection
}

func NewConnectionPool(hosts []string, keyspace string, options PoolOptions) (ConnectionPool, os.Error) {

    options.defaults()

    cp := &connectionPool {
        keyspace: keyspace,
        hosts: hosts,
        options: options,
        //available: make([]*connection, 0, size),
        //inUse: make([]*connection, 0, size)
    }

    cp.hack, _ = newConnection(hosts[0], keyspace, cp.options.Timeout)

    c := cp.acquire()
    if (c == nil) {
        return nil, os.NewError("Cannot acquire initial connection")
    }
    cp.schema = newSchema(c)
    cp.release(c)

    return cp, nil
}


func (cp *connectionPool) acquire() *connection {
    return cp.hack
}

func (cp *connectionPool) release(c *connection) {
    
}

func (cp *connectionPool) releaseWithTimeout(c *connection) {
    
}

func (cp *connectionPool) releaseWithError(c *connection) {
    
}

func (cp *connectionPool) Query() Query {
    return &query{consistencyLevel:cp.options.ReadConsistency, pool:cp}
}

func (cp *connectionPool) Keyspace() string {
    return cp.keyspace
}

func (cp *connectionPool) Schema() *Schema {
    return cp.schema
}


func (cp *connectionPool) Close() {
}

/////////////////////////////////////
// connection

type connection struct {
    socket *thrift.TNonblockingSocket
    transport *thrift.TFramedTransport
    client *cassandra.CassandraClient
    keyspace string
    lastUsage int
}

func newConnection(hostPort, keyspace string, timeout int) (*connection, os.Error) {

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
    c.client = cassandra.NewCassandraClientFactory(c.transport, protocolFactory)

    err = c.transport.Open()
    if err != nil {
        return nil, err
    }

    ire, err := c.client.SetKeyspace(keyspace)

    if err != nil {
        c.close()
        return nil, err
    }

    if ire != nil {
        c.close()
        return nil, newError("Cannot set the keyspace")
    }

    c.keyspace = keyspace

    return c, nil
}

func (c *connection) close() {
    c.transport.Close()
}
