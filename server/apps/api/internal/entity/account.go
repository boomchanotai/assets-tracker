package entity

import (
	"database/sql/driver"
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

func (at AccountType) Value() (driver.Value, error) {
	switch at {
	case AccountTypeSaving, AccountTypeFixedDeposit, AccountTypeFCD, AccountTypeMutualFund, AccountTypeStock:
		return string(at), nil
	}
	return nil, errors.Wrap(ErrInvalidAccountType, "invalid account type")
}

func (at *AccountType) Scan(value interface{}) error {
	var accountType AccountType
	if value == nil {
		*at = ""
		return nil
	}

	st, ok := value.(string)
	if !ok {
		return errors.Wrap(ErrInvalidAccountType, "invalid account type")
	}

	accountType = AccountType(st)
	switch accountType {
	case AccountTypeSaving, AccountTypeFixedDeposit, AccountTypeFCD, AccountTypeMutualFund, AccountTypeStock:
		*at = accountType
		return nil
	}

	return errors.Wrap(ErrInvalidAccountType, "invalid account type")
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
	Type AccountType
	Name string
	Bank string
}
