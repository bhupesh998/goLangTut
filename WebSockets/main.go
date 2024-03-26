package main

import (
	"context"
	"log"
	"net/http"
)

func main() {
	setupApi()

	log.Fatal(http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil))
}

func setupApi() {

	ctx := context.Background()

	manager := NewManager(ctx)


	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.serveWS)
	http.HandleFunc("/login", manager.LoginHandler)
}