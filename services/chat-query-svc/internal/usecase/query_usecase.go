package usecase

import (
	"be-realtime-chat-app/services/chat-query-svc/internal/adapter"
	"be-realtime-chat-app/services/chat-query-svc/internal/model"
	"be-realtime-chat-app/services/chat-query-svc/internal/model/converter"
	"be-realtime-chat-app/services/chat-query-svc/internal/repository"
	errorcode "be-realtime-chat-app/services/commoner/constant/errcode"
	"be-realtime-chat-app/services/commoner/helper"
	"be-realtime-chat-app/services/commoner/logs"
	"context"
)

type QueryUseCase interface {
	GetTenLatestMessage(ctx context.Context, roomID string) (*[]*model.MessageResponse, error)
	SearchMessages(ctx context.Context, params *model.SearchParams) (*[]*model.MessageResponse, error)
}

type queryUseCaseImpl struct {
	messageCQLRepository     repository.MessageCQLRepository
	messageElasticRepository repository.MessageElasticRepo
	roomAdapter              adapter.RoomAdapter
	db                       repository.DB
	log                      logs.Log
}

func NewQueryUseCase(messageCQLRepository repository.MessageCQLRepository, messageElasticRepository repository.MessageElasticRepo,
	roomAdapter adapter.RoomAdapter, db repository.DB, log logs.Log) QueryUseCase {
	return &queryUseCaseImpl{
		messageCQLRepository:     messageCQLRepository,
		messageElasticRepository: messageElasticRepository,
		roomAdapter:              roomAdapter,
		db:                       db,
		log:                      log,
	}
}

func (uc *queryUseCaseImpl) GetTenLatestMessage(ctx context.Context, roomID string) (*[]*model.MessageResponse, error) {
	messages, err := uc.messageCQLRepository.FindManyByRoomID(roomID, 10)
	if err != nil {
		return nil, err
	}

	return converter.MessagesToResponses(messages), nil
}

func (uc *queryUseCaseImpl) SearchMessages(ctx context.Context, params *model.SearchParams) (*[]*model.MessageResponse, error) {
	_, err := uc.roomAdapter.GetRoom(ctx, params.RoomID)
	if appErr, ok := err.(*helper.AppError); ok {
		if appErr.Code == errorcode.ErrUserNotFound {
			return nil, helper.NewUseCaseError(errorcode.ErrUserNotFound, "Room not found")
		}
		return nil, appErr
	}

	messages, err := uc.messageElasticRepository.SearchMessages(params)
	if err != nil {
		return nil, err
	}

	return converter.MessagesToResponses(messages), nil
}
