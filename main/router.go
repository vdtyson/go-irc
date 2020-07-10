package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	// auth
	router.HandleFunc("/register", RegisterHandler).Methods("POST")
	// channel
	router.HandleFunc("/channels/new", CreateGroupChannelHandler).Methods("POST")
	router.HandleFunc("/channels/message", NewMessageHandler).Methods("POST")
	router.HandleFunc("/channels/direct/{username1}/{username2}", CreateDMChannelHandler).Methods("POST")
	router.HandleFunc("/channels/messages", GetAllMessagesHandler).Methods("PUT")
	router.HandleFunc("/channels/users", AddUserHandler).Methods("PUT")
	// user
	router.HandleFunc("/users/{username}/channels", GetUserChannelsHandler).Methods("GET")
	// admin
	router.HandleFunc("/admin/ban", BanUserHandler).Methods("PUT")

	router.HandleFunc("/", HelloHandler).Methods("GET")
	return router
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Welcome!"))

}
