package config

import (
	"be-realtime-chat-app/services/commoner/utils"
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

const (
	keyspace = "messaging_service"
)

func NewCQLDB() (*gocql.Session, error) {
	hosts := []string{
		utils.GetEnv("SCYLLA_HOST_1"),
		utils.GetEnv("SCYLLA_HOST_2"),
		utils.GetEnv("SCYLLA_HOST_3"),
	}

	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 5 * time.Second
	cluster.ConnectTimeout = 5 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	log.Println("Connected to ScyllaDB cluster")
	return session, nil
}
