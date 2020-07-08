package main

import (
	"context"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
	"log"
	"net/http"
)

var appInstance *App
var router *mux.Router

func main() {
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
func init() {
	initAppInstance()
	initRouter()
}

func initRouter() {
	router = NewRouter()
}

func initAppInstance() {
	options := option.WithCredentialsFile(`C:\Users\Versilis\Desktop\Projects\go-irc\app\go-irc-firebase-adminsdk-e3m99-64808dbd66.json`)
	newAppInstance, err := NewAppInstance(context.Background(), nil, options)
	if err != nil {
		panic(err)
	}
	appInstance = newAppInstance
}
