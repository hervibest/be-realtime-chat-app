package config

import (
	"be-realtime-chat-app/services/commoner/logs"
	"be-realtime-chat-app/services/commoner/utils"
	"fmt"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func NewNATsConn(log logs.Log) *nats.Conn {
	host := utils.GetEnv("NATS_HOST")
	port := utils.GetEnv("NATS_PORT")
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%s", host, port))
	if err != nil {
		log.Panic("Failed to connect to NATS", zap.Error(err))
	}
	log.Info("Successfully connected to nats jetstream")
	return nc
}
