package config

import (
	"be-realtime-chat-app/services/commoner/utils"
	"fmt"
)

type ServerConfig struct {
	ChatHTTPAddr string
	ChatHTTPPort string
	ChatSvcName  string

	ConsulAddr string
}

func NewServerConfig() *ServerConfig {
	consulAddr := fmt.Sprintf("%s:%s", utils.GetEnv("CONSUL_HOST"), utils.GetEnv("CONSUL_PORT"))

	return &ServerConfig{
		ChatHTTPAddr: utils.GetEnv("CHAT_HTTP_ADDR"),
		ChatHTTPPort: utils.GetEnv("CHAT_HTTP_PORT"),
		ChatSvcName:  utils.GetEnv("CHAT_SVC_NAME"),

		ConsulAddr: consulAddr,
	}
}
