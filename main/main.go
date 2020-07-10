package main

import (
	"context"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
)

var appInstance *App
var router *mux.Router

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := http.ListenAndServe(":"+port, router)
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
	options := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	newAppInstance, err := NewAppInstance(context.Background(), nil, options)
	if err != nil {
		panic(err)
	}
	appInstance = newAppInstance
}
