package repository

import (
	"zemmai-dev/yattemasu/domain/model"
)

type UserRepository interface {
	Find(LineID string) (*model.User, error)
	Create(LineID string) (*model.User, error)
	Delete(user *model.User) (error)
	Update(LineID string, user *model.User) (*model.User, error)
}