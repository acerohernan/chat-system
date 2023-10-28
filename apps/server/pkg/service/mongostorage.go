package service

import (
	"context"

	"github.com/chat-system/server/pkg/config"
	"github.com/chat-system/server/pkg/logger"
	core "github.com/chat-system/server/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	PublicKeysCollection = "public_keys"
)

type MongoStorage struct {
	config *config.MongoConfig
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoStorage(config *config.MongoConfig, client *mongo.Client) *MongoStorage {
	return &MongoStorage{
		config: config,
		client: client,
		db:     client.Database(config.Database),
	}
}

func (s *MongoStorage) StorePublicKey(email core.UserEmail, key core.UserPublicKey) error {
	coll := s.db.Collection(PublicKeysCollection)

	_, err := coll.InsertOne(context.Background(), bson.M{"email": email, "publicKey": key})

	if err != nil {
		return err
	}

	return nil
}

func (s *MongoStorage) GetPublicKey(email core.UserEmail) (core.UserPublicKey, error) {
	coll := s.db.Collection(PublicKeysCollection)

	result := struct {
		Email     string `json:"email"`
		PublicKey string `json:"publicKey"`
	}{}

	err := coll.FindOne(context.Background(), bson.D{{Key: "email", Value: email}}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return "", ErrPublicKeyNotFound
	}

	if err != nil {
		return "", err
	}

	return core.UserPublicKey(result.PublicKey), nil
}

func (s *MongoStorage) Close() error {
	return s.client.Disconnect(context.TODO())
}

func getMongoClient(config *config.MongoConfig) (*mongo.Client, error) {
	logger.Infow("connecting to mongo db...")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.URI))

	if err != nil {
		logger.Errorw("failed at connecting to mongo", err)
		return nil, err
	}

	return client, nil
}
