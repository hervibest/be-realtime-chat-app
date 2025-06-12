package grpc

import (
	"be-realtime-chat-app/proto/querypb"
	"be-realtime-chat-app/services/chat-query-svc/internal/usecase"
	"be-realtime-chat-app/services/commoner/helper"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type QueryGRPCHandler struct {
	queryUC usecase.QueryUseCase
	querypb.UnimplementedQueryServiceServer
}

func NewQueryGRPCHandler(server *grpc.Server, queryUC usecase.QueryUseCase) {
	handler := &QueryGRPCHandler{
		queryUC: queryUC,
	}
	querypb.RegisterQueryServiceServer(server, handler)
}

func (h *QueryGRPCHandler) GetTenLatestMessage(ctx context.Context, req *querypb.GetTenLatestMessageRequest) (*querypb.GetTenLatestMessageResponse, error) {
	response, err := h.queryUC.GetTenLatestMessage(ctx, req.GetRoomId())
	if err != nil {
		appErr, ok := err.(*helper.AppError)
		if ok {
			return nil, appErr.GRPCErrorCode()
		}
	}

	messagePbs := make([]*querypb.Message, len(*response))
	for i, message := range *response {
		messagePbs[i] = &querypb.Message{
			Id:        message.ID,
			RoomId:    message.RoomID,
			UserId:    message.UserID,
			Content:   message.Content,
			CreatedAt: message.CreatedAt,
		}
	}

	return &querypb.GetTenLatestMessageResponse{
		Status:  int64(codes.OK),
		Message: messagePbs,
	}, nil
}
