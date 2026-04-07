package main

import (
	"gateway/internal/client"
	"gateway/internal/handler"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Подключаемся к gRPC серверу
	authClient, err := client.NewAuthClient("localhost:50051")
	if err != nil {
		log.Fatal(err)
	}

	// Gin роутер
	r := gin.Default()

	// Публичные маршруты (без токена)
	authHandler := handler.NewAuthHandler(authClient)
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	r.Run(":8080")
}
