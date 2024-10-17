package main

import (
	"log"
	"time"

	"github.com/coderkhushal/go-grpc-graphql-microservices/account"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	// environment variables
	var cfg config

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// check if database is connected correctly
	time.Sleep(time.Second * 2)
	r, err := account.NewPostgresRepository(cfg.DatabaseURL)
	if err != nil {
		log.Println(err)
	}
	defer r.Close()

	// starting grpc server
	log.Println("Server Started on port 8080...")

	s := account.NewService(r)
	log.Fatal(account.ListenGRPC(s, 8080))
}
