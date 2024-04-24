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
	clients map[string]*websocket.Conn
	rooms   map[int][]string
	subs    map[int]context.CancelFunc
	mu      sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		clients: make(map[string]*websocket.Conn),
		rooms:   make(map[int][]string),
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

	user := r.URL.Query().Get("user")
	if user == "" {
		http.Error(w, "Please provide a username", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	app.server.mu.Lock()
	app.server.clients[user] = conn
	if _, exists := app.server.rooms[roomNum]; !exists {
		app.server.rooms[roomNum] = []string{user}
		ctx, cancel := context.WithCancel(context.Background())
		app.server.subs[roomNum] = cancel
		go app.subscribeToRoom(ctx, roomNum)
	} else {
		app.server.rooms[roomNum] = append(app.server.rooms[roomNum], user)
	}
	app.server.mu.Unlock()

	defer func() {
		app.server.mu.Lock()
		delete(app.server.clients, user)
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

		app.broadcastMessage(message, roomNum, user)
	}
}

func (app *application) getClientsInRoom(roomNum int) []*websocket.Conn {
	app.server.mu.RLock()
	defer app.server.mu.RUnlock()

	var clients []*websocket.Conn
	for _, user := range app.server.rooms[roomNum] {
		clients = append(clients, app.server.clients[user])
	}
	return clients
}
