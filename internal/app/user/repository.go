package user

import (
	"context"
)

type Repository interface {
	SaveUser(ctx context.Context, user User) error
	GetUserByEmail(ctx context.Context, idUser string) (*User, error)
}
