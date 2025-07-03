package config

import (
	"be-realtime-chat-app/services/commoner/utils"
	"fmt"
)

type ServerConfig struct {
	QueryHTTPAddr string
	QueryHTTPPort string
	QuerySvcName  string

	QueryGRPCAddr         string
	QueryGRPCPort         string
	QueryGRPCInternalAddr string

	ConsulAddr string
}

func NewServerConfig() *ServerConfig {
	// hostname, _ := os.Hostname()
	// addrs, _ := net.LookupIP(hostname)

	// var ip string
	// for _, addr := range addrs {
	// 	if ipv4 := addr.To4(); ipv4 != nil && !ipv4.IsLoopback() {
	// 		ip = ipv4.String()
	// 		break
	// 	}
	// }

	consulAddr := fmt.Sprintf("%s:%s", utils.GetEnv("CONSUL_HOST"), utils.GetEnv("CONSUL_PORT"))

	return &ServerConfig{
		QueryHTTPAddr: utils.GetEnv("QUERY_HTTP_ADDR"),
		QueryHTTPPort: utils.GetEnv("QUERY_HTTP_PORT"),
		QuerySvcName:  utils.GetEnv("QUERY_SVC_NAME"),

		QueryGRPCAddr:         utils.GetEnv("QUERY_GRPC_ADDR"),
		QueryGRPCPort:         utils.GetEnv("QUERY_GRPC_PORT"),
		QueryGRPCInternalAddr: utils.GetEnv("QUERY_GRPC_ADDR"),

		ConsulAddr: consulAddr,
	}
}
