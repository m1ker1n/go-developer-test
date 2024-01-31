package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/m1ker1n/go-developer-test/internal/config"
	"github.com/m1ker1n/go-developer-test/internal/server"
	"github.com/m1ker1n/go-developer-test/internal/server/handlers"
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

	srv := server.New(cfg.HTTP.Addr(), walletHandler)

	go func() {
		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	quitSig := <-quit
	log.Printf("received shutdown signal: %s", quitSig)

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %s\n", err)
	}
	log.Println("server exiting")
}
