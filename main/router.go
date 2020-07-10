package main

import (
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	// auth
	router.HandleFunc("/register", RegisterHandler).Methods("POST")
	// channel
	router.HandleFunc("/channels/new", CreateGroupChannelHandler).Methods("POST")
	router.HandleFunc("/channels/message", NewMessageHandler).Methods("POST")
	router.HandleFunc("/channels/direct/{username1}/{username2}", CreateDMChannelHandler).Methods("POST")
	router.HandleFunc("/channels/messages", GetAllMessagesHandler).Methods("POST")
	router.HandleFunc("/channels/users", AddUserHandler).Methods("PUT")
	router.HandleFunc("/channels/messages/newest", GetNewestMessageHandler).Methods("POST")
	// user
	router.HandleFunc("/users/{username}/channels", GetUserChannelsHandler).Methods("GET")
	// admin
	router.HandleFunc("/admin/ban", BanUserHandler).Methods("PUT")

	router.PathPrefix("/").Handler(httpSwagger.WrapHandler)
	return router
}
