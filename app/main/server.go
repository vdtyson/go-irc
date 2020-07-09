package main

import (
	"context"
	"firebase.google.com/go/v4/auth"
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
	err := http.ListenAndServe(":8080", router)
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

// `C:\Users\Versilis\Desktop\Projects\go-irc\app\go-irc-firebase-adminsdk-e3m99-64808dbd66.json`
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
