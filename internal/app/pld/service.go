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
	CheckBlacklist(ctx context.Context, request Request) error
}
