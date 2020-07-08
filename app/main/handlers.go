package main

import (
	"cloud.google.com/go/firestore"
	"context"
)

type Handler interface {
	Create(client *firestore.Client, ctx context.Context) error
	Delete(client *firestore.Client, ctx context.Context) error
	BaseCollection() string
}
