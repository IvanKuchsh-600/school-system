package handler

import (
	"auth-service/internal/entities"
	"auth-service/internal/usecases"
	"context"
	"errors"

	pb "github.com/IvanKuchsh-600/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GrpcHandler - адаптер, который преобразует gRPC запросы в вызовы use cases
type GrpcHandler struct {
	pb.UnimplementedAuthServiceServer
	authUseCase *usecases.AuthUseCase
}

// NewGrpcHandler - конструктор адаптера
func NewGrpcHandler(authUseCase *usecases.AuthUseCase) *GrpcHandler {
	return &GrpcHandler{
		authUseCase: authUseCase,
	}
}

// RegisterAdmin - обработчик gRPC запроса
func (h *GrpcHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.Role == "" {
		return nil, status.Error(codes.InvalidArgument, "role is required")
	}

	token, err := h.authUseCase.Register(req.Email, req.Password, req.Role)

	if err != nil {
		return nil, mapDomainErrorToGrpc(err)
	}

	return &pb.AuthResponse{
		Token:   token,
		Message: "registered successfully",
	}, nil
}

// Login - обработчик gRPC запроса
func (h *GrpcHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	token, err := h.authUseCase.Login(req.Email, req.Password)
	if err != nil {
		return &pb.AuthResponse{
			Token:   "",
			Message: err.Error(),
		}, nil
	}

	return &pb.AuthResponse{
		Token:   token,
		Message: "login successful",
	}, nil
}

// ValidateToken - обработчик gRPC запроса
func (h *GrpcHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := h.authUseCase.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid:  false,
			Role:   "",
			UserId: 0,
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:  true,
		Role:   claims.Role,
		UserId: claims.UserID,
	}, nil
}

// mapDomainErrorToGrpc преобразует доменные ошибки в gRPC статусы
func mapDomainErrorToGrpc(err error) error {
	// клиентские ошибки - 400
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

	// внутренние ошибки - 500
	case errors.Is(err, entities.ErrDatabaseOperation):
		return status.Error(codes.Internal, "internal server error, please try again later")

	case errors.Is(err, entities.ErrHashFailed):
		return status.Error(codes.Internal, "internal server error, please try again later")

	default:
		return status.Error(codes.Internal, "unexpected error occurred")
	}
}
