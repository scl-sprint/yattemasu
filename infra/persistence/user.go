package persistence

import (
	"zemmai-dev/yattemasu/domain/model"
	"zemmai-dev/yattemasu/domain/repository"

	"gorm.io/gorm"
)

type UserPersistence struct {
	Conn *gorm.DB
}

func NewUserPersistence(conn *gorm.DB) repository.UserRepository {
	return &UserPersistence{Conn: conn}
}

func (up *UserPersistence) Find(LineID string) (*model.User, error) {
	var user model.User

	result := up.Conn.First(&user, "line_id = ?", LineID)

	return &user, result.Error
}

func (up *UserPersistence) Create(LineID string) (*model.User, error) {
	var user model.User

	result := up.Conn.First(&user, "line_id = ?", LineID)

	if result.Error != gorm.ErrRecordNotFound {
		return &user, nil
	}

	user = model.User{LineID: LineID}

	createResult := up.Conn.Create(&user)

	return &user, createResult.Error
}

func (up *UserPersistence) Delete(user *model.User) (error) {

	result := up.Conn.Delete(&user)

	return result.Error
}

func (up *UserPersistence) Update(LineID string, user *model.User) (*model.User, error) {
	result := up.Conn.Save(&user)

	return user, result.Error
}