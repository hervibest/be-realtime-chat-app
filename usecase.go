package berealtimechatapp

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/nats-io/nats.go"
)

type WSClient struct {
	Conn   *websocket.Conn
	RoomID string
	UserID string
}

var (
	nc      *nats.Conn
	clients = make(map[string][]*WSClient) // key = roomID
)

func main() {
	var err error
	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Drain()

	app := fiber.New()

	app.Get("/ws/:room/:user", websocket.New(handleWebSocket))

	log.Println("Server running on :3000")
	log.Fatal(app.Listen(":3000"))
}

func handleWebSocket(c *websocket.Conn) {
	roomID := c.Params("room")
	userID := c.Params("user")

	client := &WSClient{
		Conn:   c,
		RoomID: roomID,
		UserID: userID,
	}

	// Add client to memory
	clients[roomID] = append(clients[roomID], client)

	// Subscribe to room.<roomID> for receiving messages from others
	sub, err := nc.Subscribe("room."+roomID, func(msg *nats.Msg) {
		// Relay message to all clients in the room
		for _, cl := range clients[roomID] {
			if err := cl.Conn.WriteMessage(websocket.TextMessage, msg.Data); err != nil {
				log.Println("WriteMessage error:", err)
			}
		}
	})
	if err != nil {
		log.Println("NATS Subscribe error:", err)
		return
	}
	defer sub.Unsubscribe()

	log.Printf("User %s joined room %s", userID, roomID)

	// Read loop from WebSocket and forward to NATS
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("ReadMessage error for user %s: %v", userID, err)
			break
		}

		message := "[" + userID + "] " + string(msg)
		nc.Publish("room."+roomID, []byte(message))
	}

	// Cleanup on disconnect
	c.Close()
	cleanupClient(client)
}

func cleanupClient(cl *WSClient) {
	roomID := cl.RoomID
	newClients := []*WSClient{}
	for _, client := range clients[roomID] {
		if client != cl {
			newClients = append(newClients, client)
		}
	}
	clients[roomID] = newClients
	log.Printf("User %s left room %s", cl.UserID, cl.RoomID)
}
