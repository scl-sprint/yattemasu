package persistence

import (
	"testing"
	"zemmai-dev/yattemasu/domain/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestUserPersistence(t *testing.T) {
	dsn := "user1:user1-passwd@tcp(127.0.0.1:3306)/test-db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		t.Fatal(err)
	}

	db.AutoMigrate(&model.User{}, &model.Group{})

	up := NewUserPersistence(db)

	lineId := "U4e88ac334487b33946de54c48b93ff14"

	user, err := up.Create(lineId)

	if err != nil {
		t.Fatal(err)
	}

	if user.LineID != lineId {
		t.Fatal("User create failed")
	}

	user, err = up.Find(lineId)

	if err != nil {
		t.Fatal(err)
	}

	if user.LineID != lineId {
		t.Fatal("User find failed")
	}
}
