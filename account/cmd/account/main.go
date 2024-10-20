package main

import (
	"log"
	"time"

	"github.com/coderkhushal/go-grpc-graphql-microservices/account"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
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

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {

		r, err = account.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return err
	})
	defer r.Close()

	// starting grpc server
	log.Println("Server Started on port 8080...")

	s := account.NewService(r)
	log.Fatal(account.ListenGRPC(s, 8080))
}
