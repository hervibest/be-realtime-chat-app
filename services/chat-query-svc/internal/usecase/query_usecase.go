package usecase

import (
	"be-realtime-chat-app/services/chat-query-svc/internal/helper/logs"
	"be-realtime-chat-app/services/chat-query-svc/internal/model"
	"be-realtime-chat-app/services/chat-query-svc/internal/model/converter"
	"be-realtime-chat-app/services/chat-query-svc/internal/repository"
	"context"
)

type QueryUseCase interface {
	GetTenLatestMessage(ctx context.Context, roomID string) (*[]*model.MessageResponse, error)
}

type queryUseCaseImpl struct {
	messageCQLRepository     repository.MessageCQLRepository
	messageElasticRepository repository.MessageElasticRepo
	db                       repository.DB
	log                      logs.Log
}

func NewQueryUseCase(messageCQLRepository repository.MessageCQLRepository, messageElasticRepository repository.MessageElasticRepo,
	db repository.DB, log logs.Log) QueryUseCase {
	return &queryUseCaseImpl{
		messageCQLRepository:     messageCQLRepository,
		messageElasticRepository: messageElasticRepository,
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
	messages, err := uc.messageElasticRepository.SearchMessages(params)
	if err != nil {
		return nil, err
	}

	return converter.MessagesToResponses(messages), nil
}
