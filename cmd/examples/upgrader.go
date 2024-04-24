package main

import (
	"context"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Accepts all origins
	},
}

type Server struct {
	clients map[*websocket.Conn]int
	rooms   map[int]bool
	subs    map[int]context.CancelFunc
	mu      sync.Mutex
}

func NewServer() *Server {
	return &Server{
		clients: make(map[*websocket.Conn]int),
		rooms:   make(map[int]bool),
		subs:    make(map[int]context.CancelFunc),
	}
}

func (app *application) echo(w http.ResponseWriter, r *http.Request) {
	room := r.PathValue("id")

	roomNum, err := strconv.Atoi(room)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	app.server.mu.Lock()
	app.server.clients[conn] = roomNum
	if _, exists := app.server.rooms[roomNum]; !exists {
		app.server.rooms[roomNum] = true
		ctx, cancel := context.WithCancel(context.Background())
		app.server.subs[roomNum] = cancel
		go app.subscribeToRoom(ctx, roomNum)
	}
	app.server.mu.Unlock()

	defer func() {
		app.server.mu.Lock()
		delete(app.server.clients, conn)
		if len(app.getClientsInRoom(roomNum)) == 0 {
			delete(app.server.rooms, roomNum)
			if cancel, exists := app.server.subs[roomNum]; exists {
				cancel()
				delete(app.server.subs, roomNum)
			}
		}
		app.server.mu.Unlock()
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		app.broadcastMessage(message, roomNum, conn)
	}
}

func (app *application) getClientsInRoom(roomNum int) []*websocket.Conn {
	var clients []*websocket.Conn
	for conn, room := range app.server.clients {
		if roomNum == room {
			clients = append(clients, conn)
		}
	}
	return clients
}

func (server *Server) broadcastMessageDirect(message []byte, roomNum int, sender *websocket.Conn) {
	server.mu.Lock()
	defer server.mu.Unlock()

	for conn, room := range server.clients {
		if roomNum == room && conn != sender {
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				delete(server.clients, conn)
				conn.Close()
			}
		}
	}
}
