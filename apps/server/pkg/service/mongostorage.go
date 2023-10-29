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
	UsersCollection      = "users"
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

func (s *MongoStorage) StoreUser(user *core.User) error {
	coll := s.db.Collection(UsersCollection)

	_, err := coll.InsertOne(context.Background(), user)

	if err != nil {
		return err
	}

	return nil
}

func (s *MongoStorage) GetUser(id core.UserId) (*core.User, error) {
	coll := s.db.Collection(UsersCollection)

	u := &core.User{}

	err := coll.FindOne(context.Background(), bson.D{{Key: "id", Value: id}}).Decode(u)

	if err == mongo.ErrNoDocuments {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *MongoStorage) GetUserWithEmail(email core.UserEmail) (*core.User, error) {
	coll := s.db.Collection(UsersCollection)

	u := &core.User{}

	err := coll.FindOne(context.Background(), bson.D{{Key: "email", Value: email}}).Decode(u)

	if err == mongo.ErrNoDocuments {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *MongoStorage) Close() error {
	return s.client.Disconnect(context.TODO())
}

func GetMongoClient(config *config.MongoConfig) (*mongo.Client, error) {
	logger.Infow("connecting to mongo db...")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.URI))

	if err != nil {
		logger.Errorw("failed at connecting to mongo", err)
		return nil, err
	}

	return client, nil
}
