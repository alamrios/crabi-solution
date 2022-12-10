package config

import (
	"context"
	"fmt"
	"os"
)

type Config struct {
	Mongo Mongo
}

type Mongo struct {
	Protocol string
	URI      string
	Database string
	User     string
	Password string
}

func New(ctx context.Context) (*Config, error) {
	mongoProtocol := os.Getenv("MONGO_PROTOCOL")
	mongoURI := os.Getenv("MONGO_URI")
	mongoDatabase := os.Getenv("MONGO_DATABASE")
	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")

	if mongoProtocol == "" {
		return nil, fmt.Errorf("MONGO_PROTOCOL env var needed")
	}
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI env var needed")
	}
	if mongoDatabase == "" {
		return nil, fmt.Errorf("MONGO_DATABASE env var needed")
	}
	if mongoUser == "" {
		return nil, fmt.Errorf("MONGO_USER env var needed")
	}
	if mongoPassword == "" {
		return nil, fmt.Errorf("MONGO_PASSWORD env var needed")
	}

	return &Config{
		Mongo: Mongo{
			Protocol: mongoProtocol,
			URI:      mongoURI,
			Database: mongoDatabase,
			User:     mongoUser,
			Password: mongoPassword,
		},
	}, nil
}
