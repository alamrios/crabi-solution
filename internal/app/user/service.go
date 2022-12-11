package user

import (
	"context"
	"fmt"

	pld "github.com/alamrios/crabi-solution/internal/app/pld"
)

// Service struct for users service
type Service struct {
	pldService pld.Service
	userRepo   Repository
}

// NewService returns an instance of users service
func NewService(pldService pld.Service, userRepo Repository) (*Service, error) {
	if pldService == nil {
		return nil, fmt.Errorf("pld service is nil")
	}

	if userRepo == nil {
		return nil, fmt.Errorf("user repo is nil")
	}

	return &Service{
		pldService: pldService,
		userRepo:   userRepo,
	}, nil
}

// CreateUser stores given user in users repository if valid, error otherwise
func (s *Service) CreateUser(ctx context.Context, user User) (*User, error) {
	pldErr := s.pldService.CheckBlacklist(
		ctx,
		pld.Request{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		})
	if pldErr != nil {
		return nil, pldErr
	}

	duplicate, err := s.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	if duplicate != nil {
		return nil, fmt.Errorf("user with email %s already exists", user.Email)
	}

	rErr := s.userRepo.SaveUser(ctx, user)
	if rErr != nil {
		return nil, rErr
	}

	return &user, nil
}
