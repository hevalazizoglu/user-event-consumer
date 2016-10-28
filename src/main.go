package main

import (
		"log"
		"net/http"

		"github.com/gorilla/mux"
)

func route() {
		router := mux.NewRouter()
		router.HandleFunc("/user-event", HandleUserEvent).Methods("POST")
		router.HandleFunc("/user-event/_stats", ShowUserEventAPIStats).Methods("GET")

		log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
		route()
}
