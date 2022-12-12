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
		"empty first name should return error": {
			input:         user.User{},
			expectedError: fmt.Errorf("user's first name should not be empty"),
		},
		"empty last name should return error": {
			input: user.User{
				FirstName: "Dua",
			},
			expectedError: fmt.Errorf("user's last name should not be empty"),
		},
		"empty email should return error": {
			input: user.User{
				FirstName: "Dua",
				LastName:  "Lipa",
			},
			expectedError: fmt.Errorf("user's email should not be empty"),
		},
		"empty password should return error": {
			input: user.User{
				FirstName: "Dua",
				LastName:  "Lipa",
				Email:     "dua@lipa.com",
			},
			expectedError: fmt.Errorf("user's password should not be empty"),
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

func TestLogin(t *testing.T) {
	params1 := struct {
		email    string
		password string
	}{
		email:    "dua@lipa.com",
		password: "dua123lipa",
	}

	user1 := &user.User{
		FirstName: "Dua",
		LastName:  "Lipa",
		Email:     "dua@lipa.com",
		Password:  "dua123lipa",
	}

	testCases := map[string]struct {
		email         string
		password      string
		userRepoCalls []tmock.Call
		expectedError error
	}{
		"success": {
			email:    params1.email,
			password: params1.password,
			userRepoCalls: []tmock.Call{
				{
					FunctionName: "GetUserByEmailAndPassword",
					Params: []interface{}{
						params1.email,
						params1.password,
					},
					Returns: []interface{}{
						user1,
						nil,
					},
				},
			},
		},
		"empty email should return error": {
			expectedError: fmt.Errorf("user's email should not be empty"),
		},
		"empty password should return error": {
			email:         params1.email,
			expectedError: fmt.Errorf("user's password should not be empty"),
		},
		"user repo error should propagate": {
			email:    params1.email,
			password: params1.password,
			userRepoCalls: []tmock.Call{
				{
					FunctionName: "GetUserByEmailAndPassword",
					Params: []interface{}{
						params1.email,
						params1.password,
					},
					Returns: []interface{}{
						(*user.User)(nil),
						fmt.Errorf("user repo error"),
					},
				},
			},
			expectedError: fmt.Errorf("user repo error"),
		},
		"user not found should return error": {
			email:    params1.email,
			password: params1.password,
			userRepoCalls: []tmock.Call{
				{
					FunctionName: "GetUserByEmailAndPassword",
					Params: []interface{}{
						params1.email,
						params1.password,
					},
					Returns: []interface{}{
						(*user.User)(nil),
						nil,
					},
				},
			},
			expectedError: fmt.Errorf("user not exists or invalid credentials"),
		},
	}

	for name, tc := range testCases {
		ctx := context.Background()
		pldService := tmock.NewPLDService()
		userRepo := tmock.NewUserRepository().AddCall(t, tc.userRepoCalls)

		userService, err := user.NewService(pldService, userRepo)
		assert.NoError(t, err)

		email := tc.email
		password := tc.password
		expectedError := tc.expectedError

		t.Run(name, func(t *testing.T) {
			got, err := userService.Login(ctx, email, password)

			if expectedError == nil {
				assert.NotNil(t, got)
				assert.NoError(t, err)
			} else {
				assert.Equal(t, expectedError, err)
				assert.Nil(t, got)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	email1 := "dua@lipa.com"

	user1 := &user.User{
		FirstName: "Dua",
		LastName:  "Lipa",
		Email:     "dua@lipa.com",
		Password:  "dua123lipa",
	}

	testCases := map[string]struct {
		email         string
		userRepoCalls []tmock.Call
		expectedError error
	}{
		"success": {
			email: email1,
			userRepoCalls: []tmock.Call{
				{
					FunctionName: "GetUserByEmail",
					Params: []interface{}{
						email1,
					},
					Returns: []interface{}{
						user1,
						nil,
					},
				},
			},
		},
		"empty email should return error": {
			expectedError: fmt.Errorf("user's email should not be empty"),
		},
		"user repo error should propagate": {
			email: email1,
			userRepoCalls: []tmock.Call{
				{
					FunctionName: "GetUserByEmail",
					Params: []interface{}{
						email1,
					},
					Returns: []interface{}{
						(*user.User)(nil),
						fmt.Errorf("user repo error"),
					},
				},
			},
			expectedError: fmt.Errorf("user repo error"),
		},
		"user not found should return error": {
			email: email1,
			userRepoCalls: []tmock.Call{
				{
					FunctionName: "GetUserByEmail",
					Params: []interface{}{
						email1,
					},
					Returns: []interface{}{
						(*user.User)(nil),
						nil,
					},
				},
			},
			expectedError: fmt.Errorf("user not exists"),
		},
	}

	for name, tc := range testCases {
		ctx := context.Background()
		pldService := tmock.NewPLDService()
		userRepo := tmock.NewUserRepository().AddCall(t, tc.userRepoCalls)

		userService, err := user.NewService(pldService, userRepo)
		assert.NoError(t, err)

		email := tc.email
		expectedError := tc.expectedError

		t.Run(name, func(t *testing.T) {
			got, err := userService.GetUser(ctx, email)

			if expectedError == nil {
				assert.NotNil(t, got)
				assert.NoError(t, err)
			} else {
				assert.Equal(t, expectedError, err)
				assert.Nil(t, got)
			}
		})
	}
}
