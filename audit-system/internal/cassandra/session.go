package cassandra

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gocql/gocql"
)

type Config struct {
	Hosts       []string
	Keyspace    string
	Consistency gocql.Consistency
	Timeout     time.Duration
	NumConns    int
}

func NewSession(config Config) (*gocql.Session, error) {
	cluster := gocql.NewCluster(config.Hosts...)
	cluster.Keyspace = config.Keyspace
	cluster.Consistency = config.Consistency
	cluster.Timeout = config.Timeout
	cluster.NumConns = config.NumConns
	/*
		Policy defines how Go drive picks which Cassandra node to send query to in multi nodes cluster
		Inner (fallback layer) vs outer(primary layer)
	*/

	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create cassandra session: %w", err)
	}
	slog.Info("cassandra session established", "keyspace", config.Keyspace)
	return session, nil
}
