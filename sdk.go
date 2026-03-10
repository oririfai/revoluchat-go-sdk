package revoluchat

import (
	"context"
	"fmt"
	"net"

	pb "github.com/oririfai/revoluchat-go-sdk/proto/user_v1"
	"google.golang.org/grpc"
)

// User represents the user data structure required by Revoluchat.
type User struct {
	ID        uint64
	UUID      string
	Name      string
	Phone     string
	Status    string
	IsKYC     bool
	AvatarURL string
}

// UserProvider is a function that returns a User given an ID.
type UserProvider func(ctx context.Context, id uint64) (*User, error)

type server struct {
	pb.UnimplementedUserServiceServer
	provider UserProvider
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.provider(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		Id:        user.ID,
		Uuid:      user.UUID,
		Name:      user.Name,
		Phone:     user.Phone,
		Status:    user.Status,
		IsKyc:     user.IsKYC,
		AvatarUrl: user.AvatarURL,
	}, nil
}

// Config holds the SDK configuration.
type Config struct {
	GRPCPort int
	Provider UserProvider
}

// Start starts the gRPC server for Revoluchat integration.
func Start(config Config) error {
	addr := fmt.Sprintf(":%d", config.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{provider: config.Provider})

	fmt.Printf("Revoluchat Go SDK: gRPC server listening on %s\n", addr)
	return s.Serve(lis)
}
