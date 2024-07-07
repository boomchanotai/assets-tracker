package entity

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrInvalidAccountType = errors.New("INVALID_ACCOUNT_TYPE")
)

type AccountType string

const (
	AccountTypeSaving       AccountType = "SAVING"
	AccountTypeFixedDeposit AccountType = "FIXED_DEPOSIT"
	AccountTypeFCD          AccountType = "FCD"
	AccountTypeMutualFund   AccountType = "MUTUAL_FUND"
	AccountTypeStock        AccountType = "STOCK"
)

func (at AccountType) String() string {
	return string(at)
}

type Account struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Type      AccountType
	Name      string
	Bank      string
	Balance   decimal.Decimal
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a Account) String() string {
	return a.Name
}

type AccountInput struct {
	UserID  uuid.UUID
	Type    AccountType
	Name    string
	Bank    string
	Balance decimal.Decimal
}
