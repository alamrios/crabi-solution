package pld

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/alamrios/crabi-solution/config"
	"github.com/alamrios/crabi-solution/internal/app/pld"
	pb "github.com/alamrios/crabi-solution/proto/pld"
)

// Service struct for PLD service
type Service struct {
	URL string
}

// NewService PLD service constructor
func NewService(cfg *config.PLD) (*Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("pld config is nil")
	}

	return &Service{
		URL: cfg.Host + ":" + cfg.Port,
	}, nil
}

// CheckBlacklist goes to PLD Service to ckeck if data is in black list
// Returns error if user found in pld blacklist, nil otherwise
func (s *Service) CheckBlacklist(ctx context.Context, request pld.Request) error {
	conn, err := grpc.Dial(s.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("failed to connect to pld service:", err.Error())
		return fmt.Errorf("failed to connect to pld service")
	}
	defer conn.Close()

	client := pb.NewPldServiceClient(conn)
	requestBody := &pb.CheckInBlackListReq{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Birthday:  "",
	}
	res, err := client.CheckInBlacklist(context.Background(), requestBody)
	if err != nil {
		log.Println("failed to check in pld service:", err.Error())
		return fmt.Errorf("failed to check in pld service")
	}

	if res.IsInBlacklist {
		return fmt.Errorf("user was found in pld blacklist")
	}

	return nil
}
