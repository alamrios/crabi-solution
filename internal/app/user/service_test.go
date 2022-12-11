package user_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alamrios/crabi-solution/internal/app/pld"
	"github.com/alamrios/crabi-solution/internal/app/user"
	tmock "github.com/alamrios/crabi-solution/test/mock"
)

func TestNewService(t *testing.T) {
	testCases := map[string]struct {
		pldService pld.Service
		userRepo   user.Repository
		err        string
	}{
		"success": {
			pldService: &tmock.PLDService{},
			userRepo:   &tmock.UserRepository{},
		},
		"missing pld service": {
			err: "pld service is nil",
		},
		"missing user repo": {
			pldService: &tmock.PLDService{},
			err:        "user repo is nil",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := user.NewService(tc.pldService, tc.userRepo)
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

func TestCreateUser(t *testing.T) {
	input1 := user.User{
		FirstName: "Dua",
		LastName:  "Lipa",
		Email:     "dua@lipa.com",
		Password:  "dua123lipa",
	}

	testCases := map[string]struct {
		input           user.User
		pldServiceCalls []tmock.Call
		userRepoCalls   []tmock.Call
		expectedResult  *user.User
		expectedError   error
	}{
		"success": {
			input: input1,
			pldServiceCalls: []tmock.Call{
				{
					FunctionName: "CheckBlacklist",
					Params: []interface{}{
						pld.Request{
							FirstName: input1.FirstName,
							LastName:  input1.LastName,
							Email:     input1.Email,
						},
					},
					Returns: []interface{}{
						nil,
					},
				},
			},
			userRepoCalls: []tmock.Call{
				{
					FunctionName: "GetUserByEmail",
					Params: []interface{}{
						input1.Email,
					},
					Returns: []interface{}{
						(*user.User)(nil),
						nil,
					},
				},
				{
					FunctionName: "SaveUser",
					Params: []interface{}{
						input1,
					},
					Returns: []interface{}{
						nil,
					},
				},
			},
			expectedResult: &input1,
			expectedError:  nil,
		},
		"pld error should propagate": {
			input: input1,
			pldServiceCalls: []tmock.Call{
				{
					FunctionName: "CheckBlacklist",
					Params: []interface{}{
						pld.Request{
							FirstName: input1.FirstName,
							LastName:  input1.LastName,
							Email:     input1.Email,
						},
					},
					Returns: []interface{}{
						fmt.Errorf("user found in pld blacklist"),
					},
				},
			},
			expectedError: fmt.Errorf("user found in pld blacklist"),
		},
		"user repository error while GetUserByEmail should propagate": {
			input: input1,
			pldServiceCalls: []tmock.Call{
				{
					FunctionName: "CheckBlacklist",
					Params: []interface{}{
						pld.Request{
							FirstName: input1.FirstName,
							LastName:  input1.LastName,
							Email:     input1.Email,
						},
					},
					Returns: []interface{}{
						nil,
					},
				},
			},
			userRepoCalls: []tmock.Call{
				{
					FunctionName: "GetUserByEmail",
					Params: []interface{}{
						input1.Email,
					},
					Returns: []interface{}{
						(*user.User)(nil),
						fmt.Errorf("user repo error"),
					},
				},
			},
			expectedError: fmt.Errorf("user repo error"),
		},
		"duplicated email should return error": {
			input: input1,
			pldServiceCalls: []tmock.Call{
				{
					FunctionName: "CheckBlacklist",
					Params: []interface{}{
						pld.Request{
							FirstName: input1.FirstName,
							LastName:  input1.LastName,
							Email:     input1.Email,
						},
					},
					Returns: []interface{}{
						nil,
					},
				},
			},
			userRepoCalls: []tmock.Call{
				{
					FunctionName: "GetUserByEmail",
					Params: []interface{}{
						input1.Email,
					},
					Returns: []interface{}{
						&input1,
						nil,
					},
				},
			},
			expectedError: fmt.Errorf("user with email %s already exists", input1.Email),
		},
		"user repository error while SaveUser should propagate": {
			input: input1,
			pldServiceCalls: []tmock.Call{
				{
					FunctionName: "CheckBlacklist",
					Params: []interface{}{
						pld.Request{
							FirstName: input1.FirstName,
							LastName:  input1.LastName,
							Email:     input1.Email,
						},
					},
					Returns: []interface{}{
						nil,
					},
				},
			},
			userRepoCalls: []tmock.Call{
				{
					FunctionName: "GetUserByEmail",
					Params: []interface{}{
						input1.Email,
					},
					Returns: []interface{}{
						(*user.User)(nil),
						nil,
					},
				},
				{
					FunctionName: "SaveUser",
					Params: []interface{}{
						input1,
					},
					Returns: []interface{}{
						fmt.Errorf("user repository error"),
					},
				},
			},
			expectedError: fmt.Errorf("user repository error"),
		},
	}

	for name, tc := range testCases {
		ctx := context.Background()
		pldService := tmock.NewPLDService().AddCall(t, tc.pldServiceCalls)
		userRepo := tmock.NewUserRepository().AddCall(t, tc.userRepoCalls)

		userService, err := user.NewService(pldService, userRepo)
		assert.NoError(t, err)

		input := tc.input
		expectedResult := tc.expectedResult
		expectedError := tc.expectedError

		t.Run(name, func(t *testing.T) {
			got, err := userService.CreateUser(ctx, input)
			assert.Equal(t, expectedResult, got)
			assert.Equal(t, expectedError, err)
		})
	}
}
