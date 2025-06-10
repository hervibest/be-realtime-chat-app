package websockets

import (
	"be-realtime-chat-app/services/chat-realtime-svc/internal/adapter"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/model/event"
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/websocket"
	"github.com/nats-io/nats.go"
)

type UserClient struct {
	Conn      *websocket.Conn
	Messaging adapter.Messaging
	UserID    string `json:"user_id"`
	RoomID    string `json:"room_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
}

func (c *UserClient) Subscriber() {
	// room, _ := c.RoomService.GetOneById(c.RoomId) //get previous conversation/messages from database
	// messages := room.Messages

	// for _, message := range messages { //sent previous conversation/messages to this client only
	// 	_ = c.Conn.WriteJSON(message)
	// }

	// Subscribe to room.<roomID> for receiving messages from others
	sub, err := c.Messaging.Subscribe("room."+c.RoomID, func(msg *nats.Msg) {
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
			RoomID:    c.RoomID,
			UserID:    c.UserID,
			Content:   string(m),
			Username:  c.Username,
			CreatedAt: now.Format(time.RFC3339), // Format the time as needed
		}

		value, _ := sonic.ConfigFastest.Marshal(event)

		if err := c.Messaging.Publish("room."+c.RoomID, value); err != nil {
			log.Println("NATS Publish error:", err)
			return err
		}
		log.Println("Message published to room:", c.RoomID)
	}
	return nil
}
