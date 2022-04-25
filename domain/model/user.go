package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	LineID string `gorm:"unique;type:varchar(128)"`
	Location Location `gorm:"embedded"`
	Groups []Group `gorm:"foreignKey:UserLineID;references:LineID"`
}