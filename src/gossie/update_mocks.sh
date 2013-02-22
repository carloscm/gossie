#! /bin/bash -e

mockgen github.com/apesternikov/gossie/src/cassandra ICassandra >mock_cassandra/mock_cassandra.go
#mockgen -package="gossie" github.com/apesternikov/gossie/src/gossie ConnectionPool >/tmp/mock_gossie.go
#mv /tmp/mock_gossie.go .
gofmt -w mock_cassandra/mock_cassandra.go

echo >&2 "OK"
