package gossie

import (
	"errors"
	"fmt"
	"github.com/carloscm/gossie/src/cassandra"
	"github.com/pomack/thrift4go/lib/go/src/thrift"
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
	Close() error
}

// PoolOptions stores the options for the creation of a ConnectionPool
type PoolOptions struct {
	Size             int               // keep up to Size connections open and ready
	ReadConsistency  int               // default read consistency
	WriteConsistency int               // default write consistency
	Timeout          int               // socket timeout in ms
	Recycle          int               // close connections after Recycle seconds
	RecycleJitter    int               // max jitter to add to Recycle so not all connections close at the same time
	Grace            int               // if a node is blacklisted try to contact it again after Grace seconds
	Retries          int               // retry queries for Retries times before raising an error
	Authentication   map[string]string // if one or more keys are present, login() is called with the values from Authentication
}

const (
	CONSISTENCY_DEFAULT      = 0
	CONSISTENCY_ONE          = 1
	CONSISTENCY_QUORUM       = 2
	CONSISTENCY_LOCAL_QUORUM = 3
	CONSISTENCY_EACH_QUORUM  = 4
	CONSISTENCY_ALL          = 5
	CONSISTENCY_ANY          = 6
	CONSISTENCY_TWO          = 7
	CONSISTENCY_THREE        = 8
)

const (
	DEFAULT_SIZE              = 10
	DEFAULT_READ_CONSISTENCY  = CONSISTENCY_QUORUM
	DEFAULT_WRITE_CONSISTENCY = CONSISTENCY_QUORUM
	DEFAULT_TIMEOUT           = 1000
	DEFAULT_RECYCLE           = 60
	DEFAULT_RECYCLE_JITTER    = 10
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
	if o.Recycle == 0 {
		o.Recycle = DEFAULT_RECYCLE
	}
	if o.RecycleJitter == 0 {
		o.RecycleJitter = DEFAULT_RECYCLE_JITTER
	}
	if o.Grace == 0 {
		o.Grace = DEFAULT_GRACE
	}
	if o.Retries == 0 {
		o.Retries = DEFAULT_RETRIES
	}
}

type nodeInfo struct {
	lastFailure int
	node        string
}

type slot struct {
	conn      *connection
	lastUsage int
}

type connectionPool struct {
	keyspace  string
	options   PoolOptions
	schema    *Schema
	nodes     []*nodeInfo
	available chan *slot
}

// NewConnectionPool creates a new connection pool for the given nodes and keyspace.
// nodes is in the format of "host:port" strings.
func NewConnectionPool(nodes []string, keyspace string, options PoolOptions) (ConnectionPool, error) {
	if len(nodes) <= 0 {
		return nil, errors.New("At least one node is required")
	}

	options.defaults()

	cp := &connectionPool{
		keyspace:  keyspace,
		options:   options,
		nodes:     make([]*nodeInfo, len(nodes)),
		available: make(chan *slot, options.Size),
	}

	for i, n := range nodes {
		cp.nodes[i] = &nodeInfo{node: n}
	}

	for i := 0; i < options.Size; i++ {
		cp.available <- &slot{}
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
			return errors.New(ire.String())
		}

		// nonrecoverable error, but not related to availability, do not retry and pass it to the user
		if err != nil {
			cp.release(c)
			return err
		}

		// the node is timing out. This Is Bad. move it to the blacklist and try again with another connection
		if te != nil {
			cp.blacklist(c.node)
			c.close()
			c = nil
			continue
		}

		// one or more replicas are unavailable for the operation at the required consistency level. this is potentially
		// recoverable in a partitioned cluster by hoping to another connection/node and trying again
		if ue != nil {
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

func (cp *connectionPool) randomNode(now int) (string, error) {
	n := len(cp.nodes)
	i := rand.Int() % n
	var node string

	for tries := 0; tries < n; tries++ {
		nodei := cp.nodes[i]
		if nodei.lastFailure+cp.options.Grace < now {
			return nodei.node, nil
		}
		i = (i + 1) % n
	}

	return node, errors.New("All nodes are marked down, cannot acquire new connection")
}

func (cp *connectionPool) acquire() (*connection, error) {
	var c *connection

	s := <-cp.available

	now := int(time.Now().Unix())
	if s.lastUsage+cp.options.Recycle+(rand.Int()%cp.options.RecycleJitter) < now {
		if s.conn != nil {
			s.conn.close()
		}
		s.conn = nil
	}

	if s.conn == nil {
		node, err := cp.randomNode(now)
		if err != nil {
			cp.releaseEmpty()
			return nil, err
		}
		c, err = newConnection(node, cp.keyspace, cp.options.Timeout, cp.options.Authentication)
		if err == ErrorConnectionTimeout {
			cp.blacklist(node)
			return nil, err
		}
		if err != nil {
			cp.releaseEmpty()
			return nil, err
		}
	} else {
		c = s.conn
	}

	return c, nil
}

func (cp *connectionPool) release(c *connection) {
	cp.available <- &slot{conn: c, lastUsage: int(time.Now().Unix())}
}

func (cp *connectionPool) releaseEmpty() {
	cp.available <- &slot{}
}

func (cp *connectionPool) blacklist(badNode string) {
	n := len(cp.nodes)
	for i := 0; i < n; i++ {
		node := cp.nodes[i]
		if node.node == badNode {
			node.lastFailure = int(time.Now().Unix())
			break
		}
	}
	cp.releaseEmpty()
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

func (cp *connectionPool) Close() error {
	return nil
}

type connection struct {
	socket    *thrift.TNonblockingSocket
	transport *thrift.TFramedTransport
	client    *cassandra.CassandraClient
	node      string
	keyspace  string
}

func newConnection(node, keyspace string, timeout int, authentication map[string]string) (*connection, error) {

	addr, err := net.ResolveTCPAddr("tcp", node)
	if err != nil {
		return nil, err
	}

	c := &connection{node: node}

	c.socket, err = thrift.NewTNonblockingSocketAddr(addr)
	if err != nil {
		return nil, err
	}

	// socket not open yet, so no error expected. it expects nanos, we have milis, so it's 1e6
	c.socket.SetTimeout(int64(timeout) * 1e6)

	c.transport = thrift.NewTFramedTransport(c.socket)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	c.client = cassandra.NewCassandraClientFactory(c.transport, protocolFactory)

	// simulate timeout support for the underlying Dial() in .Open(). needless to say this sucks
	// restore sanity to this for Go v1 with the new DialTimeout() func
	ch := make(chan bool, 1)
	go func() {
		err = c.transport.Open()
		ch <- true
	}()
	timedOut := false
	select {
	case <-time.After(time.Duration(timeout) * time.Millisecond):
		timedOut = true
	case <-ch:
	}
	if timedOut {
		return nil, ErrorConnectionTimeout
	}
	if err != nil {
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
		ar.Credentials = thrift.NewTMap(thrift.STRING, thrift.STRING, 1)
		for k, v := range authentication {
			ar.Credentials.Set(k, v)
		}
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

	c.keyspace = keyspace

	return c, nil
}

func (c *connection) close() {
	c.transport.Close()
}
