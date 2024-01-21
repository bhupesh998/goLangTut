package main

import (
	"fmt"
	"net/http"

	"github.com/bhupesh998/gochatapp/pkg/websocket"
)

func serverWS(pool *websocket.Pool, w http.ResponseWriter, r *http.Request){
	fmt.Println("Websocket endpoint reached")

	conn , err := websocket.Upgrade(w,r)

	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn : conn,
		Pool : pool,
	}

	pool.Register <- client
	client.Read()
}

func SetupRoutes(){
	pool := websocket.NewPool()
	go pool.Start()
	
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
		serverWS(pool , w , r)
	})
}

func main(){
	fmt.Println("Full stack chat app backend")
	SetupRoutes()
	http.ListenAndServe(":9000", nil)
}