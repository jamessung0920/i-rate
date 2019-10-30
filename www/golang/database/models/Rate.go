package models

import (
	"time"
)

type (
	Rate struct {
		ID           uint `gorm:"primary_key"`
		CreatedAt    time.Time
		UpdatedAt    time.Time
		RateTime	 time.Time
		BuyRate      string  `gorm:"buy_rate"`
		SellRate     string  `gorm:"sell_rate"`
		CrawlFrom	 string  `gorm: "crawl_from"`
		Currency	 Currency `gorm: "foreignkey:CurrencyID"`
		CurrencyID	 uint  `gorm:"currency_id"`
	}
)

func (Rate) TableName() string {
	return "rate"
}
