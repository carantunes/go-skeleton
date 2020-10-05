package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const DefaultCounter int = 1

type Counter struct {
	ID    int
	count int
	date  time.Time
}

// TransferStorageService is responsible to interact with the transfer repository
type CounterStorageService struct {
	DB DBService
}

// NewCounterStorageService creates a new counter storage service instance
func NewCounterStorageService(db DBService) CounterStorageService {
	return CounterStorageService{DB: db}
}

// Increment the counter on the database
func (storage CounterStorageService) Inc(ctx context.Context, counter int) (int, error) {
	newTransfer.CreatedAt = carbon.Now().Time
	collection := storage.DB.client.Database("go-skeleton").Collection("counters")

	var got Counter
	filter := bson.M{"_id": counter}
	err := collection.FindOne(ctx, filter).Decode(&got)
	if err != nil {
		log.Fatal(err)
	}

	count = got.count + 1
	res, err := collection.UpdateOne(ctx, filter, Counter{
		ID:    got.ID,
		Count: count,
		Date:  time.Now(),
	})
	id := res.InsertedID

	return count, nil
}
