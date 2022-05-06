package usecase

import (
	"testing"
	"zemmai-dev/yattemasu/domain/model"
	"zemmai-dev/yattemasu/infra/persistence"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestUserUsecase(t *testing.T) {
	dsn := "user1:user1-passwd@tcp(127.0.0.1:3306)/test-db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		t.Fatal(err)
	}

	db.AutoMigrate(&model.User{}, &model.Group{})

	up := persistence.NewUserPersistence(db)

	uu := NewUserUsecase(up)

	lineId := "U4e88ac334487b33946de54c48b93ff14"

	user, err := uu.Create(lineId)

	if err != nil {
		t.Fatal(err)
	}

	loc := model.Location{Longitude: 43, Latitude: 135}

	newUser, err := uu.SetLocation(user, loc)

	loc, locErr := uu.GetLocation(newUser.LineID)

	if err != nil {
		t.Fatal(err)
	}

	if locErr != nil {
		t.Fatal(locErr)
	}

	if newUser.Location != loc {
		t.Fatal("Location Set is invalid")
	}
}
