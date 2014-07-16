// Experimental fork of carloscm/gossie
package gossie

import (
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/golang/glog"
	"github.com/wadey/gossie/src/cassandra"
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
	Query(mapping Mapping) Query

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
	Size             int                        // keep up to Size connections PER NODE open and ready
	ReadConsistency  cassandra.ConsistencyLevel // default read consistency
	WriteConsistency cassandra.ConsistencyLevel // default write consistency
	Timeout          time.Duration              // socket timeout
	BleederInterval  time.Duration              // <del>kill a kitten</del> Close a least used connection every BleederInterval.
	Grace            int                        // if a node is blacklisted try to contact it again after Grace seconds
	Retries          int                        // retry queries for Retries times before raising an error
	Authentication   map[string]string          // if one or more keys are present, login() is called with the values from Authentication
}

var DefaultPoolOptions = PoolOptions{
	Size:             10,
	ReadConsistency:  CONSISTENCY_QUORUM,
	WriteConsistency: CONSISTENCY_QUORUM,
	Timeout:          time.Second * 1,
	BleederInterval:  time.Second * 2,
	Grace:            5,
	Retries:          5,
	// Authentication   is empty
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
	LOWEST_COMPATIBLE_VERSION = 19
)

var (
	ErrorConnectionTimeout = errors.New("Connection timeout")
)

func (o *PoolOptions) mergeFrom(r *PoolOptions) {
	if r.Size != 0 {
		o.Size = r.Size
	}
	if r.ReadConsistency != 0 {
		o.ReadConsistency = r.ReadConsistency
	}
	if r.WriteConsistency != 0 {
		o.WriteConsistency = r.WriteConsistency
	}
	if r.Timeout != 0 {
		o.Timeout = r.Timeout
		// TODO temporary backwards compatibility fix
		if o.Timeout < 1*time.Millisecond {
			o.Timeout = o.Timeout * time.Millisecond
		}
	}
	if r.BleederInterval != 0 {
		o.BleederInterval = r.BleederInterval
	}
	if r.Grace != 0 {
		o.Grace = r.Grace
	}
	if r.Retries != 0 {
		o.Retries = r.Retries
	}
	if r.Authentication != nil {
		o.Authentication = r.Authentication
	}
}

type node struct {
	lastFailure int
	node        string
	available   lifo
}

type connectionPool struct {
	keyspace  string
	options   PoolOptions
	schema    *Schema
	nodes     []*node
	tlsConfig *tls.Config
}

var nowfunc func() time.Time = time.Now

// NewConnectionPool creates a new connection pool for the given nodes and keyspace.
// nodes is in the format of "host:port" strings.
func NewConnectionPool(nodes []string, keyspace string, options PoolOptions, tlsConfig *tls.Config) (ConnectionPool, error) {
	if len(nodes) <= 0 {
		return nil, errors.New("At least one node is required")
	}

	cp := &connectionPool{
		keyspace:  keyspace,
		options:   DefaultPoolOptions,
		nodes:     make([]*node, len(nodes)),
		tlsConfig: tlsConfig,
	}
	cp.options.mergeFrom(&options)

	for i, n := range nodes {
		cp.nodes[i] = &node{node: n}
	}

	var ksDef *cassandra.KsDef
	err := cp.run(func(c *connection) error {
		var err error
		ksDef, err = c.client.DescribeKeyspace(cp.keyspace)
		return err
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
	go cp.bleeder(options.BleederInterval)

	return cp, nil
}

func (cp *connectionPool) bleeder(d time.Duration) {
	l := len(cp.nodes)
	c := time.Tick(d)
	nodeidx := -1
	for _ = range c {
		for i := 0; i < l; i++ {
			nodeidx++
			if c, ok := cp.nodes[nodeidx%l].available.PopBottom(cp.options.Size); ok {
				glog.V(1).Info("Closing connection to ", cp.nodes[nodeidx%l].node)
				c.close()
				break
			}
		}
	}
}

type transaction func(*connection) error

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
				glog.Error("Unable to acquire cassandra connection: ", err)
				return err
			}
		}

		err := t(c)

		if err != nil {
			switch err.(type) {
			case *cassandra.InvalidRequestException:
				// nonrecoverable error, but not related to availability, do not retry and pass it to the user
				glog.Errorf("Node %s Invalid request: %s", c.node.node, err.(*cassandra.InvalidRequestException).Why)
				cp.release(c)
				return err
			case *cassandra.TimedOutException:
				// the node is timing out. This Is Bad. move it to the blacklist and try again with another connection
				glog.Infof("Node %s %s, blacklisted", c.node.node, err)
				c.node.blacklist()
				c.close()
				c = nil
				continue
			case *cassandra.UnavailableException:
				// one or more replicas are unavailable for the operation at the required consistency level. this is potentially
				// recoverable in a partitioned cluster by hoping to another connection/node and trying again
				glog.Info("Node %s %s", c.node.node, err)
				cp.release(c)
				c = nil
				continue
			default:
				// nonrecoverable error, drop the connection (but do not blacklist) and retry
				glog.Errorf("Node %s error %s", c.node.node, err)
				c.close()
				c = nil
				continue
			}
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

	now := int(nowfunc().Unix())
	n, err := cp.randomNode(now)
	if err != nil {
		return nil, err
	}
	c, ok := n.available.Pop()
	if ok {
		return c, nil
	}
	c, err = newConnection(n, cp.keyspace, cp.options.Timeout, cp.options.Authentication, cp.tlsConfig)
	if err == ErrorConnectionTimeout {
		n.blacklist()
	}
	return c, err
}

func (cp *connectionPool) release(c *connection) {
	c.node.available.Push(c)
}

func (n *node) blacklist() {
	n.lastFailure = int(nowfunc().Unix())
	//close all connections
	glog.V(1).Info("closing %d connections to blacklisted node %s", len(n.available.l), n.node)
	for c, ok := n.available.Pop(); ok; c, ok = n.available.Pop() {
		c.close()
	}
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
	sslSocket *thrift.TSSLSocket
	transport *thrift.TFramedTransport
	client    cassandra.Cassandra
	node      *node
}

func newConnection(n *node, keyspace string, timeout time.Duration, authentication map[string]string, tlsConfig *tls.Config) (*connection, error) {

	addr, err := net.ResolveTCPAddr("tcp", n.node)
	if err != nil {
		return nil, err
	}

	c := &connection{node: n}
	if tlsConfig == nil {
		c.socket = thrift.NewTSocketFromAddrTimeout(addr, timeout)
		c.transport = thrift.NewTFramedTransport(c.socket)
	} else {
		c.sslSocket = thrift.NewTSSLSocketFromAddrTimeout(addr, tlsConfig, timeout)
		c.transport = thrift.NewTFramedTransport(c.sslSocket)
	}

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
		err := c.client.Login(ar)
		if err != nil {
			c.close()
			switch err.(type) {
			case *cassandra.AuthenticationException:
				return nil, errors.New("Login error: cannot authenticate with the given credentials")
			case *cassandra.AuthorizationException:
				return nil, errors.New("Login error: the given credentials are not authorized to access the server")
			default:
				return nil, err
			}
		}
	}

	err = c.client.SetKeyspace(keyspace)
	if err != nil {
		c.close()
		switch err.(type) {
		case *cassandra.InvalidRequestException:
			err = errors.New("Cannot set the keyspace " + keyspace)
		}
		return nil, err
	}

	return c, nil
}

func (c *connection) close() {
	c.transport.Close()
}
