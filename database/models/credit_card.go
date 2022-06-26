package models

import "gorm.io/gorm"

type CreditCard struct {
	gorm.Model
	ID     uint64 `gorm:"primary_key"`
	Number string
}
