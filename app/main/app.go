package main

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type App struct {
	firebaseInstance *firebase.App
	context          context.Context
}

func NewAppInstance(ctx context.Context, config *firebase.Config, options ...option.ClientOption) (*App, error) {
	fbApp, err := firebase.NewApp(ctx, config, options...)
	if err != nil {
		return nil, err
	}

	return &App{firebaseInstance: fbApp, context: ctx}, nil
}

func (a *App) NewFirestoreClientInstance() (*firestore.Client, error) {
	return a.firebaseInstance.Firestore(a.context)
}
