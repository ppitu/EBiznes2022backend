package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	ID        uint64 `gorm:"primary_key"`
	ProductID uint64
	Product   Product
}
