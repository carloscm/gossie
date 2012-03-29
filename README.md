# About

Gossie is a Go library with a low level wrapper for the Cassandra 1.0 Thrift bindings with utilities for connection pooling, primitive type marshaling and easy query building. It also includes a higher level layer that allows mapping structs to Cassandra column famlilies, with support for advanced features like composite column names.

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

For application usage copy the sources to your GOPATH/src and issue a go install to build and copy the libraries:

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
pool := NewConnectionPool([]string{"localhost:9160"}, "Example", PoolOptions{Size: 50, Timeout: 3000})
````

The pool uses a simple randomized rule for connecting to the passed nodes, always keeping the total number of connections under PoolOptions.Size but without any guarantees on the number of connections per host. It has automatic failover and retry of operations.

### Low level queries

The Query and Mutation interfaces allow for low level queries to Cassandra and they follow the semantics of the native Thrift operations, but wrapped with much easier to use functions based on method chaining.

```Go
err = pool.Mutation().Insert("MyColumnFamily", row).Run()
row, err = pool.Query().Cf("MyColumnFamily").Get(id)
rows, err = pool.Query().Cf("MyColumnFamily").Where([]byte("MyIndexedColumn"), EQ, []byte("hi!")).IndexedGet(&IndexedRange{Count: 1000})
````

### Type marshaling

The low level interface is based on passing []byte values for everything, mirroring the Thrift API. For this reason the functions Marshal and Unmarshal provide for type conversion between native Go types and native Cassandra types.

### Struct maping

The first part of the high level Gossie interface is the Map/Unmap functions. These functions allow to convert Go structs into Row-s, and they have support of advanced features like composites or overriding column names and types.

```Go
/*
In CQL 3.0:
CREATE TABLE timeline (
    user_id varchar,
    tweet_id uuid,
    author varchar,
    body varchar,
    PRIMARY KEY (user_id, tweet_id)
);
*/

// In Gossie:
type Timeline struct {
	UserID  string  `cf:"Timeline" key:"UserID" col:"TweetID,*name" val:"*value"`
	TweetID UUID
	Author  string
	Body    string
}

row, err = Map(&Timeline{"userid", ..., "Author Name", "Hey this thing rocks!"})
````

### High level queries

Coming soon!


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
