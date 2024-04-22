package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Accepts all origins
	},
}

type Server struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func NewServer() *Server {
	return &Server{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (server *Server) echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	server.mu.Lock()
	server.clients[conn] = true
	server.mu.Unlock()

	defer func() {
		server.mu.Lock()
		delete(server.clients, conn)
		server.mu.Unlock()
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break // Exit the loop if there's an error (e.g., client disconnects)
		}

		// Broadcast the message to all clients
		server.broadcastMessage(message)
	}
}

func (server *Server) broadcastMessage(message []byte) {
	server.mu.Lock()
	defer server.mu.Unlock()

	for conn := range server.clients {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			// Optional: handle error, possibly removing the client from the map
			delete(server.clients, conn)
			conn.Close()
		}
	}
}
