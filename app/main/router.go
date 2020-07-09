package main

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	// auth
	router.HandleFunc("/register", RegisterHandler).Methods("POST")
	// channel
	router.HandleFunc("/channels/new", CreateGroupChannelHandler).Methods("POST")
	router.HandleFunc("/channels/message", NewMessageHandler).Methods("POST")
	router.HandleFunc("/channels/direct/{username1}/{username2}", CreateDMChannelHandler).Methods("POST")
	router.HandleFunc("/channels/messages", GetAllMessagesHandler).Methods("GET")
	return router
}
