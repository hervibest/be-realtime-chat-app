package repository

import (
	"be-realtime-chat-app/services/chat-query-svc/internal/entity"
	"fmt"

	"github.com/gocql/gocql"
)

type messageCQLRepoImpl struct {
	session *gocql.Session
}

type MessageCQLRepository interface {
	FindManyByRoomID(roomID string, limit int) (*[]*entity.Message, error)
}

func NewMessageCQLRepository(session *gocql.Session) MessageCQLRepository {
	return &messageCQLRepoImpl{session: session}
}

func (r *messageCQLRepoImpl) FindManyByRoomID(roomID string, limit int) (*[]*entity.Message, error) {
	var messages []*entity.Message

	query := `SELECT id, uuid, room_id, user_id, username, content, created_at, deleted_at 
	          FROM messages WHERE room_id = ? LIMIT ?`

	iter := r.session.Query(query, roomID, limit).Iter()

	var msg entity.Message
	for iter.Scan(
		&msg.ID,
		&msg.UUID,
		&msg.RoomID,
		&msg.UserID,
		&msg.Username,
		&msg.Content,
		&msg.CreatedAt,
		&msg.DeletedAt,
	) {
		// Create a new copy of the message to avoid overwriting
		m := msg
		messages = append(messages, &m)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to find messages by room ID: %w", err)
	}

	return &messages, nil
}
