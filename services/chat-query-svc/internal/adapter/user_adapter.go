package adapter

import (
	"be-realtime-chat-app/proto/userpb"
	"be-realtime-chat-app/services/commoner/discovery"
	"be-realtime-chat-app/services/commoner/helper"
	"be-realtime-chat-app/services/commoner/logs"
	"be-realtime-chat-app/services/commoner/utils"
	"context"
	"log"
)

type UserAdapter interface {
	AuthenticateUser(ctx context.Context, token string) (*userpb.AuthenticateResponse, error)
}

type userAdapter struct {
	client userpb.UserServiceClient
}

func NewUserAdapter(ctx context.Context, registry discovery.Registry, logs logs.Log) (UserAdapter, error) {
	userServiceName := utils.GetEnv("USER_SVC_NAME") + "-grpc"
	conn, err := discovery.ServiceConnection(ctx, userServiceName, registry, logs)
	if err != nil {
		return nil, err
	}

	log.Print("successfuly connected to user-svc-grpc")
	client := userpb.NewUserServiceClient(conn)

	return &userAdapter{
		client: client,
	}, nil
}

func (a *userAdapter) AuthenticateUser(ctx context.Context, token string) (*userpb.AuthenticateResponse, error) {
	processPhotoRequest := &userpb.AuthenticateRequest{
		Token: token,
	}

	response, err := a.client.AuthenticateUser(ctx, processPhotoRequest)
	if err != nil {
		return nil, helper.FromGRPCError(err)
	}

	return response, nil
}
