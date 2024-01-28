package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID      uuid.UUID
	Balance decimal.Decimal
}
