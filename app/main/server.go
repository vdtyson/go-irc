package main

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	"log"
)

type (
	Doc interface {
		Create(client *firestore.Client, ctx context.Context) error
		Delete(client *firestore.Client, ctx context.Context) error
	}
	Channel struct {
		Name string
	}
)

// cli-command
func (c *Channel) Create(client *firestore.Client, ctx context.Context) error {
	_, err := client.Collection("channels").Doc(c.Name).Set(ctx, map[string]interface{}{
		"message1": "Test message 1",
		"message2": "Test message 2",
	})
	return err
}

func (c *Channel) Delete(client *firestore.Client, ctx context.Context) error {
	_, err := client.Collection("channels").Doc(c.Name).Delete(ctx)
	return err
}

func main() {
	ctx := context.Background()
	// `C:\Users\Versilis\Desktop\Projects\go-irc\app\go-irc-firebase-adminsdk-e3m99-64808dbd66.json`
	opt := option.WithCredentialsFile(`C:\Users\Versilis\Desktop\Projects\go-irc\app\go-irc-firebase-adminsdk-e3m99-64808dbd66.json`)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app:%v\n", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing firestore:%v\n", err)
	}
	defer client.Close()
	newChannel := Channel{"#test6"}
	newChannel.Create(client, ctx)
	newChannel2 := Channel{"#test7"}
	newChannel2.Create(client, ctx)
	newChannel2.Delete(client, ctx)
}
