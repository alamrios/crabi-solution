package user

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alamrios/crabi-solution/internal/app/user"
)

const ResourceCollection = "users"

type Repository struct {
	mongoDB *mongo.Database
}

func New(mongoDB *mongo.Database) (*Repository, error) {
	if mongoDB == nil {
		return nil, fmt.Errorf("db is nil")
	}

	return &Repository{
		mongoDB: mongoDB,
	}, nil
}

// SaveUser method
func (r *Repository) SaveUser(ctx context.Context, user user.User) error {
	collection := r.mongoDB.Collection(ResourceCollection)

	_, err := collection.InsertOne(ctx, user)
	return err
}
