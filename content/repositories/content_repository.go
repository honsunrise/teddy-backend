package repositories

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type ContentRepository interface {
}

func NewContentRepository(client *mongo.Client) (ContentRepository, error) {
	return &contentRepository{
		ctx:         context.Background(),
		client:      client,
		collections: client.Database("Teddy").Collection("Content"),
	}, nil
}

type contentRepository struct {
	ctx         context.Context
	client      *mongo.Client
	collections *mongo.Collection
}
