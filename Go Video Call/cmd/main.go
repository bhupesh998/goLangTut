package main

import (
"log"
"go-video-call/internal/server"


)

func main(){
	if err:= server.Run(); err != nil {
		log.Fatal(err.Error())
	}
}