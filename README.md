# This is a fork

This is a fork of the original Gossie library created by carloscm.

Some of the changes / additions in this fork:

- the `mockgossie` package, which implements an in-memory store for your tests.
- fix singular compact column names: https://github.com/wadey/gossie/commit/4518cec59bf4ecd13323e41c3d5b8ddc289a5f04

# About

Gossie is a Go library for Apache Cassandra. It includes a wrapper for the Cassandra 1.0 Thrift bindings with utilities for connection pooling, primitive type marshaling and easy query building. It also includes a higher level layer that allows mapping structs to Cassandra column famlilies, with support for advanced features like composite column names.


# Requirements

The official Apache Thrift libraries for Go are outdated and buggy. For now the active development happens in thrift4go:
https://github.com/pomack/thrift4go

Install thrift4go:

```
go get "github.com/pomack/thrift4go/lib/go/src/thrift"
```


# Installing

There is no need to generate a Cassandra Thrift binding, I am providing one with Gossie (and the whole point is not to have to use it!).

For application usage issue two `go get` commands, one for the bindings and another for Gossie

```
go get "github.com/carloscm/gossie/src/cassandra"
go get "github.com/carloscm/gossie/src/gossie"
```

If you want to fork and do development on Gossie itself the main command you need to run is something like (from the root of the Gossie folder):

```
GOPATH=$GOPATH:`pwd` go test gossie
```


# Running the tests

Launch a Cassandra instance in localhost:9160 and execute the provided schema-test.txt with cassandra-cli to create the test keyspace and column families. Now you can run the Gossie tests.


# Quickstart

### Import

Gossie follows the Go 1.0 packaging conventions. Import Gossie into your code like this:

```Go
import (
	"github.com/carloscm/gossie/src/gossie"
)
````

### Connection pooling

To create a connection use the method NewConnectionPool, passing a list of nodes, the desired keyspace, and a PoolOptions with the various connection options you can tune.

```Go
pool, err := gossie.NewConnectionPool([]string{"localhost:9160"}, "Example", gossie.PoolOptions{Size: 50, Timeout: 3000})
if err != nil {
	// do something
}
````

The pool uses a simple randomized rule for connecting to the passed nodes, always keeping the total number of connections under PoolOptions.Size but without any guarantees on the number of connections per host. It has automatic failover and retry of operations.

### Low level queries

The Reader and Writer interfaces allow for low level queries to Cassandra and they follow the semantics of the native Thrift operations, but wrapped with much easier to use functions based on method chaining.

```Go
err = pool.Writer().Insert("MyColumnFamily", row).Run()
row, err = pool.Reader().Cf("MyColumnFamily").Get(id)
rows, err = pool.Reader().Cf("MyColumnFamily").Where([]byte("MyIndexedColumn"), gossie.EQ, []byte("hi!")).IndexedGet(&gossie.IndexedRange{Count: 1000})
````

### Type marshaling

The low level interface is based on passing []byte values for everything, mirroring the Thrift API. For this reason the functions Marshal and Unmarshal provide for type conversion between native Go types and native Cassandra types.

### Struct mapping

The Mapping interface and its implementations allow to convert Go structs into Rows, and they have support of advanced features like composites or overriding column names and types. Built-in NewMapping() returns a Mapping implementation that can map and unmap Go structs from Cassandra rows, serialized in classic key/value rows or in composited column names, with support for both sparse and compact storage. For example:

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
	UserID  string `cf:"Timeline" key:"UserID" cols:"TweetID"`
	TweetID int64
	Author  string
	Body    string
}

mapping := gossie.MustNewMapping(&Tweet{})
row, err = mapping.Map(&Tweet{"userid", 10000000000004, "Author Name", "Hey this thing rocks!"})
err = pool.Writer().Insert("Timeline", row).Run()
````

When calling [Must]NewMapping() you can tag your struct fiels with `name`, `type` and `skip`. The `name` field tag will change the column name to its value when the field it appears on is (un)marhsaled to/from a Cassandra row column. The `type` field tag allows to override the default type Go<->Cassandra type mapping used by Gossie for the field it appears on. If `skip:"true"` is present the field will be ignored by Gossie.

The tags `mapping`, `cf`, `key`, `cols` and `value` can be used in any field in the struct to document a mapping. `mapping` is optional and can have a value of `sparse` (the default) or `compact`. See [CQL3.0](http://www.datastax.com/dev/blog/whats-new-in-cql-3-0) for more information. `cf` is the column family name. `key` is the field name in the struct that stores the Cassandra row key value. `cols` is optional and it is a list of struct fiels that build up the composite column name, if there is any. `value` is the field that stores the column value for compact storage rows, and it is ignored in sparse storage rows.

Mapping instances are reusable and you are encouraged to cache them.

### Query and Result

Query allows to look up mapped structs over Cassandra rows. Pass to `Query.Components` one or more component values that all the result objects must have in common. You can also leave out the last component and use `Query.Between` to slice a range of values for it. Call `Query.Get` with the row key to get a Result. `Result.Next` reads a single struct from the Cassandra row, and returns `Done` when no more structs can be read.

```Go
query := pool.Query(TweetMapping)

// a single tweet, since we pass the row key and all possible composite values
result, err := query.Components(10000000000004).Get("username")

// all tweets for a given user
result, err := query.Get("username")

// iterating over results
for {
	t := &Tweet{}
	err := result.Next(t)
	if err != nil { // Done is also returned in err
		break
	}
}
````

### Batch

Batch is a thin interface over `Writer` which allows to directly write and delete structs in a higher level fashion. Its use is simple, for example:

```Go
mapping, err := gossie.NewMapping(&Tweet{})
batch := pool.Batch()
tweet := &Tweet{"userid", 10000000000004, "Author Name", "Hey this thing rocks!"}
err = batch.Insert(mapping, tweet).Run()
````

Use a new `Batch()` call for every batch of writes you want to perform. Its internal state may keep copies of your data so it is not reusable.


# Planned features

- Query: secondary index read with buffering
- High level mapping for Go slices
- High level mapping for Go maps


# Not planned

- Supercolumns
- Dynamic composite comparator
- CQL3. See https://github.com/tux21b/gocql


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
