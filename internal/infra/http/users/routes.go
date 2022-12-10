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

// createUser godoc
// @Description Creates a single user.
// @Param first_name query string false "User's first name"
// @Param last_name query string false "User's last name"
// @Param email query string false "User's email"
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
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}
