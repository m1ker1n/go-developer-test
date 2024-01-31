package main

import (
	"context"

	"github.com/m1ker1n/go-developer-test/internal/config"
	"github.com/m1ker1n/go-developer-test/internal/http"
	"github.com/m1ker1n/go-developer-test/internal/http/handlers"
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
	transactionService := services.NewTransaction(storage)

	walletHandler := handlers.NewWalletHandler(walletService, transactionService)

	server := http.NewServer(cfg.HTTP.Addr(), walletHandler)

	if err := server.Run(); err != nil {
		panic(err)
	}
}
