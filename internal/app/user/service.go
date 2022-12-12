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
	if user.FirstName == "" {
		return nil, fmt.Errorf("user's first name should not be empty")
	}
	if user.LastName == "" {
		return nil, fmt.Errorf("user's last name should not be empty")
	}
	if user.Email == "" {
		return nil, fmt.Errorf("user's email should not be empty")
	}
	if user.Password == "" {
		return nil, fmt.Errorf("user's password should not be empty")
	}

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

// Login returns user with given credentials
func (s *Service) Login(ctx context.Context, email, password string) (*User, error) {
	if email == "" {
		return nil, fmt.Errorf("user's email should not be empty")
	}
	if password == "" {
		return nil, fmt.Errorf("user's password should not be empty")
	}

	user, err := s.userRepo.GetUserByEmailAndPassword(ctx, email, password)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not exists or invalid credentials")
	}

	return user, nil
}

// GetUser returns user with given email
func (s *Service) GetUser(ctx context.Context, email string) (*User, error) {
	if email == "" {
		return nil, fmt.Errorf("user's email should not be empty")
	}

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not exists")
	}

	return user, nil
}
