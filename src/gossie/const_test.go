package gossie

import (
	"time"
)

var (
	invalidEndpoint    = "localhost:9999"
	localEndpoint      = "localhost:9160"
	localEndpointPool  = []string{localEndpoint}
	localEndpointsPool = []string{localEndpoint, "localhost:9170", "localhost:9180"}

	keyspace = "TestGossie"

	standardTimeout = time.Second * 3
	shortTimeout    = time.Second * 1

	poolOptions = PoolOptions{Size: 50, Timeout: standardTimeout}
)
