package grpc

import (
	"be-realtime-chat-app/proto/roompb"
	"be-realtime-chat-app/services/commoner/helper"
	"be-realtime-chat-app/services/room-svc/internal/usecase"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type RoomGRPCHandler struct {
	roomUC usecase.RoomUseCase
	roompb.UnimplementedRoomServiceServer
}

func NewRoomGRPCHandler(server *grpc.Server, roomUC usecase.RoomUseCase) {
	handler := &RoomGRPCHandler{
		roomUC: roomUC,
	}
	roompb.RegisterRoomServiceServer(server, handler)
}

func (h *RoomGRPCHandler) GetRoom(ctx context.Context, req *roompb.GetRoomRequest) (*roompb.GetRoomResponse, error) {
	response, err := h.roomUC.GetRoom(ctx, req.GetRoomId())
	if err != nil {
		appErr, ok := err.(*helper.AppError)
		if ok {
			return nil, appErr.GRPCErrorCode()
		}
	}

	room := &roompb.Room{
		Id:     response.ID,
		Uuid:   response.UUID,
		Name:   response.Name,
		UserId: response.UserID,
	}

	return &roompb.GetRoomResponse{
		Status: int64(codes.OK),
		Room:   room,
	}, nil
}
