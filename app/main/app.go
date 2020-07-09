package main

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type App struct {
	firebaseInstance *firebase.App
	authClient       *auth.Client
	firestoreClient  *firestore.Client
	authRepo         *AuthRepository
	channelRepo      *ChannelRepository
	userRepo         *UserRepository
}

func NewAppInstance(ctx context.Context, config *firebase.Config, options ...option.ClientOption) (*App, error) {
	fbApp, err := firebase.NewApp(ctx, config, options...)
	if err != nil {
		return nil, err
	}

	authClient, err := fbApp.Auth(ctx)
	if err != nil {
		return nil, err
	}

	firestoreClient, err := fbApp.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	authRepo := &AuthRepository{firestoreClient: firestoreClient, authClient: authClient}
	channelRepo := &ChannelRepository{fsClient: firestoreClient}
	userRepo := &UserRepository{firestoreClient}

	return &App{
		firebaseInstance: fbApp,
		authClient:       authClient,
		firestoreClient:  firestoreClient,
		authRepo:         authRepo,
		channelRepo:      channelRepo,
		userRepo:         userRepo,
	}, nil
}
