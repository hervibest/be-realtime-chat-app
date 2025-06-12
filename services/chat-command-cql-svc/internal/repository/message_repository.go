package repository

import (
	"be-realtime-chat-app/services/chat-command-cql-svc/internal/entity"
	"fmt"
	"time"
)

type messageRepository struct {
	session DB
}

type MessageRepository interface {
	Insert(message *entity.Message) error
	FindManyByRoomID(roomID string, limit int) (*[]*entity.Message, error)
	SoftDelete(id string) error
}

func NewMessageRepository(session DB) MessageRepository {
	return &messageRepository{session: session}
}

func (r *messageRepository) Insert(message *entity.Message) error {
	query := `INSERT INTO messages (id, uuid, room_id, user_id, username, content, created_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`

	q := r.session.Query(query,
		message.ID,
		message.UUID,
		message.RoomID,
		message.UserID,
		message.Username,
		message.Content,
		message.CreatedAt,
	)

	if err := q.Exec(); err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}
	return nil
}

func (r *messageRepository) FindManyByRoomID(roomID string, limit int) (*[]*entity.Message, error) {
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

func (r *messageRepository) SoftDelete(id string) error {
	query := `UPDATE messages SET deleted_at = ? WHERE id = ?`

	q := r.session.Query(query, time.Now(), id)

	if err := q.Exec(); err != nil {
		return fmt.Errorf("failed to soft delete message: %w", err)
	}
	return nil
}
