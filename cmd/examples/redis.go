package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
)

type StructuredMessage struct {
	Sender  string `json:"sender"`
	Content []byte `json:"content"`
}

func (app *application) broadcastMessage(message []byte, roomNum int, sender *websocket.Conn) {
	structuredMessage := StructuredMessage{
		Sender:  sender.RemoteAddr().String(),
		Content: message,
	}

	msgBytes, err := json.Marshal(structuredMessage)
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
		var message StructuredMessage

		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			log.Println(err)
			continue
		}

		senderConn := app.findConnectionByAddr(message.Sender)

		app.sendToRoom(message.Content, roomNum, senderConn)
	}
}

func (app *application) findConnectionByAddr(addr string) *websocket.Conn {
	app.server.mu.Lock()
	defer app.server.mu.Unlock()

	for conn := range app.server.clients {
		if conn.RemoteAddr().String() == addr {
			return conn
		}
	}

	return nil
}

func (app *application) sendToRoom(message []byte, roomNum int, sender *websocket.Conn) {
	app.server.mu.Lock()
	defer app.server.mu.Unlock()

	for conn, room := range app.server.clients {
		if roomNum == room && conn != sender {
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				delete(app.server.clients, conn)
				conn.Close()
			}
		}
	}
}
