package usecase

import (
	"be-realtime-chat-app/services/chat-realtime-svc/internal/adapter"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/entity"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/helper/logs"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/model"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/model/converter"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/repository"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/websockets"
	"context"
	"errors"

	"github.com/gofiber/contrib/websocket"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
)

type RoomUseCase interface {
	CreateRoom(ctx context.Context, request *model.CreateRoomRequest) (*model.RoomResponse, error)
	JoinRoom(conn *websocket.Conn, request *model.JoinRoomRequest) error
}

type roomUseCaseImpl struct {
	db             repository.DB
	roomRepository repository.RoomRepository
	mesaging       adapter.Messaging
	log            logs.Log
}

func NewRoomUseCase(db repository.DB, roomRepository repository.RoomRepository, mesaging adapter.Messaging, log logs.Log) RoomUseCase {
	return &roomUseCaseImpl{
		db:             db,
		roomRepository: roomRepository,
		mesaging:       mesaging,
		log:            log,
	}
}

func (uc *roomUseCaseImpl) CreateRoom(ctx context.Context, request *model.CreateRoomRequest) (*model.RoomResponse, error) {
	response := new(entity.Room)
	if err := repository.BeginTxx(ctx, uc.db, func(tx repository.TX) error {
		uc.log.Info("Create room request")
		room := &entity.Room{
			ID:     ulid.Make().String(),
			Name:   request.Name,
			UserID: request.UserID,
		}

		txResponse, err := uc.roomRepository.Insert(ctx, uc.db, room)
		if err != nil {
			uc.log.Warn("Failed to create room", zap.Error(err))
			return errors.New("failed to create room: " + err.Error())
		}

		response = txResponse
		topic := "room." + room.ID
		initialMessage := []byte("room created by " + room.UserID)
		if err := uc.mesaging.Publish(topic, initialMessage); err != nil {
			uc.log.Warn("Failed to publish to NATS", zap.String("topic", topic), zap.Error(err))
		}
		return nil
	}); err != nil {
		uc.log.Warn("Transaction failed", zap.Error(err))
		return nil, errors.New("transaction failed: " + err.Error())
	}

	return converter.RoomToResponse(response), nil
}

func (uc *roomUseCaseImpl) JoinRoom(conn *websocket.Conn, request *model.JoinRoomRequest) error {
	uc.log.Info("Join room request", zap.String("roomID", request.RoomID), zap.String("userID", request.UserID))

	room, err := uc.roomRepository.FindByID(context.Background(), uc.db, request.RoomID)
	if err != nil {
		uc.log.Warn("Failed to find room", zap.Error(err))
		return errors.New("failed to find room: " + err.Error())
	}

	if room == nil {
		return errors.New("room not found")
	}

	client := &websockets.UserClient{
		Conn:   conn,
		RoomID: room.ID,
		UserID: request.UserID,
	}

	go client.Subscriber()
	client.Publisher()

	return nil
}
