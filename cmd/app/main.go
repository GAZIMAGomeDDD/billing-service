package main

import (
	"context"
	"fmt"

	"github.com/GAZIMAGomeDDD/billing-service/internal/currency"
	"github.com/GAZIMAGomeDDD/billing-service/internal/handler"
	"github.com/GAZIMAGomeDDD/billing-service/internal/server"
	"github.com/GAZIMAGomeDDD/billing-service/internal/storage/postgres"
	"github.com/GAZIMAGomeDDD/billing-service/pkg/database/inmemory"
	"github.com/GAZIMAGomeDDD/billing-service/pkg/database/postgresdb"
	"github.com/GAZIMAGomeDDD/billing-service/pkg/exchangeratesapi"
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
	cr := currency.New("", cache, 10)
	cr.GetCurrencyRateOfRuble = exchangeratesapi.GetCurrencyRateOfRuble
	h := handler.New(s, cr)
	srv := server.NewServer(h.Init())
	srv.Run()
}
