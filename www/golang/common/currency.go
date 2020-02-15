package common

import (
	"app/database"
	"app/database/models"
)

func GetCurrencyList() (error, []string) {
	var currency []models.Currency
	if queryErr := database.GetDB().Find(&currency).Error; queryErr != nil {
		Log.Error("query currency error")
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

func GetCurrencyListWithKey() (error, map[uint]string) {
	var currency []models.Currency
	if queryErr := database.GetDB().Find(&currency).Error; queryErr != nil {
		Log.Error("query currency error")
		return queryErr, nil
	} else {
		// var currencyList []string
		currencyList := make(map[uint]string)
		for _, value := range currency {
			// currencyList = append(currencyList, value.Currency)
			currencyList[value.ID] = value.Currency
		}
		// fmt.Println(currency.Currency)
		return nil, currencyList
	}
}
