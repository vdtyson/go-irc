package main

import (
	"context"
	"google.golang.org/api/option"
	"log"
)

func main() {
	ctx := context.Background()
	options := option.WithCredentialsFile(`C:\Users\Versilis\Desktop\Projects\go-irc\app\go-irc-firebase-adminsdk-e3m99-64808dbd66.json`)
	app, err := NewAppInstance(ctx, nil, options)
	if err != nil {
		log.Fatalf("error initializing app:%v\n", err)
	}

	firestoreClient, err := app.NewFirestoreClientInstance()
	if err != nil {
		log.Fatalf("error initializing firestore:%v\n", err)
	}
	defer firestoreClient.Close()
}
