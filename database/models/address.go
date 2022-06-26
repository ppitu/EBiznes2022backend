package models

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	ID   uint64 `gorm:"primary_key"`
	City string
}
