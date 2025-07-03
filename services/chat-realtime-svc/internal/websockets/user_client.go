package websockets

import (
	"be-realtime-chat-app/services/chat-realtime-svc/internal/adapter"
	"be-realtime-chat-app/services/chat-realtime-svc/internal/model/event"
	"be-realtime-chat-app/services/commoner/constant/enum"
	"be-realtime-chat-app/services/commoner/logs"
	"context"
	"errors"
	"log"
	"sync"
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
	writeMu      sync.Mutex
	closed       bool // Flag untuk menandai koneksi sudah ditutup
}

// SafeWriteJSON safely writes JSON to the WebSocket connection
func (c *UserClient) SafeWriteJSON(v interface{}) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	if c.closed {
		return errors.New("connection closed")
	}
	return c.Conn.WriteJSON(v)
}

// SafeClose safely closes the WebSocket connection
func (c *UserClient) SafeClose() error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	if c.closed {
		return nil
	}
	c.closed = true
	c.Log.Info("Closing WebSocket connection", zap.String("roomID", c.RoomID), zap.String("userID", c.UserID))
	// Send close message to client (optional but proper)
	_ = c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "room closed"))

	return c.Conn.Close()
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

	if len(latestMessagesPb.Message) != 0 && latestMessagesPb.Message != nil {
		for _, msg := range latestMessagesPb.Message {
			event := &event.Message{
				ID:        msg.GetId(),
				RoomID:    msg.GetRoomId(),
				UserID:    msg.GetUserId(),
				Username:  msg.GetUsername(),
				Content:   msg.GetContent(),
				CreatedAt: msg.GetCreatedAt(),
			}

			if err := c.SafeWriteJSON(event); err != nil { // Gunakan SafeWriteJSON di sini
				log.Println("Failed to send message:", err)
				return
			}
		}
	}

	for {
		select {
		case <-done:
			c.Log.Info("Exiting subscriber loop due to WebSocket close")
			return
		default:
			msg, err := sub.NextMsg(2 * time.Second)
			if err != nil {
				if err == nats.ErrTimeout {
					continue
				}
				log.Println("Error receiving message from NATS:", err)
				return
			}

			var event event.Message
			if err := sonic.ConfigFastest.Unmarshal(msg.Data, &event); err != nil {
				log.Println("Failed to unmarshal event:", err)
				continue
			}

			c.Log.Info("Received & parsed message",
				zap.String("eventID", event.ID),
				zap.String("roomID", event.RoomID),
				zap.String("userID", event.UserID),
				zap.String("content", event.Content))

			if event.RoomStatus == enum.RoomStatusEnumClosed {
				c.Log.Info("Room has been closed", zap.String("roomID", event.RoomID), zap.String("userID", event.UserID))
				closeMsg := map[string]string{
					"type":    "room_deleted",
					"message": "Room has been deleted",
				}
				if err := c.SafeWriteJSON(closeMsg); err != nil {
					log.Println("Failed to send close message:", err)
				}
				c.SafeClose()
				log.Println("WebSocket connection closed for room:", event.RoomID, "user:", event.UserID)
				return
			}

			if err := c.SafeWriteJSON(event); err != nil {
				log.Println("WriteMessage error:", err)
				return
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
