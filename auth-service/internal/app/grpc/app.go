package grpcapp

import (
	"auth-service/deployment/config"
	"auth-service/internal/adatpers/hasher"
	"auth-service/internal/adatpers/jwt"
	"auth-service/internal/adatpers/repository"
	"auth-service/internal/ports/grpc/auth"
	"auth-service/internal/usecases"
	"errors"
	"fmt"
	"log"
	"net"

	pb "github.com/IvanKuchsh-600/proto"
	"google.golang.org/grpc"
)

type App struct {
	cfg *config.Config
}

func NewApp(cfg *config.Config) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("create app: %w", errors.New("config is required"))
	}

	return &App{cfg: cfg}, nil
}

func (a *App) Run() error {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Println(cfg.JWTSecret)
	fmt.Println(cfg.JWTExpirationHours)
	fmt.Println(cfg.ServerPort)
	userRepo := repository.NewInMemoryUserRepo()

	jwtManager, err := jwt.NewJWTManager(cfg.JWTSecret, cfg.JWTExpirationHours)
	if err != nil {
		return fmt.Errorf("initialize JWT manager: %w", err)
	}
	hasher := hasher.NewBcryptHasher()

	authService, err := usecases.NewAuthUseCase(userRepo, jwtManager, hasher)
	if err != nil {
		return fmt.Errorf("initialize auth use case: %w", err)
	}
	grpcHandler := auth.NewGrpcServer(authService)

	lis, err := net.Listen("tcp", cfg.ServerPort)
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
