package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/m1ker1n/go-developer-test/internal/config"
	"github.com/m1ker1n/go-developer-test/internal/storage/postgres"
	"github.com/shopspring/decimal"
)

func main() {
	cfg := config.MustLoad()

	ctx := context.Background()
	storage, err := postgres.NewStorage(ctx, cfg.Postgres.ConnectionString)
	if err != nil {
		panic(err)
	}
	defer storage.Close(ctx)

	wallet, err := storage.CreateWallet(ctx, decimal.NewFromFloat(100.33))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", wallet)

	id, _ := uuid.Parse("194a9fa0-98d3-4443-bf9a-e2f96c39618e")
	w2, _ := storage.GetWallet(ctx, id)
	fmt.Printf("%+v\n", w2)
	fmt.Printf("%s", w2.Balance)
}
