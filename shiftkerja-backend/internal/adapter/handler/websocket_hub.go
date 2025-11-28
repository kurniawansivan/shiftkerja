package handler

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Upgrader turns a normal HTTP request into a WebSocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// ALLOW CORS for WebSockets
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Hub maintains the set of active clients
type Hub struct {
	// A map to track active clients (Thread-safe)
	Clients map[*websocket.Conn]bool
	// Mutex to lock the map when adding/removing clients
	Mutex sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[*websocket.Conn]bool),
	}
}

// HandleWS is the endpoint: ws://localhost:8080/ws
func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request) {
	// 1. Upgrade HTTP -> WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("❌ Failed to upgrade WS: %v\n", err)
		return
	}

	// 2. Register Client
	h.Mutex.Lock()
	h.Clients[conn] = true
	h.Mutex.Unlock()
	
	fmt.Println("🟢 New Client Connected!")

	// Ensure connection is closed and removed when function exits
	defer func() {
		h.Mutex.Lock()
		delete(h.Clients, conn)
		h.Mutex.Unlock()
		conn.Close()
		fmt.Println("🔴 Client Disconnected")
	}()

	// 3. Listen for Messages (The Broadcast Loop)
	for {
		// Read incoming JSON message
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			// If error (client disconnected), break the loop
			break 
		}

		fmt.Printf("⚡ Broadcasting: %v\n", msg)

		// 4. BROADCAST to ALL connected clients
		h.Mutex.Lock()
		for client := range h.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Printf("❌ Failed to write to client: %v\n", err)
				client.Close()
				delete(h.Clients, client)
			}
		}
		h.Mutex.Unlock()
	}
}