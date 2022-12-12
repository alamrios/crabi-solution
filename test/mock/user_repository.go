// Source: internal/app/user/repository.go
package mock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/alamrios/crabi-solution/internal/app/user"
)

// UserRepository is a mock of Repository interface
type UserRepository struct {
	mock.Mock
}

// NewUserRepository creates new user mock repository
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// AddCall adds new call to the mock
func (m *UserRepository) AddCall(t *testing.T, calls []Call) *UserRepository {
	t.Helper()

	for _, call := range calls {
		m.On(call.FunctionName, call.Params...).Return(call.Returns...)
	}

	return m
}

// SaveUser method mock
func (m *UserRepository) SaveUser(_ context.Context, user user.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// GetUserByEmail method mock
func (m *UserRepository) GetUserByEmail(_ context.Context, email string) (*user.User, error) {
	args := m.Called(email)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *UserRepository) GetUserByEmailAndPassword(_ context.Context, email, password string) (*user.User, error) {
	args := m.Called(email, password)
	return args.Get(0).(*user.User), args.Error(1)
}
