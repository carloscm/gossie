package gossie

import (
	"flag"
	"log"
	"os/exec"
	"testing"
	"time"
)

var cassandraRestart = flag.Bool("cassandra-restart", false, "run tests to ensure ConnectionPool tolerates Cassandra restarting")

func TestCassandraRestart(t *testing.T) {
	if !*cassandraRestart {
		t.Skipf("use -cassandra-restart to ensure ConnectionPool tolerates Cassandra restarting")
	}
	rPool, err := NewConnectionPool(
		localEndpointPool,
		keyspace,
		PoolOptions{Size: 1, Timeout: standardTimeout},
	)
	if err != nil {
		t.Fatal(err)
	}
	wPool, err := NewConnectionPool(
		localEndpointPool,
		keyspace,
		PoolOptions{Size: 1, Timeout: standardTimeout},
	)
	if err != nil {
		t.Fatal(err)
	}
	row := &Row{
		Key: []byte("test"),
		Columns: []*Column{
			&Column{
				Name:      []byte("test"),
				Value:     []byte("test"),
				Timestamp: time.Now().UnixNano(),
			},
		},
	}
	if _, err := rPool.Reader().Cf("AllTypes").Get([]byte("test")); err != nil {
		t.Fatal(err)
	}
	if err := wPool.Writer().Insert("AllTypes", row).Run(); err != nil {
		t.Fatal(err)
	}
	var argv []string
	if false {
		argv = []string{"sh", "-c", "launchctl unload /Library/LaunchDaemons/org.apache.cassandra.plist; sudo launchctl load /Library/LaunchDaemons/org.apache.cassandra.plist"}
	} else {
		argv = []string{"/etc/init.d/cassandra", "restart"}
	}
	log.Println("restarting Cassandra...")
	if err := exec.Command("sudo", argv...).Run(); err != nil {
		t.Fatal(err)
	}
	log.Println("restarted Cassandra; sleeping...")
	time.Sleep(10 * time.Second)
	if _, err := rPool.Reader().Cf("AllTypes").Get([]byte("test")); err != nil {
		t.Fatal(err)
	}
	if err := wPool.Writer().Insert("AllTypes", row).Run(); err != nil {
		t.Fatal(err)
	}
}
