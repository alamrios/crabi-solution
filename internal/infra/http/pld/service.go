package pld

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alamrios/crabi-solution/config"
	"github.com/alamrios/crabi-solution/internal/app/pld"
)

// Service struct for PLD service
type Service struct {
	URL string
}

// NewService PLD service constructor
func NewService(cfg *config.PLD) (*Service, error) {
	service := &Service{
		URL: cfg.Protocol + cfg.Host + ":" + cfg.Port + cfg.URI,
	}

	return service, nil
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
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	contentType := "application/json"
	bytesBuffer := bytes.NewBuffer(jsonData)
	response, err := http.Post(s.URL, contentType, bytesBuffer)
	if err != nil {
		return err
	}

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
