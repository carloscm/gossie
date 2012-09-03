package gossie

var (
	invalidEndpoint    = "localhost:9999"
	localEndpoint      = "localhost:9160"
	localEndpointPool  = []string{localEndpoint}
	localEndpointsPool = []string{localEndpoint, "localhost:9170", "localhost:9180"}

	keyspace = "TestGossie"

	standardTimeout = 3000
	shortTimeout    = 1000

	poolOptions = PoolOptions{Size: 50, Timeout: standardTimeout}
)
