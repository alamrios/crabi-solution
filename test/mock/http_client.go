package mock

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
)

// HttpClient is a mock of Service interface.
type HttpClient struct {
	mock.Mock
}

// NewHttpClient creates new pld mock service
func NewHttpClient() *HttpClient {
	return &HttpClient{}
}

// AddCall adds new call to the mock
func (m *HttpClient) AddCall(t *testing.T, calls []Call) *HttpClient {
	t.Helper()

	for _, call := range calls {
		m.On(call.FunctionName, call.Params...).Return(call.Returns...)
	}

	return m
}

// Post method mock
func (m *HttpClient) Post(URL string, _ []byte, contentType string) (*http.Response, error) {
	args := m.Called(URL, contentType)
	return args.Get(0).(*http.Response), args.Error(1)
}
