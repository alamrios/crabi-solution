package pld

import (
	"context"

	"github.com/alamrios/crabi-solution/internal/app/pld"
)

type Service struct {
}

func NewService() (*Service, error) {
	return &Service{}, nil
}

func (s *Service) CheckBlackList(ctx context.Context, request pld.Request) error {
	return nil
}
