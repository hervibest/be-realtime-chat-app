package repository

import (
	"be-realtime-chat-app/services/chat-query-svc/internal/entity"
	"fmt"

	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
)

type messageCQLRepoImpl struct {
	session gocqlx.Session
}

type MessageCQLRepository interface {
	FindManyByRoomID(roomID string, limit int) (*[]*entity.Message, error)
}

func NewMessageCQLRepository(session gocqlx.Session) MessageCQLRepository {
	return &messageCQLRepoImpl{session: session}
}

func (r *messageCQLRepoImpl) FindManyByRoomID(roomID string, limit int) (*[]*entity.Message, error) {
	var messages []*entity.Message

	query := qb.Select("messages").
		Where(qb.Eq("room_id")).
		Limit(uint(limit)).
		Query(r.session).
		BindMap(qb.M{
			"room_id": roomID,
		})

	if err := query.SelectRelease(&messages); err != nil {
		return nil, fmt.Errorf("failed to find messages by room ID: %w", err)
	}

	return &messages, nil
}
