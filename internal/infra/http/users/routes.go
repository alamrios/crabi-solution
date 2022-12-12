package users

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alamrios/crabi-solution/internal/app/user"
	"github.com/gorilla/mux"
)

type userService interface {
	CreateUser(ctx context.Context, user user.User) (*user.User, error)
	Login(ctx context.Context, email, password string) (*user.User, error)
	GetUser(ctx context.Context, email string) (*user.User, error)
}

// Router infraestructure
type Router struct {
	userService userService
}

// New Router constructor
func New(userService userService) (*Router, error) {
	if userService == nil {
		return nil, fmt.Errorf("user service is nil")
	}

	return &Router{
		userService: userService,
	}, nil
}

// AppendRoutes adds all func handlers
func (h *Router) AppendRoutes(rb *mux.Router) {
	rb.HandleFunc("/api/v1/users/", h.createUser).Methods("POST")
	rb.HandleFunc("/api/v1/users/{email}", h.getUser).Methods("GET")
	rb.HandleFunc("/api/v1/login/", h.login).Methods("POST")
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, DELETE, PUT, POST, OPTIONS")
}

type createUserRequest struct {
	FirstName string `json:"first_name" validate:"required" example:"Joaquin"`
	LastName  string `json:"last_name" validate:"required" example:"Guzman"`
	Email     string `json:"email" validate:"required" example:"joaquin@guzman.com"`
	Password  string `json:"password" validate:"required" example:"joaquin123"`
}

type createUserResponse struct {
	FirstName string `json:"first_name" validate:"required" example:"Joaquin"`
	LastName  string `json:"last_name" validate:"required" example:"Guzman"`
	Email     string `json:"email" validate:"required" example:"joaquin@guzman.com"`
}

// createUser godoc
// @Description Creates a single user.
// @Param first_name query string false "User's first name"
// @Param last_name query string false "User's last name"
// @Param email query string false "User's email"
// @Param password query string false "User's password"
// @Success 200 {object} jsonapi.Response{user.User}
func (h *Router) createUser(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	ctx := r.Context()

	var userToCreate createUserRequest
	var err error
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &userToCreate)

	user, err := h.userService.CreateUser(ctx, user.User{
		FirstName: userToCreate.FirstName,
		LastName:  userToCreate.LastName,
		Email:     userToCreate.Email,
		Password:  userToCreate.Password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		response := createUserResponse{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

type loginRequest struct {
	Email    string `json:"email" validate:"required" example:"joaquin@guzman.com"`
	Password string `json:"password" validate:"required" example:"joaquin123"`
}

type loginResponse struct {
	FirstName string `json:"first_name" validate:"required" example:"Joaquin"`
	LastName  string `json:"last_name" validate:"required" example:"Guzman"`
	Email     string `json:"email" validate:"required" example:"joaquin@guzman.com"`
}

// login godoc
// @Description User login.
// @Param email query string false "User's email"
// @Param password query string false "User's password"
// @Success 200 {object} jsonapi.Response{loginResponse}
func (h *Router) login(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	ctx := r.Context()

	var request loginRequest
	var err error
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &request)

	user, err := h.userService.Login(ctx, request.Email, request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		response := loginResponse{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

type getUserResponse struct {
	FirstName string `json:"first_name" validate:"required" example:"Joaquin"`
	LastName  string `json:"last_name" validate:"required" example:"Guzman"`
	Email     string `json:"email" validate:"required" example:"joaquin@guzman.com"`
}

// getUser godoc
// @Description retreive user data for given email.
// @Param email query string false "User's email"
// @Param password query string false "User's password"
// @Success 200 {object} jsonapi.Response{user.User}
func (h *Router) getUser(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	ctx := r.Context()

	vars := mux.Vars(r)
	email, ok := vars["email"]
	if !ok {
		http.Error(w, "user's email needed", http.StatusBadRequest)
	}

	user, err := h.userService.GetUser(ctx, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		response := getUserResponse{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
