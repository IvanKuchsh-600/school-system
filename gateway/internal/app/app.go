package app

import (
	"fmt"
	gprcauth "gateway/internal/adapters/client/gprc"
	"gateway/internal/ports/http/gateway"
)

type App struct {
}

func NewApp() (*App, error) {
	return &App{}, nil
}

func (a *App) Run() error {
	authClient, err := gprcauth.NewAuthClient("localhost:50051")
	if err != nil {
		return fmt.Errorf("Failed to create auth client: %v", err)
	}

	server, err := gateway.NewServer(":8080", authClient)
	if err != nil {
		return fmt.Errorf("Failed to create server: %v", err)
	}

	err = server.Run()
	if err != nil {
		return fmt.Errorf("Failed to run server: %v", err)
	}

	return nil
}
