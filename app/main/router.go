package main

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/register", RegisterHandler).Methods("POST")
	return router
}
