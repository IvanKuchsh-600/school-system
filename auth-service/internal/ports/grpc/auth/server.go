package auth

import (
	"auth-service/internal/entities"
	"context"
	"errors"
	"fmt"

	pb "github.com/IvanKuchsh-600/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	pb.UnimplementedAuthServiceServer
	authService AuthService
}

func NewGrpcServer(authService AuthService) *GrpcServer {
	return &GrpcServer{
		authService: authService,
	}
}

func (gs *GrpcServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.Role == "" {
		return nil, status.Error(codes.InvalidArgument, "role is required")
	}

	token, err := gs.authService.Register(req.Email, req.Password, req.Role)

	if err != nil {
		return nil, mapEntitiesErrorToGrpc(err)
	}

	return &pb.AuthResponse{
		Token:   token,
		Message: "registered successfully",
	}, nil
}

func (gs *GrpcServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	token, err := gs.authService.Login(req.Email, req.Password)
	if err != nil {
		return nil, mapEntitiesErrorToGrpc(err)
	}

	return &pb.AuthResponse{
		Token:   token,
		Message: "login successful",
	}, nil
}

func mapEntitiesErrorToGrpc(err error) error {
	fmt.Println(err)
	// клиентские ошибки
	switch {
	case errors.Is(err, entities.ErrEmailRequired):
		return status.Error(codes.InvalidArgument, "email is required")
	case errors.Is(err, entities.ErrEmailInvalid):
		return status.Error(codes.InvalidArgument, "invalid email format")

	case errors.Is(err, entities.ErrPasswordRequired):
		return status.Error(codes.InvalidArgument, "password is required")
	case errors.Is(err, entities.ErrRoleRequired):
		return status.Error(codes.InvalidArgument, "role is required")
	case errors.Is(err, entities.ErrRoleInvalid):
		return status.Error(codes.InvalidArgument, "invalid role")

	// бизнес-ошибки
	case errors.Is(err, entities.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, "user with this email already exists")
	case errors.Is(err, entities.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, "invalid email or password")

	// внутренние ошибки - 500
	case errors.Is(err, entities.ErrDatabaseOperation):
		return status.Error(codes.Internal, "internal server error, please try again later")

	case errors.Is(err, entities.ErrHashFailed):
		return status.Error(codes.Internal, "internal server error, please try again later")

	case errors.Is(err, entities.ErrInternalError):
		return status.Error(codes.Internal, "internal server error, please try again later")

	default:
		return status.Error(codes.Internal, "unexpected error occurred")
	}
}
