package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/m1ker1n/go-developer-test/internal/config"
	"github.com/m1ker1n/go-developer-test/internal/storage/postgres"
	"github.com/shopspring/decimal"
)

func main() {
	cfg := config.MustLoad()
	//TODO: remove printf
	fmt.Printf("%+v\n", cfg)

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.Postgres.ConnectionString)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)
	pgxdecimal.Register(conn.TypeMap())

	queries := postgres.New(conn)

	balance, err := pgxdecimal.Decimal(decimal.NewFromFloat32(100.333)).NumericValue()
	if err != nil {
		panic(err)
	}

	createdWallet, err := queries.CreateWallet(ctx, balance)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", createdWallet)

	id, uuidErr := uuid.Parse("47e3eca5-5365-458d-8ef3-1f2e7573dfa7")
	if uuidErr != nil {
		panic(uuidErr)
	}
	gotWallet, err := queries.GetWallet(ctx, id)
	if err != nil {
		panic(err)
	}

	var decimalFromNumeric pgxdecimal.Decimal
	decimalFromNumericErr := decimalFromNumeric.ScanNumeric(gotWallet.Balance)
	if decimalFromNumericErr != nil {
		panic(decimalFromNumericErr)
	}
	fmt.Printf("%s", decimal.Decimal(decimalFromNumeric))
}
