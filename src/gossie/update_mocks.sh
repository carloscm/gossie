#! /bin/bash -e

go get github.com/golang/mock/mockgen
mockgen github.com/apesternikov/gossie/src/cassandra Cassandra >mock_cassandra/mock_cassandra.go
#mockgen -package="gossie" github.com/apesternikov/gossie/src/gossie ConnectionPool >/tmp/mock_gossie.go
#mv /tmp/mock_gossie.go .
gofmt -w mock_cassandra/mock_cassandra.go

echo >&2 "OK"
