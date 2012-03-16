# About

Gossie is (for now) a Go library with a low level wrapper for the Cassandra 1.0 thrift bidings with with utilities for connection pooling, primitive type marshaling and easy query building (much easier to use than the generated thrift bindings). A higher level layer will be implemented on top of the current code to allow for struct marshalling into rows and composites, among other things.


# Requeriments

The official Apache Thrift libraries for Go are outdated and buggy. For now the active development happens in thrift4go:
https://github.com/pomack/thrift4go

Install from the tip of master, make inside thrift4go/lib/go for installing the Thrift Go library.

Once Go 1.0 is released the author will submit a new version of the Go Thrift libraries/generator to Apache.


# Installing

There is no need to generate a Cassandra biding, I am providing one with Gossie (and the whole point is not to have to use it!)

I am using godag, a "go command"-like wrapper for compiling/linking/etc so there is no Makefile in Gossie. I do not plan on providing one, and Go 1.0 would make it obsolete anyway.

Gossie is written in Go r60.3 for now. I am waiting for the official Go 1.0 release to port it.


# Running the tests

Launch a Cassandra instance in localhost:9160, create a keyspace named TestGossie, and execute the provided schema-test.txt to create the test column families. Now you can run the Gossie tests.


# Example

I will provide a full example once the higher level marshaling is implemented. For examples of the low level layer check src/gossie/query_test.go


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
