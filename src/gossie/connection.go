package gossie

import (
	"errors"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/apesternikov/gossie/src/cassandra"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

/*
   to do:
   auth
   timeout while waiting for an available connection slot
   panic handling inside run()?
   maybe more pooling options
   Close()
*/

// ConnectionPool implements a pool of Cassandra connections to one or more nodes
type ConnectionPool interface {
	// Keyspace returns the keyspace name this ConnectionPool is connected to
	Keyspace() string

	// Schema returns the parsed schema for the keyspace this ConnectionPool is connected to
	Schema() *Schema

	// Reader returns a new query builder for read operations
	Reader() Reader

	// Writer returns a new mutation builder for write operations
	Writer() Writer

	// Query returns a high level interface for read operations over structs
	Query(Mapping) Query

	// Batch returns a high level interface for write operations over structs
	Batch() Batch

	// Close all the connections in the pool
	Close()
}

//package-private methods
type connectionRunner interface {
	run(t transaction) error
	runWithRetries(t transaction, retries int) error
}

// PoolOptions stores the options for the creation of a ConnectionPool
type PoolOptions struct {
	Size             int                        // keep up to Size connections open and ready
	ReadConsistency  cassandra.ConsistencyLevel // default read consistency
	WriteConsistency cassandra.ConsistencyLevel // default write consistency
	Timeout          time.Duration              // socket timeout
	Grace            int                        // if a node is blacklisted try to contact it again after Grace seconds
	Retries          int                        // retry queries for Retries times before raising an error
	Authentication   map[string]string          // if one or more keys are present, login() is called with the values from Authentication
}

const (
	CONSISTENCY_DEFAULT      cassandra.ConsistencyLevel = 0
	CONSISTENCY_ONE                                     = cassandra.ConsistencyLevel_ONE
	CONSISTENCY_QUORUM                                  = cassandra.ConsistencyLevel_QUORUM
	CONSISTENCY_LOCAL_QUORUM                            = cassandra.ConsistencyLevel_LOCAL_QUORUM
	CONSISTENCY_EACH_QUORUM                             = cassandra.ConsistencyLevel_EACH_QUORUM
	CONSISTENCY_ALL                                     = cassandra.ConsistencyLevel_ALL
	CONSISTENCY_ANY                                     = cassandra.ConsistencyLevel_ANY
	CONSISTENCY_TWO                                     = cassandra.ConsistencyLevel_TWO
	CONSISTENCY_THREE                                   = cassandra.ConsistencyLevel_THREE
)

const (
	DEFAULT_SIZE              = 10
	DEFAULT_READ_CONSISTENCY  = CONSISTENCY_QUORUM
	DEFAULT_WRITE_CONSISTENCY = CONSISTENCY_QUORUM
	DEFAULT_TIMEOUT           = time.Second * 1
	DEFAULT_GRACE             = 5
	DEFAULT_RETRIES           = 5
)

const (
	LOWEST_COMPATIBLE_VERSION = 19
)

var (
	ErrorConnectionTimeout = errors.New("Connection timeout")
)

func (o *PoolOptions) defaults() {
	if o.Size == 0 {
		o.Size = DEFAULT_SIZE
	}
	if o.ReadConsistency == 0 {
		o.ReadConsistency = DEFAULT_READ_CONSISTENCY
	}
	if o.WriteConsistency == 0 {
		o.WriteConsistency = DEFAULT_WRITE_CONSISTENCY
	}
	if o.Timeout == 0 {
		o.Timeout = DEFAULT_TIMEOUT
	}
	if o.Grace == 0 {
		o.Grace = DEFAULT_GRACE
	}
	if o.Retries == 0 {
		o.Retries = DEFAULT_RETRIES
	}
}

type node struct {
	lastFailure int
	node        string
	available   lifo
}

type connectionPool struct {
	keyspace string
	options  PoolOptions
	schema   *Schema
	nodes    []*node
}

// NewConnectionPool creates a new connection pool for the given nodes and keyspace.
// nodes is in the format of "host:port" strings.
func NewConnectionPool(nodes []string, keyspace string, options PoolOptions) (ConnectionPool, error) {
	if len(nodes) <= 0 {
		return nil, errors.New("At least one node is required")
	}

	options.defaults()

	cp := &connectionPool{
		keyspace: keyspace,
		options:  options,
		nodes:    make([]*node, len(nodes)),
	}

	for i, n := range nodes {
		cp.nodes[i] = &node{node: n}
	}

	var ksDef *cassandra.KsDef
	err := cp.run(func(c *connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, error) {
		var ire *cassandra.InvalidRequestException
		var nfe *cassandra.NotFoundException
		var err error
		ksDef, nfe, ire, err = c.client.DescribeKeyspace(cp.keyspace)
		if nfe != nil {
			ksDef = nil
		}
		return ire, nil, nil, err
	})

	if err != nil {
		return nil, err
	}

	if ksDef == nil {
		return nil, errors.New("Keyspace not found while trying to parse schema")
	}

	cp.schema = newSchema(ksDef)
	if cp.schema == nil {
		return nil, errors.New("Cannot parse schema")
	}

	return cp, nil
}

type transaction func(*connection) (*cassandra.InvalidRequestException, *cassandra.UnavailableException, *cassandra.TimedOutException, error)

func (cp *connectionPool) run(t transaction) error {
	return cp.runWithRetries(t, cp.options.Retries)
}

func (cp *connectionPool) runWithRetries(t transaction, retries int) error {
	var c *connection
	var err error

	for tries := 0; tries < retries; tries++ {

		// acquire a new connection if we are just starting out or after discarding one
		if c == nil {
			c, err = cp.acquire()
			// nothing to do, cannot acquire a connection
			if err != nil {
				return err
			}
		}

		ire, ue, te, err := t(c)

		// nonrecoverable error, but not related to availability, do not retry and pass it to the user
		if ire != nil {
			cp.release(c)
			return errors.New(ire.Why)
		}

		// nonrecoverable error, drop the connection (but do not blacklist) and retry
		if err != nil {
			log.Printf("Node %s error %s", c.node, err.Error())
			c.close()
			c = nil
			continue
		}

		// the node is timing out. This Is Bad. move it to the blacklist and try again with another connection
		if te != nil {
			log.Printf("Node %s timed out, blacklisted", c.node)
			c.node.blacklist()
			c.close()
			c = nil
			continue
		}

		// one or more replicas are unavailable for the operation at the required consistency level. this is potentially
		// recoverable in a partitioned cluster by hoping to another connection/node and trying again
		if ue != nil {
			log.Printf("Unavailable exception")
			cp.release(c)
			c = nil
			continue
		}

		// no errors, release connection and return
		cp.release(c)
		return nil
	}

	// loop exited normally so it hit the retry limit
	return errors.New("Max retries hit trying to run a Cassandra transaction")
}

func (cp *connectionPool) randomNode(now int) (*node, error) {
	n := len(cp.nodes)
	i := rand.Int() % n

	for tries := 0; tries < n; tries++ {
		nodei := cp.nodes[i]
		if nodei.lastFailure+cp.options.Grace < now {
			return nodei, nil
		}
		i = (i + 1) % n
	}

	//TODO: try to acquire one anyway
	return nil, errors.New("All nodes are marked down, cannot acquire new connection")
}

func (cp *connectionPool) acquire() (*connection, error) {

	now := int(time.Now().Unix())
	n, err := cp.randomNode(now)
	if err != nil {
		return nil, err
	}
	c, ok := n.available.Pop()
	if ok {
		return c, nil
	}
	c, err = newConnection(n, cp.keyspace, cp.options.Timeout, cp.options.Authentication)
	if err == ErrorConnectionTimeout {
		n.blacklist()
	}
	return c, err
}

func (cp *connectionPool) release(c *connection) {
	c.node.available.Push(c)
}

func (n *node) blacklist() {
	n.lastFailure = int(time.Now().Unix())
}

func (cp *connectionPool) Reader() Reader {
	return newReader(cp, cp.options.ReadConsistency)
}

func (cp *connectionPool) Writer() Writer {
	return newWriter(cp, cp.options.WriteConsistency)
}

func (cp *connectionPool) Query(m Mapping) Query {
	return newQuery(cp, m)
}

func (cp *connectionPool) Batch() Batch {
	return newBatch(cp)
}

func (cp *connectionPool) Keyspace() string {
	return cp.keyspace
}

func (cp *connectionPool) Schema() *Schema {
	return cp.schema
}

func (cp *connectionPool) Close() {
}

type connection struct {
	socket    *thrift.TSocket
	transport *thrift.TFramedTransport
	client    cassandra.Cassandra
	node      *node
}

func newConnection(n *node, keyspace string, timeout time.Duration, authentication map[string]string) (*connection, error) {

	addr, err := net.ResolveTCPAddr("tcp", n.node)
	if err != nil {
		return nil, err
	}

	c := &connection{node: n}

	c.socket = thrift.NewTSocketFromAddrTimeout(addr, timeout)

	c.transport = thrift.NewTFramedTransport(c.socket)
	protocol := thrift.NewTBinaryProtocolTransport(c.transport)
	c.client = cassandra.NewCassandraClientProtocol(c.transport, protocol, protocol)

	if err = c.transport.Open(); err != nil {
		return nil, err
	}

	version, err := c.client.DescribeVersion()
	if err != nil {
		c.close()
		return nil, err
	}
	versionComponents := strings.Split(version, ".")
	if len(versionComponents) < 1 {
		return nil, errors.New(fmt.Sprint("Cannot parse the Thrift API version number: ", version))
	}
	majorVersion, err := strconv.Atoi(versionComponents[0])
	if err != nil {
		return nil, errors.New(fmt.Sprint("Cannot parse the Thrift API version number: ", version))
	}
	if majorVersion < LOWEST_COMPATIBLE_VERSION {
		return nil, errors.New(fmt.Sprint("Unsupported Thrift API version, lowest supported is ", LOWEST_COMPATIBLE_VERSION,
			", server reports ", majorVersion))
	}

	if len(authentication) > 0 {
		ar := cassandra.NewAuthenticationRequest()
		ar.Credentials = authentication
		autE, auzE, err := c.client.Login(ar)
		if autE != nil {
			return nil, errors.New("Login error: cannot authenticate with the given credentials")
		}
		if auzE != nil {
			return nil, errors.New("Login error: the given credentials are not authorized to access the server")
		}
		if err != nil {
			return nil, err
		}
	}

	ire, err := c.client.SetKeyspace(keyspace)
	if err != nil {
		c.close()
		return nil, err
	}
	if ire != nil {
		c.close()
		return nil, errors.New("Cannot set the keyspace")
	}

	return c, nil
}

func (c *connection) close() {
	c.transport.Close()
}
