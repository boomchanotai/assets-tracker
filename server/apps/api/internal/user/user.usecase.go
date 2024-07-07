package user

type usecase struct {
	userRepo Repository
}

func NewUsecase(userRepo Repository) *usecase {
	return &usecase{
		userRepo: userRepo,
	}
}
