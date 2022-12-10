package pld

import (
	"context"
)

type Request struct {
	FirstName string
	LastName  string
	Email     string
}

type Service interface {
	CheckBlackList(ctx context.Context, request Request) error
}
