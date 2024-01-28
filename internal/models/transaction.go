package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	ID         uuid.UUID
	Time       time.Time
	WalletFrom uuid.UUID
	WalletTo   uuid.UUID
	Amount     decimal.Decimal
}
