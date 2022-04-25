package usecase

import (
	"zemmai-dev/yattemasu/domain/model"
	"zemmai-dev/yattemasu/domain/repository"
)

type UserUsecase interface {
	Find(LineID string) (*model.User, error)
	Create(LineID string) (*model.User, error)
	Delete(user *model.User) error
	SetLocation(user *model.User, loc model.Location) (*model.User, error)
	GetLocation(LineID string) (model.Location, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func (uu *userUsecase) Create(LineID string) (*model.User, error) {
	user, err := uu.userRepo.Create(LineID)

	return user, err
}

func (uu *userUsecase) Delete(user *model.User) error {
	err := uu.userRepo.Delete(user)

	return err
}

func (uu *userUsecase) Find(LineID string) (*model.User, error) {
	user, err := uu.userRepo.Find(LineID)

	return user, err
}
func (uu *userUsecase) SetLocation(user *model.User, loc model.Location) (*model.User, error) {

	user.Location = loc

	user, err := uu.userRepo.Update(user.LineID, user)

	return user, err
}
func (uu *userUsecase) GetLocation(LineID string) (model.Location, error) {

	user, err := uu.Find(LineID)

	return user.Location, err
}
