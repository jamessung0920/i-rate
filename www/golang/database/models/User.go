package models

import "time"

type (
	User struct {
		ID           uint `gorm:"primary_key"`
		CreatedAt    time.Time
		UpdatedAt    time.Time
		Email        string  `gorm:"email"`
		Password     string  `gorm:"password"`
		Name         string  `gorm:"name"`
		Username     string  `gorm:"username"`
		Points       int64   `gorm:"points"`
		Profit       float64 `gorm:"profit"`
		Active       string  `gorm:"active"`
		Status       string  `gorm:"status"`
		PendingPoint int64   `gorm:"pending_point"`
		TmpPoint     int64   `gorm:"tmp_point"`
		TmpProfit    float64 `gorm:"tmp_profit"`
	}
)

func (User) TableName() string {
	return "user"
}
