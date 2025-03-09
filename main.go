package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

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
	db    *sql.DB
	// Map to store connected clients
	clients = make(map[*websocket.Conn]string)
)

var jwtSecret = []byte("7FADN3XFg3Xz7jPSuqoUsWsIxL2uLUGwPGOCjmomcCk=")

// Claims for JWT
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = RegisterUser(db, user.Username, user.Password)
	if err != nil {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = AuthenticateUser(db, user.Username, user.Password)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// Handle WebSocket connections
func handleConnections(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Register the new client
	mutex.Lock()
	clients[conn] = claims.Username
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
		broadcastMessage(claims.Username + ": " + string(message))
	}
}

// Broadcast message to all connected clients
func broadcastMessage(message string) {
	mutex.Lock()
	defer mutex.Unlock()

	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Error writing message:", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}

func main() {
	db = initDB()
	defer db.Close()
	// Serve static files from the "static" folder
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	// WebSocket endpoint
	http.HandleFunc("/ws", handleConnections)

	// Start the server
	log.Println("Server starting on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
