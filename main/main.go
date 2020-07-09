package main

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
)

var appInstance *App
var router *mux.Router
var user *auth.UserRecord

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Errorf("port must be set")
	}
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal(err)
	}
}
func init() {
	initAppInstance()
	initRouter()
	initUser()
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

func initUser() {
	newUser, err := appInstance.authClient.GetUserByEmail(context.Background(), "versilistyson@gmail.com")
	if err != nil {
		panic(err)
	}
	user = newUser
}
