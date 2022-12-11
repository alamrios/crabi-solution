package pld

import (
	"context"
)

// Request struct for PLD service
type Request struct {
	FirstName string
	LastName  string
	Email     string
}

// Service contract for PLD service
type Service interface {
	CheckBlacklist(ctx context.Context, request Request) error
}
