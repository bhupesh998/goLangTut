package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Hello Mod in goLang")
	greet()
	r := mux.NewRouter()
	r.HandleFunc("/", servHome).Methods("GET")

	log.Fatal(http.ListenAndServe(":4000", r))

}

func greet() {
	fmt.Println("Hello Folks")
}

func servHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>I am learning go</h1>"))
}
