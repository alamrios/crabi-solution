package pld

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alamrios/crabi-solution/config"
	"github.com/alamrios/crabi-solution/internal/app/pld"
)

type HttpClient interface {
	Post(URL string, data []byte, contentType string) (*http.Response, error)
}

// Service struct for PLD service
type Service struct {
	URL        string
	httpClient HttpClient
}

// NewService PLD service constructor
func NewService(cfg *config.PLD, httpClient HttpClient) (*Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("pld config is nil")
	}

	if httpClient == nil {
		return nil, fmt.Errorf("http client is nil")
	}

	return &Service{
		URL:        cfg.Protocol + cfg.Host + ":" + cfg.Port + cfg.URI,
		httpClient: httpClient,
	}, nil
}

// CheckBlacklistRequest struct for CheckBlacklist request
type CheckBlacklistRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// CheckBlacklistResponse struct for CheckBlacklist response
type CheckBlacklistResponse struct {
	IsInBlacklist bool `json:"is_in_blacklist"`
}

// CheckBlacklist goes to PLD Service to ckeck if data is in black list
// Returns error if user found in pld blacklist, nil otherwise
func (s *Service) CheckBlacklist(ctx context.Context, request pld.Request) error {
	requestBody := CheckBlacklistRequest{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
	}

	data, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	response, err := s.httpClient.Post(s.URL, data, "application/json")
	if err != nil {
		return err
	}

	fmt.Println(response)

	if response.StatusCode != 201 {
		return fmt.Errorf("pld server returned %d status code", response.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var responseObject CheckBlacklistResponse
	json.Unmarshal(bodyBytes, &responseObject)

	if responseObject.IsInBlacklist {
		return fmt.Errorf("user was found in pld blacklist")
	}

	return nil
}
