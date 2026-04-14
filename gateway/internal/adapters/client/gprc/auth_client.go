package gprcauth

import (
	"context"

	pb "github.com/IvanKuchsh-600/proto"

	"google.golang.org/grpc"
)

type AuthClient struct {
	client pb.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthClient(addr string) (*AuthClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := pb.NewAuthServiceClient(conn)
	return &AuthClient{client: client, conn: conn}, nil
}

func (c *AuthClient) Close() {
	c.conn.Close()
}

func (c *AuthClient) Register(email, password, role string) (string, error) {
	resp, err := c.client.Register(context.Background(), &pb.RegisterRequest{
		Email:    email,
		Password: password,
		Role:     role,
	})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (c *AuthClient) Login(email, password string) (string, error) {
	resp, err := c.client.Login(context.Background(), &pb.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (c *AuthClient) ValidateToken(token string) (bool, string, int64) {
	resp, err := c.client.ValidateToken(context.Background(), &pb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return false, "", 0
	}
	return resp.Valid, resp.Role, resp.UserId
}
