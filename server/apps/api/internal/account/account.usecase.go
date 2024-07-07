package account

type usecase struct {
	accountRepo Repository
}

func NewUsecase(accountRepo Repository) *usecase {
	return &usecase{
		accountRepo: accountRepo,
	}
}
