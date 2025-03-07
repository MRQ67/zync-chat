package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

var (
	// Upgrader to convert HTTP to WebSocket
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow connections from any origin
		},
	}

	// Mutex to protect the clients map
	mutex sync.Mutex

	// Map to store connected clients
	clients = make(map[*websocket.Conn]bool)
)

var jwtSecret = []byte("7FADN3XFg3Xz7jPSuqoUsWsIxL2uLUGwPGOCjmomcCk=")

// Claims for JWT
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Handle WebSocket connections
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Register the new client
	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	// Listen for messages from this client
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			// Remove client on disconnect
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}

		// Broadcast the message to all clients
		broadcastMessage(message)
	}
}

// Broadcast message to all connected clients
func broadcastMessage(message []byte) {
	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error writing message:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func main() {
	// Serve static files from the "static" folder
	http.Handle("/", http.FileServer(http.Dir("static")))

	// WebSocket endpoint
	http.HandleFunc("/ws", handleConnections)

	// Start the server
	log.Println("Server starting on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
