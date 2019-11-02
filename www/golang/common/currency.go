package common

import (
	"app/database"
	"app/database/models"
)

func GetCurrencyList() (error, []string) {
	var currency []models.Currency
	if queryErr := database.GetDB().Find(&currency).Error; queryErr != nil {
		Log.Warn("query currency error")
		return queryErr, nil
	} else {
		// var currencyList []string
		var currencyList []string 
		for _, value := range currency {
			currencyList = append(currencyList, value.Currency)
		}
		// fmt.Println(currency.Currency)
		return nil, currencyList
	}
}