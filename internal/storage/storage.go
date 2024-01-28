package storage

import "errors"

// Errors of wallets
var (
	ErrWalletNotFound = errors.New("wallet not found")
)

// Errors of transactions
var (
	ErrWalletFromNotFound = errors.New("source wallet not found")
	ErrWalletToNotFound   = errors.New("destination wallet not found")
	// Maybe it should be checked in service layer, but I don't know
	ErrNotEnoughMoney = errors.New("not enough money to send")
)
