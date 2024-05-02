package main

type Message struct {
	MessageType string `json:"type"`
	Sender      string `json:"sender"`
	Message     string `json:"message"`
}
