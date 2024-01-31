package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/m1ker1n/go-developer-test/internal/models"
	"github.com/m1ker1n/go-developer-test/internal/storage"
	"github.com/shopspring/decimal"
)

type walletStorage interface {
	CreateWallet(ctx context.Context, balance decimal.Decimal) (models.Wallet, error)
	GetWallet(ctx context.Context, id uuid.UUID) (models.Wallet, error)
}

type Wallet struct {
	initialWalletBalance decimal.Decimal
	walletStorage        walletStorage
}

var (
	ErrWalletNotFound  = errors.New("could not found wallet")
	ErrWalletUndefined = errors.New("some error happened")
)

func NewWallet(initialWalletBalance decimal.Decimal, walletStorage walletStorage) *Wallet {
	return &Wallet{initialWalletBalance: initialWalletBalance, walletStorage: walletStorage}
}

func (w *Wallet) CreateWallet(ctx context.Context) (models.Wallet, error) {
	wallet, err := w.walletStorage.CreateWallet(ctx, w.initialWalletBalance)
	if err != nil {
		//TODO: log internal error of wallet creation
		return models.Wallet{}, ErrWalletUndefined
	}
	return wallet, nil
}

func (w *Wallet) GetWallet(ctx context.Context, walletId uuid.UUID) (models.Wallet, error) {
	wallet, err := w.walletStorage.GetWallet(ctx, walletId)
	if err != nil {
		if errors.Is(err, storage.ErrWalletNotFound) {
			return models.Wallet{}, ErrWalletNotFound
		}
		//TODO: log undefined error getting wallet
		return models.Wallet{}, ErrWalletUndefined
	}
	return wallet, nil
}
