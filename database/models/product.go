package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ID         uint64 `gorm:"primary_key"`
	Name       string
	CategoryID uint64
	Category   Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
