package mongo

import (
	"context"
	"time"

	"github.com/midnightrun/hexagonal-architecture-url-shortener-example/shortener"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func (m *mongoRepository) Find() {

}

func (m *mongoRepository) Find() {

}

func newMongoClient(mongoURL string, mongoTimeout time.Duration) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, err
}

func NewMongoRepository(mongoURL, mongoDB string, mongoTimeout time.Duration) (shortener.RedirectRepository, error) {
	repository := &mongoRepository{
		database: mongoDB,
		timeout:  time.Duration(mongoTimeout) * time.Second,
	}

	client, err := newMongoClient(mongoURL, repository.timeout)
	if err != nil {
		return nil, err
	}

	repository.client = client

	return client, nil
}
