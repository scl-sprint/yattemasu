package model

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	LineID string
	UserLineID string
}