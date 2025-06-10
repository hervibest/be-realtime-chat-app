package repository

import (
	"be-realtime-chat-app/services/chat-realtime-svc/internal/entity"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/helper/logs"
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"go.uber.org/zap"
)

type RoomRepository interface {
	Insert(ctx context.Context, db Querier, room *entity.Room) (*entity.Room, error)
	FindByUUID(ctx context.Context, db Querier, uuid string) (*entity.Room, error)
	FindByID(ctx context.Context, db Querier, id string) (*entity.Room, error)
	FindManyByUserID(ctx context.Context, db Querier, userID string) (*[]*entity.Room, error)
	FindMany(ctx context.Context, db Querier) (*[]*entity.Room, error)
	DeleteByUUID(ctx context.Context, db Querier, uuid string) error
}
type roomRepositoryImpl struct {
	log logs.Log
}

func NewRoomRepository(log logs.Log) RoomRepository {
	return &roomRepositoryImpl{log: log}
}

func (r *roomRepositoryImpl) Insert(ctx context.Context, db Querier, room *entity.Room) (*entity.Room, error) {
	query := `
	INSERT INTO rooms
		(id, name, user_id)
	VALUES
		($1, $2, $3)`

	_, err := db.Exec(ctx, query, room.ID, room.Name)
	if err != nil {
		r.log.Error("failed to exec insert query", zap.String("query", query), zap.Error(err))
		return nil, err
	}
	return room, nil
}

func (r *roomRepositoryImpl) FindByUUID(ctx context.Context, db Querier, uuid string) (*entity.Room, error) {
	room := new(entity.Room)
	query := `SELECT * FROM rooms WHERE uuid = $1 AND deleted_at IS NULL`
	if err := pgxscan.Get(ctx, db, room, query, uuid); err != nil {
		r.log.Error("failed to get query", zap.String("query", query), zap.Error(err))
		return nil, err
	}
	return room, nil
}

func (r *roomRepositoryImpl) FindByID(ctx context.Context, db Querier, id string) (*entity.Room, error) {
	room := new(entity.Room)
	query := `SELECT * FROM rooms WHERE id = $1 AND deleted_at IS NULL`
	if err := pgxscan.Get(ctx, db, room, query, id); err != nil {
		r.log.Error("failed to get query", zap.String("query", query), zap.Error(err))
		return nil, err
	}
	return room, nil
}

func (r *roomRepositoryImpl) DeleteByUUID(ctx context.Context, db Querier, uuid string) error {
	query := `UPDATE rooms SET deleted_at = now() WHERE uuid = $1 AND deleted_at IS NOT NULL`
	row, err := db.Exec(ctx, query, uuid)
	if err != nil {
		r.log.Error("failed to exec delete query", zap.String("query", query), zap.Error(err))
		return err
	}

	if row.RowsAffected() == 0 {
		return errors.New("invalid room uuid")
	}

	return nil
}

// TODO pagination
func (r *roomRepositoryImpl) FindManyByUserID(ctx context.Context, db Querier, userID string) (*[]*entity.Room, error) {
	query := `SELECT * FROM rooms WHERE user_id = $1 AND deleted_at IS NULL`
	var rooms []*entity.Room
	if err := pgxscan.Select(ctx, db, &rooms, query, userID); err != nil {
		r.log.Error("failed to get query", zap.String("query", query), zap.Error(err))
		return nil, err
	}
	return &rooms, nil
}

// TODO pagination
func (r *roomRepositoryImpl) FindMany(ctx context.Context, db Querier) (*[]*entity.Room, error) {
	query := `SELECT * FROM rooms WHERE deleted_at IS NULL`
	var rooms []*entity.Room
	if err := pgxscan.Select(ctx, db, &rooms, query); err != nil {
		r.log.Error("failed to get query", zap.String("query", query), zap.Error(err))
		return nil, err
	}
	return &rooms, nil
}
