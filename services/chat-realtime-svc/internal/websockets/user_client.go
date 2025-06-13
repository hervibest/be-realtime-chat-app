package websockets

import (
	"be-realtime-chat-app/services/chat-realtime-svc/internal/adapter"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/model/event"
	"be-realtime-chat-app/services/commoner/constant/enum"
	"be-realtime-chat-app/services/commoner/logs"
	"context"
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
)

type UserClient struct {
	Log          logs.Log
	Conn         *websocket.Conn
	Messaging    adapter.MessagingAdapter
	QueryAdapter adapter.QueryAdapter
	UserID       string `json:"user_id"`
	RoomID       string `json:"room_id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
}

func (c *UserClient) Subscriber(done chan struct{}) {
	c.Log.Info("UserClient Subscriber started", zap.String("roomID", c.RoomID), zap.String("userID", c.UserID))

	sub, err := c.Messaging.SubscribeSync("room." + c.RoomID)
	if err != nil {
		log.Println("NATS SubscribeSync error:", err)
		return
	}
	defer sub.Unsubscribe()

	latestMessagesPb, err := c.QueryAdapter.GetTenLatestMessage(context.Background(), c.RoomID)
	if err != nil {
		c.Log.Error("Error when get ten latest message form query svc", zap.Error(err))
	}

	if len(latestMessagesPb.Message) != 0 {
		for _, msg := range latestMessagesPb.Message {
			event := &event.Message{
				ID: msg.GetId(),
				// UUID:       msg.Get(),
				RoomID: msg.GetRoomId(),
				// RoomStatus: msg.String(),
				UserID:    msg.GetUserId(),
				Username:  msg.GetUsername(),
				Content:   msg.GetContent(),
				CreatedAt: msg.String(),
			}

			if err := c.Conn.WriteJSON(event); err != nil {
				log.Println("Failed to send close message:", err)
			}
		}
	}

	for {
		select {
		case <-done:
			c.Log.Info("Exiting subscriber loop due to WebSocket close")
			return
		default:
			// Tunggu pesan NATS dengan timeout
			msg, err := sub.NextMsg(10 * time.Second)
			if err != nil {
				if err == nats.ErrTimeout {
					continue // tidak ada pesan, ulangi
				}
				log.Println("Error receiving message from NATS:", err)
				return
			}

			var event event.Message
			if err := sonic.ConfigFastest.Unmarshal(msg.Data, &event); err != nil {
				log.Println("Failed to unmarshal event:", err)
				continue
			}

			c.Log.Info("Received & parsed message", zap.String("eventID", event.ID), zap.String("roomID", event.RoomID), zap.String("userID", event.UserID), zap.String("content", event.Content))

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

			if err := c.Conn.WriteJSON(event); err != nil {
				log.Println("WriteMessage error:", err)
			}
		}
	}
}

func (c *UserClient) Publisher(done chan struct{}) error {
	defer func() {
		_ = c.Conn.Close()
	}()

	for {
		c.Log.Info("UserClient Publisher started", zap.String("roomID", c.RoomID), zap.String("userID", c.UserID))
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("ReadMessage error:", err)
				done <- struct{}{} // Signal the subscriber to stop
			}
			break
		}
		c.Log.Info("Received message from WebSocket", zap.String("roomID", c.RoomID), zap.String("userID", c.UserID), zap.String("message", string(m)))
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
			CreatedAt:  now.Format(time.RFC3339Nano), // Format the time as needed
		}
		c.Log.Info("Publishing message", zap.String("eventID", event.ID), zap.String("username", event.Username), zap.String("userID", event.UserID), zap.String("content", event.Content))
		if err := c.Messaging.PublishMessage(context.TODO(), "room."+c.RoomID, event); err != nil {
			log.Println("NATS Publish error:", err)
			return err
		}
		log.Println("Message published to room:", c.RoomID)
	}
	return nil
}
