package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/contacts", GetContacts).Methods("GET")
	router.HandleFunc("/api/v1/contacts/{id:[0-9]+}", GetContact).Methods("GET")
	router.HandleFunc("/api/v1/contacts", CreateContact).Methods("POST")
	router.HandleFunc("/api/v1/contacts/{id:[0-9]+}", UpdateContact).Methods("PUT")
	router.HandleFunc("/api/v1/contacts", DeleteContacts).Methods("DELETE")
	router.HandleFunc("/api/v1/contacts/{id:[0-9]+}", DeleteContact).Methods("DELETE")

	http.ListenAndServe(":3001", router)
}
