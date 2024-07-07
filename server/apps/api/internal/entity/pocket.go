package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PocketType string

const (
	PocketTypeCashBox PocketType = "CASHBOX"
	PocketTypeNormal  PocketType = "NORMAL"
)

type Pocket struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	Name      string
	Type      PocketType
	Balance   decimal.Decimal
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p Pocket) String() string {
	return p.Name + " " + p.Balance.String()
}

type PocketInput struct {
	UserID    uuid.UUID
	AccountID uuid.UUID
	Name      string
	Type      PocketType
	Balance   decimal.Decimal
}
