package account

import (
	"context"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/interfaces"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type usecase struct {
	accountRepo     interfaces.AccountRepository
	pocketRepo      interfaces.PocketRepository
	transactionRepo interfaces.TransactionRepository
}

func NewUsecase(accountRepo interfaces.AccountRepository, pocketRepo interfaces.PocketRepository, transactionRepo interfaces.TransactionRepository) *usecase {
	return &usecase{
		accountRepo:     accountRepo,
		pocketRepo:      pocketRepo,
		transactionRepo: transactionRepo,
	}
}

func (u *usecase) GetAccounts(ctx context.Context, userID uuid.UUID) ([]entity.Account, error) {
	accounts, err := u.accountRepo.GetUserAccounts(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get accounts")
	}

	return accounts, nil
}

func (u *usecase) GetAccount(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*entity.Account, error) {
	account, err := u.accountRepo.GetUserAccount(ctx, userID, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}

	return account, nil
}

func (u *usecase) CreateAccount(ctx context.Context, input entity.AccountInput) (*entity.Account, error) {
	account, err := u.accountRepo.CreateAccount(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create account")
	}

	_, err = u.pocketRepo.CreatePocket(ctx, entity.PocketInput{
		UserID:    input.UserID,
		AccountID: account.ID,
		Name:      "Cashbox",
		Type:      entity.PocketTypeCashBox,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cashbox pocket")
	}

	return account, nil
}

func (u *usecase) UpdateAccount(ctx context.Context, userID uuid.UUID, id uuid.UUID, input entity.AccountInput) (*entity.Account, error) {
	// Check ownership
	if _, err := u.accountRepo.GetUserAccount(ctx, userID, id); err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}

	account, err := u.accountRepo.UpdateAccount(ctx, id, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update account")
	}

	return account, nil
}

func (u *usecase) DeleteAccount(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	// Check ownership
	if _, err := u.accountRepo.GetUserAccount(ctx, userID, id); err != nil {
		return errors.Wrap(err, "failed to get account")
	}

	err := u.accountRepo.DeleteAccount(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete account")
	}

	return nil
}

func (u *usecase) getCashboxPocket(ctx context.Context, accountID uuid.UUID) (*entity.Pocket, error) {
	pockets, err := u.pocketRepo.GetPocketsByAccountID(ctx, accountID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pockets")
	}

	if len(pockets) == 0 {
		return nil, errors.New("pockets not found")
	}

	// Get cashbox pocket
	var cashbox entity.Pocket
	for _, pocket := range pockets {
		if pocket.Type == entity.PocketTypeCashBox {
			cashbox = pocket
			break
		}
	}

	return &cashbox, nil
}

func (u *usecase) Deposit(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, amount decimal.Decimal) (*entity.Account, error) {
	// TODO: Lock db transaction

	// Check ownership
	if _, err := u.accountRepo.GetUserAccount(ctx, userID, accountID); err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}

	// Update account balance
	account, err := u.accountRepo.Deposit(ctx, accountID, amount)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deposit")
	}

	cashbox, err := u.getCashboxPocket(ctx, accountID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cashbox pocket")
	}

	// Deposit to cashbox pocket
	if err = u.pocketRepo.Deposit(ctx, cashbox.ID, amount); err != nil {
		return nil, errors.Wrap(err, "failed to deposit to cashbox pocket")
	}

	// Create Transaction
	if _, err = u.transactionRepo.CreateTransaction(ctx, entity.TransactionInput{
		AccountID:    accountID,
		FromPocketID: nil,
		ToPocketID:   &cashbox.ID,
		Type:         entity.TxTypeDeposit,
		Amount:       amount,
	}); err != nil {
		return nil, errors.Wrap(err, "failed to create transaction")
	}

	return account, nil
}

func (u *usecase) UpdateBalance(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, amount decimal.Decimal) (*entity.Account, error) {
	// TODO: Lock db transaction

	// Check ownership
	if _, err := u.accountRepo.GetUserAccount(ctx, userID, accountID); err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}

	// Update account balance
	account, differenceBalance, err := u.accountRepo.UpdateBalance(ctx, accountID, amount)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update account")
	}

	cashbox, err := u.getCashboxPocket(ctx, accountID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cashbox pocket")
	}

	// Deposit to cashbox pocket
	if err = u.pocketRepo.Deposit(ctx, cashbox.ID, differenceBalance); err != nil {
		return nil, errors.Wrap(err, "failed to deposit to cashbox pocket")
	}

	//  Create Transaction
	if _, err = u.transactionRepo.CreateTransaction(ctx, entity.TransactionInput{
		AccountID:    accountID,
		FromPocketID: nil,
		ToPocketID:   &cashbox.ID,
		Type:         entity.TxTypeDeposit,
		Amount:       differenceBalance,
	}); err != nil {
		return nil, errors.Wrap(err, "failed to create transaction")
	}

	return account, nil
}
