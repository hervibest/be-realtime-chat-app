package archiverepo

// import (
// 	"be-realtime-chat-app/services/chat-command-cql-svc/internal/entity"
// 	"fmt"
// 	"time"

// 	"github.com/scylladb/gocqlx/v2"
// 	"github.com/scylladb/gocqlx/v2/qb"
// )

// type messageRepository struct {
// 	session gocqlx.Session
// }

// type MessageRepository interface {
// 	Insert(message *entity.Message) error
// 	FindManyByRoomID(roomID string, limit int) (*[]*entity.Message, error)
// 	SoftDelete(id string) error
// }

// func NewMessageRepository(session gocqlx.Session) MessageRepository {
// 	return &messageRepository{session: session}
// }

// func (r *messageRepository) Insert(message *entity.Message) error {
// 	query := qb.Insert("messages").
// 		Columns("id", "uuid", "room_id", "user_id", "username", "content", "created_at").
// 		Query(r.session).
// 		BindStruct(message)

// 	if err := query.ExecRelease(); err != nil {
// 		return fmt.Errorf("failed to insert message: %w", err)
// 	}
// 	return nil
// }

// func (r *messageRepository) FindManyByRoomID(roomID string, limit int) (*[]*entity.Message, error) {
// 	var messages []*entity.Message

// 	query := qb.Select("messages").
// 		Where(qb.Eq("room_id")).
// 		Limit(uint(limit)).
// 		Query(r.session).
// 		BindMap(qb.M{
// 			"room_id": roomID,
// 		})

// 	if err := query.SelectRelease(&messages); err != nil {
// 		return nil, fmt.Errorf("failed to find messages by room ID: %w", err)
// 	}

// 	return &messages, nil
// }

// func (r *messageRepository) SoftDelete(id string) error {
// 	query := qb.Update("messages").
// 		Set("deleted_at").
// 		Where(qb.Eq("id")).
// 		Query(r.session).
// 		BindMap(qb.M{
// 			"id":         id,
// 			"deleted_at": time.Now(),
// 		})

// 	if err := query.ExecRelease(); err != nil {
// 		return fmt.Errorf("failed to soft delete message: %w", err)
// 	}
// 	return nil
// }
