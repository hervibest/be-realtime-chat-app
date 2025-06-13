package adapter

import (
	"be-realtime-chat-app/proto/roompb"
	"be-realtime-chat-app/services/commoner/discovery"
	"be-realtime-chat-app/services/commoner/helper"
	"be-realtime-chat-app/services/commoner/logs"
	"be-realtime-chat-app/services/commoner/utils"
	"context"
	"log"
)

type RoomAdapter interface {
	GetRoom(ctx context.Context, token string) (*roompb.GetRoomResponse, error)
}

type roomAdapter struct {
	client roompb.RoomServiceClient
}

func NewRoomAdapter(ctx context.Context, registry discovery.Registry, logs logs.Log) (RoomAdapter, error) {
	roomServiceName := utils.GetEnv("ROOM_SVC_NAME") + "-grpc"
	conn, err := discovery.ServiceConnection(ctx, roomServiceName, registry, logs)
	if err != nil {
		return nil, err
	}

	log.Print("successfuly connected to room-svc-grpc")
	client := roompb.NewRoomServiceClient(conn)

	return &roomAdapter{
		client: client,
	}, nil
}

func (a *roomAdapter) GetRoom(ctx context.Context, roomID string) (*roompb.GetRoomResponse, error) {
	processPhotoRequest := &roompb.GetRoomRequest{
		RoomId: roomID,
	}
	response, err := a.client.GetRoom(ctx, processPhotoRequest)
	if err != nil {
		return nil, helper.FromGRPCError(err)
	}

	return response, nil
}
