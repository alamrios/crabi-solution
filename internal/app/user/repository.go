package user

import (
	"context"
)

// Repository contract for users repository
type Repository interface {
	SaveUser(ctx context.Context, user User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}
