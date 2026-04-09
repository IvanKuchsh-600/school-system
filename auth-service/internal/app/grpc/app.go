package grpcapp

import (
	"auth-service/internal/adatpers/hasher"
	"auth-service/internal/adatpers/jwt"
	"auth-service/internal/adatpers/repository"
	"auth-service/internal/ports/grpc/auth"
	"auth-service/internal/usecases"
	"fmt"
	"log"
	"net"

	pb "github.com/IvanKuchsh-600/proto"
	"google.golang.org/grpc"
)

type App struct{}

func NewApp() (*App, error) {
	return &App{}, nil
}

func (a *App) Run() error {
	userRepo := repository.NewInMemoryUserRepo()
	jwtManager, err := jwt.NewJWTManager("my-secret-key", 8)
	if err != nil {
		return fmt.Errorf("initialize JWT manager: %w", err)
	}
	hasher := hasher.NewBcryptHasher()

	authService, err := usecases.NewAuthUseCase(userRepo, jwtManager, hasher)
	if err != nil {
		return fmt.Errorf("initialize auth use case: %w", err)
	}
	grpcHandler := auth.NewGrpcServer(authService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, grpcHandler)

	log.Println("Auth service running on :50051")
	err = grpcServer.Serve(lis)
	if err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}
