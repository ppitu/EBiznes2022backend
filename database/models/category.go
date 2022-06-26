package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	ID   uint64 `gorm:"primary_key"`
	Name string
}
