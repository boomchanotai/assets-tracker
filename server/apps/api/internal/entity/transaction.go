package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TxType string

const (
	TxTypeDeposit  TxType = "DEPOSIT"
	TxTypeWithdraw TxType = "WITHDRAW"
	TxTypeTransfer TxType = "TRANSFER"
)

func (tt TxType) String() string {
	return string(tt)
}

type Transaction struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	FromPocketID *uuid.UUID // Deposit == nil, Withdraw == PocketID, Transfer == FromPocketID
	ToPocketID   *uuid.UUID // Deposit == PocketID, Withdraw == nil, Transfer == ToPocketID
	Type         TxType
	Amount       decimal.Decimal
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (t Transaction) String() string {
	return t.ID.String() + " " + t.Type.String()
}

type TransactionInput struct {
	UserID       uuid.UUID
	FromPocketID *uuid.UUID
	ToPocketID   *uuid.UUID
	Type         TxType
	Amount       decimal.Decimal
}
