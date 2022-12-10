package config

import (
	"context"
	"fmt"
	"os"
)

type Config struct {
	Mongo Mongo
	PLD   PLD
}

type Mongo struct {
	Protocol string
	URI      string
	Database string
	User     string
	Password string
}

type PLD struct {
	Protocol string
	Host     string
	Port     string
	URI      string
}

func New(ctx context.Context) (*Config, error) {
	mongo := Mongo{
		Protocol: os.Getenv("MONGO_PROTOCOL"),
		URI:      os.Getenv("MONGO_URI"),
		Database: os.Getenv("MONGO_DATABASE"),
		User:     os.Getenv("MONGO_USER"),
		Password: os.Getenv("MONGO_PASSWORD"),
	}

	if mongo.Protocol == "" {
		return nil, fmt.Errorf("MONGO_PROTOCOL env var needed")
	}
	if mongo.URI == "" {
		return nil, fmt.Errorf("MONGO_URI env var needed")
	}
	if mongo.Database == "" {
		return nil, fmt.Errorf("MONGO_DATABASE env var needed")
	}
	if mongo.User == "" {
		return nil, fmt.Errorf("MONGO_USER env var needed")
	}
	if mongo.Password == "" {
		return nil, fmt.Errorf("MONGO_PASSWORD env var needed")
	}

	pld := PLD{
		Protocol: os.Getenv("PLD_PROTOCOL"),
		Host:     os.Getenv("PLD_HOST"),
		Port:     os.Getenv("PLD_PORT"),
		URI:      os.Getenv("PLD_URI"),
	}

	if pld.Protocol == "" {
		return nil, fmt.Errorf("PLD_PROTOCOL env var needed")
	}
	if pld.Host == "" {
		return nil, fmt.Errorf("PLD_HOST env var needed")
	}
	if pld.Port == "" {
		return nil, fmt.Errorf("PLD_PORT env var needed")
	}
	if pld.URI == "" {
		return nil, fmt.Errorf("PLD_URI env var needed")
	}

	return &Config{
		Mongo: mongo,
		PLD:   pld,
	}, nil
}
