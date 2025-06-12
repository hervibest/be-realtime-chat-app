package websockets

import (
	"be-realtime-chat-app/services/chat-realtime-svc/internal/adapter"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/model/event"
	"be-realtime-chat-app/services/commoner/constant/enum"
	"context"
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/oklog/ulid/v2"
)

type UserClient struct {
	Conn      *websocket.Conn
	Messaging adapter.MessagingAdapter
	UserID    string `json:"user_id"`
	RoomID    string `json:"room_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
}

//TODO get previous conversation/messages from database from query service (CQL DB)

func (c *UserClient) Subscriber() {
	// room, _ := c.RoomService.GetOneById(c.RoomId) //get previous conversation/messages from database
	// messages := room.Messages

	// for _, message := range messages { //sent previous conversation/messages to this client only
	// 	_ = c.Conn.WriteJSON(message)
	// }

	// Subscribe to room.<roomID> for receiving messages from others
	sub, err := c.Messaging.Subscribe("room."+c.RoomID, func(msg *nats.Msg) {
		var event event.Message
		if err := sonic.ConfigFastest.Unmarshal(msg.Data, &event); err != nil {
			log.Println("Failed to unmarshal event:", err)
			return
		}

		if event.RoomStatus == enum.RoomStatusEnumClosed {
			closeMsg := map[string]string{
				"type":    "room_deleted",
				"message": "Room has been deleted",
			}
			if err := c.Conn.WriteJSON(closeMsg); err != nil {
				log.Println("Failed to send close message:", err)
			}

			c.Conn.Close()
			return
		}

		if err := c.Conn.WriteJSON(msg.Data); err != nil {
			log.Println("WriteMessage error:", err)
		}
	})

	if err != nil {
		log.Println("NATS Subscribe error:", err)
		return
	}

	_ = c.Conn.Close()
	sub.Unsubscribe()
}

func (c *UserClient) Publisher() error {
	defer func() {
		_ = c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("ReadMessage error:", err)
			}
			break
		}

		log.Println("Received message:", string(m))
		loc, _ := time.LoadLocation("Asia/Jakarta") // Set your desired timezone
		now := time.Now().In(loc)                   // Get current time in the specified timezone
		event := &event.Message{
			ID:         ulid.Make().String(),
			UUID:       uuid.NewString(),
			RoomID:     c.RoomID,
			RoomStatus: enum.RoomStatusEnumActive,
			UserID:     c.UserID,
			Content:    string(m),
			Username:   c.Username,
			CreatedAt:  now.Format(time.RFC3339), // Format the time as needed
		}

		if err := c.Messaging.PublishMessage(context.TODO(), "room."+c.RoomID, event); err != nil {
			log.Println("NATS Publish error:", err)
			return err
		}
		log.Println("Message published to room:", c.RoomID)
	}
	return nil
}
