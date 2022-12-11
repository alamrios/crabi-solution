// Source: internal/app/pld/service.go
package mock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/alamrios/crabi-solution/internal/app/pld"
)

// PLDService is a mock of Service interface.
type PLDService struct {
	mock.Mock
}

// NewPLDService creates new pld mock service
func NewPLDService() *PLDService {
	return &PLDService{}
}

// AddCall adds new call to the mock
func (m *PLDService) AddCall(t *testing.T, calls []Call) *PLDService {
	t.Helper()

	for _, call := range calls {
		m.On(call.FunctionName, call.Params...).Return(call.Returns...)
	}

	return m
}

// CheckBlacklist method mock
func (m *PLDService) CheckBlacklist(_ context.Context, request pld.Request) error {
	args := m.Called(request)
	return args.Error(0)
}
