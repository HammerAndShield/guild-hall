package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
)

type StructuredMessage struct {
	Sender  string `json:"sender"`
	Content []byte `json:"content"`
}

func (app *application) broadcastMessage(message []byte, roomNum int, user string) {
	msg := Message{
		MessageType: "message",
		Sender:      user,
		Message:     string(message),
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}

	app.rdb.Publish(context.Background(), strconv.Itoa(roomNum), msgBytes)
}

func (app *application) subscribeToRoom(ctx context.Context, roomNum int) {
	pubsub := app.rdb.Subscribe(ctx, strconv.Itoa(roomNum))
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		var message Message

		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			log.Println(err)
			continue
		}

		app.sendToRoom([]byte(msg.Payload), roomNum, message.Sender)
	}
}

func (app *application) sendToRoom(message []byte, roomNum int, sender string) {
	app.server.mu.RLock()
	defer app.server.mu.RUnlock()

	for _, user := range app.server.rooms[roomNum] {
		if user != sender {
			conn, err := app.findConnectionByUser(user)
			if err != nil {
				log.Println(err)
				continue
			}
			if err = conn.WriteMessage(websocket.TextMessage, message); err != nil {
				app.server.mu.Lock()
				conn.Close()
				delete(app.server.clients, user)
				app.server.mu.Unlock()
			}
		}
	}
}

func (app *application) findConnectionByUser(user string) (*websocket.Conn, error) {
	app.server.mu.RLock()
	defer app.server.mu.RUnlock()

	conn, exists := app.server.clients[user]
	if !exists {
		return nil, errors.New("user not found")
	}

	return conn, nil
}
