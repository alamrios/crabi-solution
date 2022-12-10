package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/alamrios/crabi-solution/config"
	"github.com/alamrios/crabi-solution/internal/app/user"
	"github.com/alamrios/crabi-solution/internal/infra/http/pld"
	userRouter "github.com/alamrios/crabi-solution/internal/infra/http/users"
	"github.com/alamrios/crabi-solution/internal/infra/repository/mongo"
	userRepo "github.com/alamrios/crabi-solution/internal/infra/repository/mongo/user"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New(ctx)
	if err != nil {
		fmt.Printf("failed to set up configuration: %v", err)
		os.Exit(1)
	}

	mongoDB, err := mongo.NewClient(ctx, &cfg.Mongo)
	if err != nil {
		log.Fatalf("failed to setup mongoDB client: %v", err)
	}

	pldService, err := pld.NewService()
	if err != nil {
		log.Fatalf("failed to setup pld client: %v", err)
	}

	userRepo, err := userRepo.New(mongoDB)
	if err != nil {
		log.Fatalf("failed to setup user repo: %v", err)
	}

	userService, err := user.NewService(pldService, userRepo)
	if err != nil {
		log.Fatalf("failed to setup user service: %v", err)
	}

	usersRouter, err := userRouter.New(userService)
	if err != nil {
		log.Fatalf("failed to setup user router: %v", err)
	}

	router := mux.NewRouter()
	usersRouter.AppendRoutes(router)

	log.Print("Starting crabi-solution server...")
	err = http.ListenAndServe(":8080", router)
	log.Fatal(err)
}
