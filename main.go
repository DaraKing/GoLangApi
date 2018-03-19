package GoLang_API_2

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

	log.Fatal(http.ListenAndServe(":8082", myRouter))
}
func main() {
	fmt.Println("GoLang API for parsing messages")
	handleRequest()
}
