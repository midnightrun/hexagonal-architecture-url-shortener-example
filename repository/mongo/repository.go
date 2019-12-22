package mongo

import (
	"context"
	"time"

	"github.com/midnightrun/hexagonal-architecture-url-shortener-example/shortener"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func (m *mongoRepository) Find(code string) (*shortener.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	redirect := &shortener.Redirect{}
	collection := m.client.Database(m.database).Collection("redirects")
	filter := bson.M{"code": code}

	err := collection.FindOne(ctx, filter).Decode(&redirect)
	if err != nil {
		return nil, err
	}

	return redirect, nil
}

func (m *mongoRepository) Store(redirect *shortener.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	collection := m.client.Database(m.database).Collection("redirects")

	_, err := collection.InsertOne(
		ctx,
		bson.M{
			"code":       redirect.Code,
			"url":        redirect.URL,
			"created_at": redirect.CreatedAt,
		})
	if err != nil {
		return err
	}

	return nil
}

func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
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

func NewMongoRepository(mongoURL, mongoDB string, mongoTimeout int) (shortener.RedirectRepository, error) {
	repository := &mongoRepository{
		database: mongoDB,
		timeout:  time.Duration(mongoTimeout) * time.Second,
	}

	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, err
	}

	repository.client = client

	return repository, nil
}
