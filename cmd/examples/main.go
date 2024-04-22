package main

import (
	"fmt"
	"net/http"
)

func main() {
	server := NewServer()
	http.HandleFunc("/", server.echo)

	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
