// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"github.com/chat-system/server/pkg/config"
	"github.com/chat-system/server/pkg/controllers"
	"github.com/chat-system/server/pkg/service"
	"github.com/chat-system/server/pkg/service/auth"
	"go.mongodb.org/mongo-driver/mongo"
)

// Injectors from wire.go:

func InitializeServer(conf *config.Config) (*ChatServer, error) {
	rtcController := controllers.NewRTCController()
	verifier := auth.NewVerifier(conf)
	client, err := service.GetMongoClient(conf)
	if err != nil {
		return nil, err
	}
	persistentStorage := createStorage(conf, client)
	authController := controllers.NewAuthController(conf, verifier, persistentStorage)
	userController := controllers.NewUserController(persistentStorage, verifier)
	chatServer, err := NewChatServer(conf, rtcController, authController, userController, persistentStorage)
	if err != nil {
		return nil, err
	}
	return chatServer, nil
}

// wire.go:

func createStorage(conf *config.Config, mc *mongo.Client) service.PersistentStorage {
	if mc != nil {
		return service.NewMongoStorage(conf, mc)
	}
	return service.NewLocalStorage()
}
