package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
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

var jwtSecret []byte

// Claims for JWT
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Message struct {
	Username  string `json:"username"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"` // Optional, for client-side display
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

	//log.Printf("Attempting to register: username='%s', password='%s'", user.Username, user.Password)

	if user.Username == "" {
		log.Println("Empty username received")
		http.Error(w, "Username cannot be empty", http.StatusBadRequest)
		return
	}

	err = RegisterUser(db, user.Username, user.Password)
	if err != nil {
		log.Println("Registration error:", err)
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Println("User registered successfully:", user.Username)
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
	username := clients[conn] // Assuming username is stored in clients map
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}
		broadcastMessage(username, string(msg))
	}
}

// Broadcast message to all connected clients
func broadcastMessage(username, message string) {
	mutex.Lock()
	defer mutex.Unlock()
	msg := Message{
		Username:  username,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339), // ISO 8601 format
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Println("JSON marshal error:", err)
		return
	}
	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			log.Println("Write error:", err)
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

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}
	jwtSecret = []byte(secret)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback for local development
	}

	// Start the server
	log.Println("Server started on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
