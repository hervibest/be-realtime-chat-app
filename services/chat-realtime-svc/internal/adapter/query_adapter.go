package adapter

import (
	"be-realtime-chat-app/proto/querypb"
	"be-realtime-chat-app/services/commoner/discovery"
	"be-realtime-chat-app/services/commoner/helper"
	"be-realtime-chat-app/services/commoner/logs"
	"be-realtime-chat-app/services/commoner/utils"
	"context"
	"log"
)

type QueryAdapter interface {
	GetTenLatestMessage(ctx context.Context, roomID string) (*querypb.GetTenLatestMessageResponse, error)
}

type queryAdapter struct {
	client querypb.QueryServiceClient
}

func NewQueryAdapter(ctx context.Context, registry discovery.Registry, logs logs.Log) (QueryAdapter, error) {
	roomServiceName := utils.GetEnv("QUERY_SVC_NAME") + "-grpc"
	conn, err := discovery.ServiceConnection(ctx, roomServiceName, registry, logs)
	if err != nil {
		return nil, err
	}

	log.Print("successfuly connected to query-svc-grpc")
	client := querypb.NewQueryServiceClient(conn)

	return &queryAdapter{
		client: client,
	}, nil
}

func (a *queryAdapter) GetTenLatestMessage(ctx context.Context, roomID string) (*querypb.GetTenLatestMessageResponse, error) {
	processPhotoRequest := &querypb.GetTenLatestMessageRequest{
		RoomId: roomID,
	}

	response, err := a.client.GetTenLatestMessage(ctx, processPhotoRequest)
	if err != nil {
		return nil, helper.FromGRPCError(err)
	}

	return response, nil
}
