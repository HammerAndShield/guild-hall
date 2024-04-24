package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type application struct {
	rdb    *redis.Client
	server *Server
}

func main() {
	server := NewServer()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	app := application{
		rdb:    rdb,
		server: server,
	}

	http.HandleFunc("GET /{id}", app.echo)

	fmt.Println("server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
