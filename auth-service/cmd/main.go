package main

import (
	grpcapp "auth-service/internal/app/grpc"
	"log"
)

func main() {
	app, err := grpcapp.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
