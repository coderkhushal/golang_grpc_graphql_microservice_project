package main

import (
	"log"
	"time"

	"github.com/coderkhushal/go-grpc-graphql-microservices/catalog"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	var r catalog.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = catalog.NewElasticRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return err

	})
	defer r.Close()
	log.Println("listening on port 8080...")

	s := catalog.NewCatalogService(r)

	log.Fatal(catalog.ListenGRPC(s, 8080))
}
