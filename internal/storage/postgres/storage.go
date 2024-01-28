package postgres

import (
	"context"
	"errors"
	"github.com/google/uuid"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/m1ker1n/go-developer-test/internal/models"
	"github.com/m1ker1n/go-developer-test/internal/storage"
	"github.com/shopspring/decimal"
)

type Storage struct {
	queries *Queries
	conn    *pgx.Conn
}

func NewStorage(ctx context.Context, connectionString string) (*Storage, error) {
	s := &Storage{}

	conn, err := pgx.Connect(ctx, connectionString)
	if err != nil {
		return nil, err
	}
	s.conn = conn
	pgxdecimal.Register(conn.TypeMap())

	queries := New(conn)
	s.queries = queries
	return s, nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}

func (s *Storage) CreateWallet(ctx context.Context, balance decimal.Decimal) (models.Wallet, error) {
	wallet, err := s.queries.CreateWallet(ctx, balance)
	if err != nil {
		//TODO: log internal error for creating wallet
		return models.Wallet{}, err
	}

	return models.Wallet(wallet), nil
}

func (s *Storage) GetWallet(ctx context.Context, id uuid.UUID) (models.Wallet, error) {
	wallet, err := s.queries.GetWallet(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Wallet{}, storage.ErrWalletNotFound
		}
		//TODO: log internal error for not getting wallet
		return models.Wallet{}, err
	}

	return models.Wallet(wallet), nil
}

func (s *Storage) Send(ctx context.Context, from, to uuid.UUID, amount decimal.Decimal) error {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		//TODO: log internal error for transaction
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			//TODO: log error of rollback
		}
	}(tx, ctx)

	qtx := s.queries.WithTx(tx)
	{
		fromWallet, err := qtx.GetWallet(ctx, from)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return storage.ErrWalletFromNotFound
			}
			//TODO: log internal error for something else happened
			return err
		}

		if enoughMoney := fromWallet.Balance.GreaterThanOrEqual(amount); !enoughMoney {
			return storage.ErrNotEnoughMoney
		}

		toWallet, err := qtx.GetWallet(ctx, to)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return storage.ErrWalletToNotFound
			}
			//TODO: log internal error for something else happened
			return err
		}

		newFromBalance := fromWallet.Balance.Sub(amount)
		err = qtx.SetBalance(ctx, SetBalanceParams{
			ID:      fromWallet.ID,
			Balance: newFromBalance,
		})
		if err != nil {
			//TODO: log something happened while setting balance
			return err
		}

		newToBalance := toWallet.Balance.Add(amount)
		err = qtx.SetBalance(ctx, SetBalanceParams{
			ID:      toWallet.ID,
			Balance: newToBalance,
		})
		if err != nil {
			//TODO: log something happened while setting balance
			return err
		}

		_, err = qtx.CreateTransaction(ctx, CreateTransactionParams{
			WalletFromID: fromWallet.ID,
			WalletToID:   toWallet.ID,
			Amount:       amount,
		})
		if err != nil {
			//TODO: log something happened while creating transaction
			return err
		}
	}
	return tx.Commit(ctx)
}

func (s *Storage) GetTransactions(ctx context.Context, walletId uuid.UUID) ([]models.Transaction, error) {
	transactions, err := s.queries.ListTransactionsByWalletId(ctx, walletId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrWalletNotFound
		}
		//TODO: log something wrong happened
		return nil, err
	}

	result := make([]models.Transaction, len(transactions))
	for i := range result {
		result[i] = transactions[i].toModel()
	}
	return result, err
}

func (t Transaction) toModel() models.Transaction {
	return models.Transaction{
		ID:         t.ID,
		Time:       t.Time.Time,
		WalletFrom: t.WalletFromID,
		WalletTo:   t.WalletToID,
		Amount:     t.Amount,
	}
}
