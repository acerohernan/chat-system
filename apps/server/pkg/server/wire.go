//go:build wireinject
// +build wireinject

package server

import (
	"github.com/chat-system/server/pkg/config"
	"github.com/chat-system/server/pkg/controllers"
	"github.com/chat-system/server/pkg/service"
	"github.com/chat-system/server/pkg/service/auth"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitializeServer(conf *config.Config) (*ChatServer, error) {
	wire.Build(
		auth.NewVerifier,
		service.GetMongoClient,
		createStorage,
		controllers.NewAuthController,
		controllers.NewRTCController,
		controllers.NewUserController,
		NewChatServer,
	)
	return &ChatServer{}, nil
}

func createStorage(conf *config.Config, mc *mongo.Client) service.PersistentStorage {
	if mc != nil {
		return service.NewMongoStorage(conf, mc)
	}
	return service.NewLocalStorage()
}
