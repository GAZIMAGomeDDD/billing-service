package main

import (
	"context"
	"os"

	"github.com/GAZIMAGomeDDD/billing-service/internal/currency"
	"github.com/GAZIMAGomeDDD/billing-service/internal/handler"
	"github.com/GAZIMAGomeDDD/billing-service/internal/server"
	"github.com/GAZIMAGomeDDD/billing-service/internal/storage/postgres"
	"github.com/GAZIMAGomeDDD/billing-service/pkg/database/inmemory"
	"github.com/GAZIMAGomeDDD/billing-service/pkg/database/postgresdb"
	"github.com/sirupsen/logrus"
)

func main() {
	connString := os.Getenv("postgres_connection_string")
	serveAddress := os.Getenv("serve_address")
	exchangeratesapiToken := os.Getenv("exchangeratesapi_token")

	db, err := postgresdb.NewDB(context.Background(), connString)
	if err != nil {
		logrus.Fatal(err)
	}

	s, err := postgres.New(db)
	if err != nil {
		logrus.Fatal(err)
	}

	cache := inmemory.New()
	cr := currency.New(exchangeratesapiToken, cache, 10)
	h := handler.New(s, cr)

	srv := server.NewServer(h.Init(), serveAddress)
	if err := srv.Run(); err != nil {
		logrus.Fatal(err)
	}
}
