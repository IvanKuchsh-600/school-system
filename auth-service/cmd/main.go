package main

import (
	"auth-service/deployment/config"
	grpcapp "auth-service/internal/app/grpc"
	"log"
)

func main() {
	cfg := config.MustLoad()

	app, err := grpcapp.NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
