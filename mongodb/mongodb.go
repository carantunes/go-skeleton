package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// DB represents the database type
type DBService struct {
	client *mongo.Client
}

// New creates a new mongodb db service
func New(client *mongo.Client) (DBService, error) {
	return DBService{databaseClient}, nil
}
