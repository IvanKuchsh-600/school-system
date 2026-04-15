package gateway

import (
	"errors"
	"fmt"
	gprcauth "gateway/internal/adapters/client/gprc"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type Server struct {
	port       string
	authClient *gprcauth.AuthClient
	r          *gin.Engine
}

func NewServer(port string, authClient *gprcauth.AuthClient) (*Server, error) {
	if port == "" {
		return nil, errors.New("port is required")
	}

	if authClient == nil {
		return nil, errors.New("auth client is nil")
	}

	r := gin.Default()

	return &Server{port: port, authClient: authClient, r: r}, nil
}

func (s *Server) Run() error {
	s.r.POST("/register", s.Register)
	s.r.POST("/login", s.Login)

	err := s.r.Run(":8080")
	if err != nil {
		return fmt.Errorf("could not start server: %w", err)
	}

	return nil
}

func (s *Server) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := s.authClient.Register(req.Email, req.Password, req.Role)
	if err != nil {
		handleGrpcError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *Server) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := s.authClient.Login(req.Email, req.Password)
	if err != nil {
		handleGrpcError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func handleGrpcError(c *gin.Context, err error) {
	grpcStatus, ok := status.FromError(err)
	if ok {
		switch grpcStatus.Code() {
		// 400 Bad Request - клиентские ошибки
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{"error": grpcStatus.Message()})
			return

		// 401 Unauthorized - неверные учетные данные
		case codes.Unauthenticated:
			c.JSON(http.StatusUnauthorized, gin.H{"error": grpcStatus.Message()})
			return

		// 403 Forbidden - недостаточно прав
		case codes.PermissionDenied:
			c.JSON(http.StatusForbidden, gin.H{"error": grpcStatus.Message()})
			return

		// 404 Not Found - ресурс не найден
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": grpcStatus.Message()})
			return

		// 409 Conflict - пользователь уже существует
		case codes.AlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": grpcStatus.Message()})
			return

		// 500 Internal Server Error - внутренние ошибки
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return

		// 503 Unavailable - сервис недоступен
		case codes.Unavailable:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "service unavailable"})
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
			return
		}
	}

	// Если не удалось распарсить gRPC статус
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
