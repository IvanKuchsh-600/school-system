package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	handler "auth-service/internal/adatpers/handler"
	"auth-service/internal/adatpers/hasher"
	"auth-service/internal/adatpers/jwt"
	"auth-service/internal/adatpers/repository"
	"auth-service/internal/usecases"

	pb "github.com/IvanKuchsh-600/proto"
)

func main() {
	userRepo := repository.NewInMemoryUserRepo()
	jwtManager := jwt.NewJWTManager("my-secret-key", 24)
	hasher := hasher.NewBcryptHasher()

	authService := usecases.NewAuthUseCase(userRepo, jwtManager, hasher)

	grpcHandler := handler.NewGrpcHandler(authService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, grpcHandler)

	log.Println("Auth service running on :50051")
	log.Fatal(grpcServer.Serve(lis))
}
