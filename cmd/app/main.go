package main

import (
	"context"
	"fmt"

	"github.com/GAZIMAGomeDDD/billing-service/internal/currencycache"
	"github.com/GAZIMAGomeDDD/billing-service/internal/handler"
	"github.com/GAZIMAGomeDDD/billing-service/internal/server"
	"github.com/GAZIMAGomeDDD/billing-service/internal/storage/inmemory"
	"github.com/GAZIMAGomeDDD/billing-service/internal/storage/postgres"
	"github.com/GAZIMAGomeDDD/billing-service/pkg/database/postgresdb"
)

func main() {
	connString := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	db, err := postgresdb.NewDB(context.Background(), connString)
	if err != nil {
		fmt.Println(err)
		return
	}
	s, err := postgres.New(db)
	if err != nil {
		fmt.Println(err)

		return
	}
	cache := inmemory.New()
	crCache := currencycache.New("", cache, 10)
	h := handler.New(s, crCache)
	srv := server.NewServer(h.Init())
	srv.Run()
}
