package models

import "time"

type (
	FollowList struct {
		ID            uint `gorm:"primary_key"`
		CreatedAt     time.Time
		UpdatedAt     time.Time
		UserID        string   `gorm:"user_id"`
		Currency      Currency `gorm: "foreignkey:CurrencyID"`
		CurrencyID    uint     `gorm:"currency_id"`
		WishBuyInRate string   `gorm:"wish_buy_in_rate"`
	}
)

func (FollowList) TableName() string {
	return "follow_list"
}
