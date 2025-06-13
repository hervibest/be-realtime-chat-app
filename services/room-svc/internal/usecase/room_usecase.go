package usecase

import (
	errorcode "be-realtime-chat-app/services/commoner/constant/errcode"
	"be-realtime-chat-app/services/commoner/constant/message"
	"be-realtime-chat-app/services/commoner/helper"
	"be-realtime-chat-app/services/room-svc/internal/adapter"
	"be-realtime-chat-app/services/room-svc/internal/entity"
	"be-realtime-chat-app/services/room-svc/internal/helper/logs"
	"be-realtime-chat-app/services/room-svc/internal/model"
	"be-realtime-chat-app/services/room-svc/internal/model/converter"
	"be-realtime-chat-app/services/room-svc/internal/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RoomUseCase interface {
	CreateRoom(ctx context.Context, request *model.CreateRoomRequest) (*model.RoomResponse, error)
	GetActiveRooms(ctx context.Context) (*[]*model.RoomResponse, error)
	GetRoom(ctx context.Context, roomID string) (*model.RoomResponse, error)
	UserDeleteRoom(ctx context.Context, request *model.UserDeleteRoomRequest) error
	UserGetActiveRooms(ctx context.Context, userID string) (*[]*model.RoomResponse, error)
}

type roomUseCaseImpl struct {
	db              repository.DB
	roomRepository  repository.RoomRepository
	messaging       adapter.MessagingAdapter
	cacheAdapter    adapter.CacheAdapter
	customValidator helper.CustomValidator
	log             logs.Log
}

func NewRoomUseCase(db repository.DB, roomRepository repository.RoomRepository, messaging adapter.MessagingAdapter,
	cacheAdapter adapter.CacheAdapter, customValidator helper.CustomValidator,
	log logs.Log) RoomUseCase {
	return &roomUseCaseImpl{
		db:              db,
		roomRepository:  roomRepository,
		messaging:       messaging,
		cacheAdapter:    cacheAdapter,
		customValidator: customValidator,
		log:             log,
	}
}

func (uc *roomUseCaseImpl) GetActiveRooms(ctx context.Context) (*[]*model.RoomResponse, error) {
	uc.log.Info("Fetching active rooms")
	rooms, err := uc.roomRepository.FindMany(ctx, uc.db)
	if err != nil {
		return nil, helper.WrapInternalServerError(uc.log, "failed to find active rooms", err)
	}

	if len(*rooms) == 0 {
		uc.log.Info("No active rooms found")
		return nil, nil
	}

	return converter.RoomsToResponses(rooms), nil
}

func (uc *roomUseCaseImpl) UserGetActiveRooms(ctx context.Context, userID string) (*[]*model.RoomResponse, error) {
	uc.log.Info("Fetching room active rooms")
	rooms, err := uc.roomRepository.FindManyByUserID(ctx, uc.db, userID)
	if err != nil {
		return nil, helper.WrapInternalServerError(uc.log, "failed to find active rooms", err)
	}

	if len(*rooms) == 0 {
		uc.log.Info("No active rooms found")
		return nil, nil
	}

	return converter.RoomsToResponses(rooms), nil
}

func (uc *roomUseCaseImpl) UserDeleteRoom(ctx context.Context, request *model.UserDeleteRoomRequest) error {
	if validatonErrs := uc.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return validatonErrs
	}

	if err := repository.BeginTxx(ctx, uc.db, func(tx repository.TX) error {
		uc.log.Info("Deleting room", zap.String("roomUUID", request.RoomUUID), zap.String("roomID", request.UserID))
		room, err := uc.roomRepository.FindByUUIDAndUserID(context.Background(), uc.db, request.RoomUUID, request.UserID)
		if err != nil {
			if strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
				return helper.NewUseCaseError(errorcode.ErrResourceNotFound, message.RoomNotFound)
			}
			return helper.WrapInternalServerError(uc.log, "failed to FindByUUIDAndUserID", err)
		}

		if err := uc.roomRepository.DeleteByUUID(ctx, uc.db, request.RoomUUID); err != nil {
			if strings.Contains(err.Error(), message.InternalNoRowsAffected) {
				return helper.NewUseCaseError(errorcode.ErrInvalidArgument, "room not found or already deleted")
			}
			return helper.WrapInternalServerError(uc.log, "failed to delete room by uuid", err)
		}

		topic := "room." + room.ID
		initialMessage := []byte("room deleted by " + room.UserID)
		if err := uc.messaging.PublishMessage(ctx, topic, initialMessage); err != nil {
			uc.log.Warn("Failed to publish to NATS", zap.String("topic", topic), zap.Error(err))
			return helper.WrapInternalServerError(uc.log, "failed to publish new room to nats", err)
		}
		return nil
	}); err != nil {
		uc.log.Warn("Transaction failed", zap.Error(err))
		return helper.WrapInternalServerError(uc.log, "transaction failed", err)
	}

	return nil
}

func (uc *roomUseCaseImpl) CreateRoom(ctx context.Context, request *model.CreateRoomRequest) (*model.RoomResponse, error) {
	if validatonErrs := uc.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return nil, validatonErrs
	}

	response := new(entity.Room)
	if err := repository.BeginTxx(ctx, uc.db, func(tx repository.TX) error {
		uc.log.Info("Create room request")
		room := &entity.Room{
			ID:     ulid.Make().String(),
			Name:   request.Name,
			UserID: request.UserID,
		}

		txRoom, err := uc.roomRepository.Insert(ctx, uc.db, room)
		if err != nil {
			uc.log.Warn("Failed to create room", zap.Error(err))
			return helper.WrapInternalServerError(uc.log, "failed to insert new room", err)
		}

		response = txRoom
		topic := "room." + room.ID
		initialMessage := []byte("room created by " + room.UserID)
		if err := uc.messaging.PublishMessage(ctx, topic, initialMessage); err != nil {
			uc.log.Warn("Failed to publish to NATS", zap.String("topic", topic), zap.Error(err))
			return helper.WrapInternalServerError(uc.log, "failed to publish new room to nats", err)
		}

		return nil
	}); err != nil {
		uc.log.Warn("Transaction failed", zap.Error(err))
		return nil, helper.WrapInternalServerError(uc.log, "transaction failed", err)
	}

	return converter.RoomToResponse(response), nil
}

func (uc *roomUseCaseImpl) GetRoom(ctx context.Context, roomID string) (*model.RoomResponse, error) {
	uc.log.Info("Fetching room active rooms")
	room, err := uc.findCachedByRoomID(ctx, roomID)
	if err != nil {
		return nil, helper.WrapInternalServerError(uc.log, "failed to find active room", err)
	}

	return converter.RoomToResponse(room), nil
}

func (uc *roomUseCaseImpl) findCachedByRoomID(ctx context.Context, roomID string) (*entity.Room, error) {
	cachedRoom, err := uc.cacheAdapter.Get(ctx, roomID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, helper.WrapInternalServerError(uc.log, "failed to get cached room", err)
	}

	var room *entity.Room

	//fallback to database if cache is stale or not found
	if errors.Is(err, redis.Nil) {
		room, err = uc.roomRepository.FindByID(context.Background(), uc.db, roomID)
		if err != nil {
			if strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
				return nil, helper.NewUseCaseError(errorcode.ErrResourceNotFound, message.RoomNotFound)
			}
			return nil, helper.WrapInternalServerError(uc.log, "failed to FindByID", err)
		}
		if err := uc.saveRoomToCache(ctx, room, time.Now().Add(24*time.Hour)); err != nil {
			return nil, helper.WrapInternalServerError(uc.log, "failed to save room to cache", err)
		}
	} else {
		room = &entity.Room{}
		if err := sonic.ConfigFastest.Unmarshal([]byte(cachedRoom), room); err != nil {
			return nil, helper.WrapInternalServerError(uc.log, "failed to unmarshal room body from cached", err)
		}
	}

	return room, nil
}

func (uc *roomUseCaseImpl) saveRoomToCache(ctx context.Context, room *entity.Room, expiresAt time.Time) error {
	jsonValue, err := sonic.ConfigFastest.Marshal(room)
	if err != nil {
		return fmt.Errorf("marshal room : %+v", err)
	}

	if err := uc.cacheAdapter.Set(ctx, room.ID, jsonValue, time.Until(expiresAt)); err != nil {
		return fmt.Errorf("save room body into cache : %+v", err)
	}

	return nil
}
