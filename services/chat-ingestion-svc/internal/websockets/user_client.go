package websockets

import (
	"be-realtime-chat-app/services/chat-ingestion-svc/internal/adapter"
	"log"

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

		if err := c.Messaging.Publish("room."+c.RoomID, m); err != nil {
			log.Println("NATS Publish error:", err)
			return err
		}
		log.Println("Message published to room:", c.RoomID)
	}
	return nil
}
