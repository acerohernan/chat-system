//go:build mage
// +build mage

package main

import (
	"context"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/chat-system/server/pkg/config"
	"github.com/chat-system/server/pkg/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Proto() error {
	grpcProtoFiles := []string{
		"models.proto",
		"rtc.proto",
	}

	location := "../../packages/protocol"
	target := "./proto"

	fmt.Println("generating protobuf...")

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	protoc, err := getToolPath("protoc")
	if err != nil {
		return err
	}

	protocGoPath, err := getToolPath("protoc-gen-go")
	if err != nil {
		return err
	}

	protocGrpcGoPath, err := getToolPath("protoc-gen-go-grpc")
	if err != nil {
		return err
	}

	args := append([]string{
		"--go_out", target,
		"--go-grpc_out", target,
		"--go_opt=paths=source_relative",
		"--go-grpc_opt=paths=source_relative",
		"--plugin=go=" + protocGoPath,
		"--plugin=go-grpc=" + protocGrpcGoPath,
		"-I=" + location,
	}, grpcProtoFiles...)
	cmd := exec.Command(protoc, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func Mongo() error {
	// This script will created the first data needed and indexes for query performance
	config := config.NewConfig()

	client, err := service.GetMongoClient(config.Mongo)

	if err != nil {
		return err
	}

	db := client.Database(config.Mongo.Database)

	pubKeysColl := db.Collection(service.PublicKeysCollection)

	_, err = pubKeysColl.Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.M{
		"useremail": 1,
	}, Options: options.Index().SetUnique(true)})

	if err != nil {
		return err
	}

	return nil
}

func getToolPath(name string) (string, error) {
	if p, err := exec.LookPath(name); err == nil {
		return p, nil
	}
	// check under gopath
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	p := filepath.Join(gopath, "bin", name)
	if _, err := os.Stat(p); err != nil {
		return "", err
	}
	return p, nil
}
