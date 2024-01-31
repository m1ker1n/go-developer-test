package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/m1ker1n/go-developer-test/internal/models"
	"github.com/m1ker1n/go-developer-test/internal/storage"
	"github.com/shopspring/decimal"
)

type transactionStorage interface {
	CreateTransaction(ctx context.Context, to, from uuid.UUID, amount decimal.Decimal) (models.Transaction, error)
	GetTransactions(ctx context.Context, walletId uuid.UUID) ([]models.Transaction, error)
}

type Transaction struct {
	transactionStorage transactionStorage
}

var (
	ErrTransactionUndefined          = errors.New("some error happened")
	ErrTransactionWalletFromNotFound = errors.New("source wallet not found")
	ErrTransactionWalletToNotFound   = errors.New("destination wallet not found")
	ErrNotEnoughMoney                = errors.New("not enough money to send")
)

func NewTransaction(transactionStorage transactionStorage) *Transaction {
	return &Transaction{transactionStorage: transactionStorage}
}

func (t *Transaction) CreateTransaction(ctx context.Context, to, from uuid.UUID, amount decimal.Decimal) (models.Transaction, error) {
	tr, err := t.transactionStorage.CreateTransaction(ctx, to, from, amount)
	if err != nil {
		if errors.Is(err, storage.ErrWalletFromNotFound) {
			return models.Transaction{}, ErrTransactionWalletFromNotFound
		}
		if errors.Is(err, storage.ErrNotEnoughMoney) {
			return models.Transaction{}, ErrNotEnoughMoney
		}
		if errors.Is(err, storage.ErrWalletToNotFound) {
			return models.Transaction{}, ErrTransactionWalletToNotFound
		}

		//TODO: log
		return models.Transaction{}, ErrTransactionUndefined
	}
	return tr, nil
}

func (t *Transaction) GetTransactions(ctx context.Context, walletId uuid.UUID) ([]models.Transaction, error) {
	tr, err := t.transactionStorage.GetTransactions(ctx, walletId)
	if err != nil {
		if errors.Is(err, storage.ErrWalletNotFound) {
			return nil, ErrWalletNotFound
		}

		//TODO: log
		return nil, ErrTransactionUndefined
	}
	return tr, err
}
