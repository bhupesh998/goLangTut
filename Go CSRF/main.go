package main

import (
	"fmt"
	"log"
	"github.com/bhupesh998/go-csrf/db"
	"github.com/bhupesh998/go-csrf/server"
	"github.com/bhupesh998/go-csrf/middleware/myJWT"
)

var host ="localhost"
var port ="9000"

func main(){

	fmt.Println("CSRF PROJECT")
	db.InitDB()

	jwtErr := myJWT.InitJWT()
	if jwtErr != nil {
		log.Println("Error Intializing JWT ")
		log.Fatal(jwtErr)
	}

	serverErr := server.StartServer(host , port)
	if serverErr != nil {
		log.Println("Error Starting Server ")
		log.Fatal(serverErr)
	}
}