package config

import (
	"be-realtime-chat-app/services/commoner/utils"
	"fmt"
)

type ServerConfig struct {
	RoomHTTPAddr string
	RoomHTTPPort string
	RoomSvcName  string

	RoomGRPCAddr         string
	RoomGRPCPort         string
	RoomGRPCInternalAddr string

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
		RoomHTTPAddr: utils.GetEnv("ROOM_HTTP_ADDR"),
		RoomHTTPPort: utils.GetEnv("ROOM_HTTP_PORT"),
		RoomSvcName:  utils.GetEnv("ROOM_SVC_NAME"),

		RoomGRPCAddr:         utils.GetEnv("ROOM_GRPC_ADDR"),
		RoomGRPCPort:         utils.GetEnv("ROOM_GRPC_PORT"),
		RoomGRPCInternalAddr: utils.GetEnv("ROOM_GRPC_ADDR"),

		ConsulAddr: consulAddr,
	}
}
