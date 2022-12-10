package user

import (
	"context"
)

type Repository interface {
	SaveUser(ctx context.Context, user User) error
}
