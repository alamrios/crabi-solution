package user

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alamrios/crabi-solution/internal/app/user"
)

// ResourceCollection name of mongo collection
const ResourceCollection = "users"

// Repository struct for users mongo repository
type Repository struct {
	mongoDB *mongo.Database
}

// New returns an instance of users mongo repository
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

// GetUserByEmail returns user in mongo collection with given email
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	collection := r.mongoDB.Collection(ResourceCollection)

	query := collection.FindOne(
		ctx,
		bson.M{
			"email": email,
		},
	)

	if query.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}

	var user user.User
	err := query.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
