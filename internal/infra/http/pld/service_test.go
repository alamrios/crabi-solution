package pld_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alamrios/crabi-solution/config"
	model "github.com/alamrios/crabi-solution/internal/app/pld"
	"github.com/alamrios/crabi-solution/internal/infra/http/pld"
	tmock "github.com/alamrios/crabi-solution/test/mock"
)

func TestNewService(t *testing.T) {
	testCases := map[string]struct {
		cfg        *config.PLD
		httpClient *tmock.HttpClient
		err        string
	}{
		"success": {
			cfg:        &config.PLD{},
			httpClient: &tmock.HttpClient{},
		},
		"nil config should return error": {
			httpClient: &tmock.HttpClient{},
			err:        "pld config is nil",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := pld.NewService(tc.cfg, tc.httpClient)
			if tc.err != "" {
				assert.EqualError(t, err, tc.err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func TestCheckBlacklist(t *testing.T) {
	ctx := context.Background()
	cfg := config.PLD{
		Protocol: "http://",
		Host:     "crabi-pld",
		Port:     "3000",
		URI:      "/check-blacklist",
	}

	request1 := model.Request{
		FirstName: "Dua",
		LastName:  "Lipa",
		Email:     "dua@lipa.com",
	}

	response1 := pld.CheckBlacklistResponse{
		IsInBlacklist: false,
	}

	data1, _ := json.Marshal(response1)
	reader1 := io.NopCloser(bytes.NewReader(data1))

	request2 := model.Request{
		FirstName: "Joaquin",
		LastName:  "Guzman",
		Email:     "joaquin@guzman.com",
	}

	response2 := pld.CheckBlacklistResponse{
		IsInBlacklist: true,
	}

	data2, _ := json.Marshal(response2)
	reader2 := io.NopCloser(bytes.NewReader(data2))

	tests := map[string]struct {
		httpClientCalls []tmock.Call
		request         model.Request
		err             error
	}{
		"success": {
			httpClientCalls: []tmock.Call{
				{
					FunctionName: "Post",
					Params: []interface{}{
						cfg.Protocol + cfg.Host + ":" + cfg.Port + cfg.URI,
						"application/json",
					},
					Returns: []interface{}{
						&http.Response{
							StatusCode: 201,
							Body:       reader1,
						},
						nil,
					},
				},
			},
			request: request1,
		},
		"data found in pld service should return error": {
			httpClientCalls: []tmock.Call{
				{
					FunctionName: "Post",
					Params: []interface{}{
						cfg.Protocol + cfg.Host + ":" + cfg.Port + cfg.URI,
						"application/json",
					},
					Returns: []interface{}{
						&http.Response{
							StatusCode: 201,
							Body:       reader2,
						},
						nil,
					},
				},
			},
			request: request2,
			err:     fmt.Errorf("user was found in pld blacklist"),
		},
		"http client error should propagate": {
			httpClientCalls: []tmock.Call{
				{
					FunctionName: "Post",
					Params: []interface{}{
						cfg.Protocol + cfg.Host + ":" + cfg.Port + cfg.URI,
						"application/json",
					},
					Returns: []interface{}{
						(*http.Response)(nil),
						fmt.Errorf("http client error"),
					},
				},
			},
			request: request1,
			err:     fmt.Errorf("http client error"),
		},
		"http status different to 201 should return error": {
			httpClientCalls: []tmock.Call{
				{
					FunctionName: "Post",
					Params: []interface{}{
						cfg.Protocol + cfg.Host + ":" + cfg.Port + cfg.URI,
						"application/json",
					},
					Returns: []interface{}{
						&http.Response{
							StatusCode: 400,
						},
						nil,
					},
				},
			},
			request: request1,
			err:     fmt.Errorf("pld server returned 400 status code"),
		},
	}
	for name, tt := range tests {

		request := tt.request
		expected := tt.err
		httpClientCalls := tt.httpClientCalls

		t.Run(name, func(t *testing.T) {
			httpClient := tmock.NewHttpClient().AddCall(t, httpClientCalls)
			pldService, err := pld.NewService(&cfg, httpClient)
			assert.NoError(t, err)

			got := pldService.CheckBlacklist(ctx, request)
			assert.Equal(t, got, expected)
		})
	}
}
