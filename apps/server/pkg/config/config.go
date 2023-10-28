package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port   string
	Host   string
	Redis  *RedisConfig
	Logger *LoggerConfig
	Auth   *AuthConfig
	Mongo  *MongoConfig
}

type RedisConfig struct {
	Host string
}

type LoggerConfig struct {
	Level string
}

type AuthConfig struct {
	JWTIssuer          string
	JWTSecret          string
	GoogleClientId     string
	GoogleClientSecret string
}

type MongoConfig struct {
	URI      string
	Database string
}

func NewConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Port: os.Getenv("PORT"),
		Host: os.Getenv("HOST"),
		Redis: &RedisConfig{
			Host: os.Getenv("REDIS_HOST"),
		},
		Logger: &LoggerConfig{
			Level: os.Getenv("LOG_LEVEL"),
		},
		Auth: &AuthConfig{
			JWTIssuer:          os.Getenv("AUTH_JWT_ISSUER"),
			JWTSecret:          os.Getenv("AUTH_JWT_SECRET"),
			GoogleClientId:     os.Getenv("AUTH_GOOGLE_CLIENT_ID"),
			GoogleClientSecret: os.Getenv("AUTH_GOOGLE_CLIENT_SECRET"),
		},
		Mongo: &MongoConfig{
			URI:      os.Getenv("MONGO_URI"),
			Database: os.Getenv("MONGO_DB"),
		},
	}
}
