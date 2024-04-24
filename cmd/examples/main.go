package main

import (
	"flag"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type application struct {
	rdb    *redis.Client
	server *Server
	port   string
}

func main() {
	port := flag.String("port", ":8080", "Port for the server to listen on")
	flag.Parse()

	server := NewServer()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	app := application{
		rdb:    rdb,
		server: server,
		port:   *port,
	}

	http.HandleFunc("GET /{id}", app.echo)

	fmt.Printf("server started on %s", app.port)
	err := http.ListenAndServe(app.port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
