package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func handleRequest() {

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/match", getMessage).Methods("POST")
	myRouter.HandleFunc("/match", optionsRequest).Methods("OPTIONS")

	myRouter.HandleFunc("/match/{id}", getMatch).Methods("GET")
	myRouter.HandleFunc("/match/{id}", optionsRequest).Methods("OPTIONS")

	myRouter.HandleFunc("/player/{steamID}", getPlayer).Methods("GET")
	myRouter.HandleFunc("/player/{steamID}", optionsRequest).Methods("OPTIONS")

	myRouter.HandleFunc("/map/{name}", mapInfo).Methods("GET")
	myRouter.HandleFunc("/map/{name}", optionsRequest).Methods("OPTIONS")

	log.Fatal(http.ListenAndServe(":8082", myRouter))
}
func main() {
	fmt.Println("GoLang API for parsing messages")
	handleRequest()
}
