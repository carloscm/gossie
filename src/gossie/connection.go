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

const (
    DEFAULT_READ_CONSISTENCY = CONSISTENCY_QUORUM
    DEFAULT_WRITE_CONSISTENCY = CONSISTENCY_QUORUM
    DEFAULT_TIMEOUT = 3000
    DEFAULT_RECYCLE = 60
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

type connectionPool struct {
    keyspace string
    readConsistency int
    writeConsistency int
    hosts string[]
    timeout int
    recycle int
    size int
    available Connection[]
    inUse Connection[]
}

func NewConnectionPoolDefaults(keyspace string, hosts []string, size int) (ConnectionPool, os.Error) {
    return NewConnectionPool(keyspace, hostPort, size, DEFAULT_READ_CONSISTENCY, DEFAULT_WRITE_CONSISTENCY, DEFAULT_TIMEOUT, DEFAULT_RECYCLE)
}

func NewConnectionPool(keyspace string, hosts []string, size, readConsistency, writeConsistency, timeout, recycle int) (ConnectionPool, os.Error) {

    cp := &connectionPool {
        keyspace: keyspace
        readConsistency: readConsistency
        writeConsistency: writeConsistency
        hosts: hosts
        timeout: timeout
        recycle: recycle
        size: size
        available: make([]*connection, 0, size)
        inUse: make([]*connection, 0, size)
    }

    return cp, nil
}

/*
func (*connectionPool) Acquire() Connection {
    
}

func (*connectionPool) Release(Connection c) {
    
}
*/

/////////////////////////////////////
// Connection

type Connection interface {
    Execute(Query) Result
    Close()
}

type connection struct {
    socket *thrift.TNonblockingSocket
    transport *thrift.TFramedTransport
    client *Cassandra.CassandraClient
    lastUsage int
}

func NewConnection(hostPort, keyspace string, timeout int) (Connection, os.Error) {

    addr, err := net.ResolveTCPAddr("tcp", hostPort)
    if err != nil {
        return nil, err
    }

    c := &connection{}

    c.socket, err := thrift.NewTNonblockingSocketAddr(addr)
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

    return c, nil
}

func (c *connection) Close() {
    c.transport.Close()
}
