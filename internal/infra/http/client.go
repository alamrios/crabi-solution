package http

import (
	"bytes"
	"net/http"
)

// Client struct for http
type Client struct{}

// NewClient constructor for htttp client
func NewClient() (*Client, error) {
	return &Client{}, nil
}

// Post send a post request
func (c *Client) Post(URL string, data []byte, contentType string) (*http.Response, error) {
	bytesBuffer := bytes.NewBuffer(data)
	response, err := http.Post(URL, contentType, bytesBuffer)
	if err != nil {
		return nil, err
	}

	return response, nil
}
