package common

import (
	"app/database"
	"app/database/models"
)

func GetUnsureFollowList(userID string) (error, models.FollowList) {
	var followList models.FollowList
	if queryErr := database.GetDB().Where("user_id = ? AND (currency_id = ? OR wish_buy_in_rate = ?)", userID, 0, "").Find(&followList).Error; queryErr != nil {
		Log.Error("query follow list error")
		return queryErr, followList
	} else {
		return nil, followList
	}
}

func GetAllFollowList() (error, []models.FollowList) {
	var followList []models.FollowList
	if queryErr := database.GetDB().Where("currency_id != ? AND wish_buy_in_rate != ?", 0, "").Preload("Currency").Find(&followList).Error; queryErr != nil {
		Log.Error("query follow list error")
		return queryErr, followList
	} else {
		return nil, followList
	}
}

func HasCurrencyFollowed(userID string, currencyID uint) (error, bool, models.FollowList) {
	var followList models.FollowList
	var count int
	if queryErr := database.GetDB().Where("user_id = ? AND currency_id = ? AND wish_buy_in_rate != ?", userID, currencyID, "").Find(&followList).Count(&count).Error; queryErr != nil {
		Log.Error("query follow list error")
		return queryErr, count != 0, followList
	} else {
		return nil, count != 0, followList
	}
}

func CleanUserFollowData(followList models.FollowList) {
	if database.GetDB().Delete(&followList).Error != nil {
		Log.Error("delete follow list data error")
	}
}

func CreateUserFollowData(userID string) {
	if database.GetDB().Create(&models.FollowList{UserID: userID}).Error != nil {
		Log.Error("create follow list data fail!")
	}
}

func UpdateUserFollowData(followList models.FollowList) {
	if database.GetDB().Save(&followList).Error != nil {
		Log.Error("save follow list currency id fail!")
	}
}
