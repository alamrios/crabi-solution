package user

import (
	"context"
	"fmt"

	pld "github.com/alamrios/crabi-solution/internal/app/pld"
)

type Service struct {
	pldService pld.Service
	userRepo   Repository
}

func NewService(pldService pld.Service, userRepo Repository) (*Service, error) {
	return &Service{
		pldService: pldService,
		userRepo:   userRepo,
	}, nil
}

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
