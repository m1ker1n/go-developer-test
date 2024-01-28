package main

import (
	"context"
	"fmt"
	"github.com/m1ker1n/go-developer-test/internal/config"
	"github.com/m1ker1n/go-developer-test/internal/services"
	"github.com/m1ker1n/go-developer-test/internal/storage/postgres"
)

func main() {
	cfg := config.MustLoad()

	ctx := context.Background()
	storage, err := postgres.NewStorage(ctx, cfg.Postgres.ConnectionString)
	if err != nil {
		panic(err)
	}
	defer storage.Close(ctx)

	walletService := services.NewWallet(cfg.Wallet.InitialBalance, storage)
	w, _ := walletService.CreateWallet(ctx)
	fmt.Printf("%+v", w)
}
