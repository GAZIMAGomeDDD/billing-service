package main

import (
	"context"

	"github.com/GAZIMAGomeDDD/billing-service/internal/currency"
	"github.com/GAZIMAGomeDDD/billing-service/internal/handler"
	"github.com/GAZIMAGomeDDD/billing-service/internal/server"
	"github.com/GAZIMAGomeDDD/billing-service/internal/storage/postgres"
	"github.com/GAZIMAGomeDDD/billing-service/pkg/database/inmemory"
	"github.com/GAZIMAGomeDDD/billing-service/pkg/database/postgresdb"
	"github.com/sirupsen/logrus"
)

func main() {
	connString := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	db, err := postgresdb.NewDB(context.Background(), connString)
	if err != nil {
		logrus.Fatal(err)
	}

	s, err := postgres.New(db)
	if err != nil {
		logrus.Fatal(err)
	}

	cache := inmemory.New()
	cr := currency.New("207e0d99dc8df832c4921e5af54e56e4", cache, 100)
	h := handler.New(s, cr)

	srv := server.NewServer(h.Init(), ":8080")
	if err := srv.Run(); err != nil {
		logrus.Fatal(err)
	}
}
