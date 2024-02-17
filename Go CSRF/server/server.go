package server

import (
	"log"
	"net/http"
	"github.com/bhupesh998/go-csrf/middleware"

)

func StartServer(hostname string , port string) error {
	host := hostname +":"+ port 
	log.Printf("Server Running on %s", host)

	handler := middleware.Newhandler()

	http.Handle("/", handler)
	return http.ListenAndServe(host, nil)
}
