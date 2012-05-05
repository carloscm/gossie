# About

Gossie is a Go library for Apache Cassandra. It includes a wrapper for the Cassandra 1.0 Thrift bindings with utilities for connection pooling, primitive type marshaling and easy query building. It also includes a higher level layer that allows mapping structs to Cassandra column famlilies, with support for advanced features like composite column names.


# Requeriments

The official Apache Thrift libraries for Go are outdated and buggy. For now the active development happens in thrift4go:
https://github.com/pomack/thrift4go

Installing thrift4go under GOPATH in Go 1:

```
1) cd lib/go/src
2) cp -R thrift $GOPATH/src
3) go install thrift
```


# Installing

There is no need to generate a Cassandra Thrift biding, I am providing one with Gossie (and the whole point is not to have to use it!)

For application usage copy the sources to your $GOPATH/src and issue a go install to build and copy the libraries:

```
1) cp -R src/* $GOPATH/src
2) go install cassandra gossie
```

If you want to fork and do development on Gossie itself the main command you need to run is something like (from the root of the Gossie folder):

```
# locally install inside a pkg folder the depedencies, ie the cassandra bindings
GOPATH=$GOPATH:`pwd` go test -i gossie
# actually build and run the gossie tests
GOPATH=$GOPATH:`pwd` go test gossie
```


# Running the tests

Launch a Cassandra instance in localhost:9160, create a keyspace named TestGossie, and execute the provided schema-test.txt to create the test column families. Now you can run the Gossie tests.


# Quickstart

### Connection pooling

To create a connection use the method NewConnectionPool, passing a list of nodes, the desired keyspace, and a PoolOptions with the various connection options you can tune.

```Go
pool := gossie.NewConnectionPool([]string{"localhost:9160"}, "Example", PoolOptions{Size: 50, Timeout: 3000})
````

The pool uses a simple randomized rule for connecting to the passed nodes, always keeping the total number of connections under PoolOptions.Size but without any guarantees on the number of connections per host. It has automatic failover and retry of operations.

### Low level queries

The Reader and Writer interfaces allow for low level queries to Cassandra and they follow the semantics of the native Thrift operations, but wrapped with much easier to use functions based on method chaining.

```Go
err = pool.Writer().Insert("MyColumnFamily", row).Run()
row, err = pool.Reader().Cf("MyColumnFamily").Get(id)
rows, err = pool.Reader().Cf("MyColumnFamily").Where([]byte("MyIndexedColumn"), EQ, []byte("hi!")).IndexedGet(&IndexedRange{Count: 1000})
````

### Type marshaling

The low level interface is based on passing []byte values for everything, mirroring the Thrift API. For this reason the functions Marshal and Unmarshal provide for type conversion between native Go types and native Cassandra types.

### Struct mapping

The Mapping interface and its implementations allow to convert Go structs into Rows, and they have support of advanced features like composites or overriding column names and types. NewSparse() returns a Mapping for the new CQL 3.0 pattern of composite "primary keys":

```Go
/*
In CQL 3.0:
CREATE TABLE Timeline (
    UserID varchar,
    TweetID bigint,
    Author varchar,
    Body varchar,
    PRIMARY KEY (UserID, TweetID)
);
*/

// In Gossie:
type Tweet struct {
	UserID  string
	TweetID int64
	Author  string
	Body    string
}

mapping := gossie.NewSparse("Timeline", "UserID", "TweetID")
row, err = mapping.Map(&Tweet{"userid", 10000000000004, "Author Name", "Hey this thing rocks!"})
err = pool.Writer().Insert("Timeline", row).Run()
````

When calling Mapping.Map() you can tag your struct fiels with `name`, `type` and `skip`. The `name` field tag will change the column name to its value when the field it appears on is (un)marhsaled to/from a Cassandra row column. The `type` field tag allows to override the default type Go<->Cassandra type mapping used by Gossie for the field it appears on. If `skip:"true"` is present the field will be ignored by Gossie.

The tags `cf`, `key` and `cols` can be used in any field in the struct to document a mapping. It can later be extracted with `MappingFromTags()` by passing any instance of the struct, even an empty one. For example this is equivalent to the mapping created  with `NewSparse()` in the previous example:

```Go
type Tweet struct {
	UserID  string `cf:"Timeline" key:"UserID" cols:"TweetID"`
	TweetID int64
	Author  string
	Body    string
}
mapping := gossie.MappingFromTags(&Tweet{})
```

### Query and Result interfaces (planned)

High level queries with transparent paging and buffering. This is still WIP, a possible example:

```Go

query := pool.Query(TweetMapping)

// a single tweet, since we pass the row key and all possible composite values
result, err := query.Get("username", 10000000000004)

// all tweets for a given user
result, err := query.Get("username")

// all tweets for a given user, starting at a certain TweetID
result, err := query.Where("TweetID", ">=", 10000000000004).Get("username")

// iterating over results
for {
	var t Tweet
	err := result.Next(&t)
	if err != nil {
		break
	}
	...
}
````

# Planned features

- Query: range reads for composites with buffering and paging
- Query: secondary index read with buffering and paging
- Query: multiget reads with buffering and paging
- A higher level abstraction for writes (Batch interface)
- High level mapping for Go slices
- High level mapping for Go maps


# License

Copyright (C) 2012 by Carlos Carrasco

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
