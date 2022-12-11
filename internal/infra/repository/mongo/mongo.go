package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/alamrios/crabi-solution/config"
)

// NewClient returns conection to mongo database
func NewClient(ctx context.Context, cfg *config.Mongo) (*mongo.Database, error) {
	var auth string
	if cfg.User != "" && cfg.Password != "" {
		auth = cfg.User + ":" + cfg.Password + "@"
	}
	uri := cfg.Protocol + "://" + auth + cfg.URI

	timeout := 2 * time.Second
	clientOpts := options.Client().
		ApplyURI(uri).
		SetConnectTimeout(timeout).
		SetServerSelectionTimeout(timeout)

	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %v", err)
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialise mongo client: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to verify mongo client: %v", err)
	}

	return client.Database(cfg.Database), nil
}
