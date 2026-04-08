package grpcapp

import (
	"auth-service/internal/adatpers/hasher"
	"auth-service/internal/adatpers/jwt"
	"auth-service/internal/adatpers/repository"
	"auth-service/internal/ports/grpc/auth"
	"auth-service/internal/usecases"
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
	jwtManager := jwt.NewJWTManager("my-secret-key", 24)
	hasher := hasher.NewBcryptHasher()

	authService := usecases.NewAuthUseCase(userRepo, jwtManager, hasher)

	grpcHandler := auth.NewGrpcServer(authService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, grpcHandler)

	log.Println("Auth service running on :50051")
	log.Fatal(grpcServer.Serve(lis))

	return nil
}
