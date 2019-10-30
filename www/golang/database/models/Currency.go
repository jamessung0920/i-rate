package models

import "time"

type (
	Currency struct {
		ID           uint `gorm:"primary_key"`
		CreatedAt    time.Time
		UpdatedAt    time.Time
		Currency      string  `gorm:"currency"`
	}
)

func (Currency) TableName() string {
	return "currency"
}
