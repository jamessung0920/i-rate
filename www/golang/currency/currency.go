package currency

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/jinzhu/gorm"
	"net/http"
	
	"app/common"
	"app/database"
	"app/database/models"
	"app/crawler"
)

var Log = common.NewLogger()

func AddRoute(r *gin.Engine) {
	// Log.info("hahaha")
	currency := r.Group("currency")
	{
		currency.GET("taiwan-bank", getTaiwanBankCurrencyRate)
		currency.GET("crawl", doAutoCrawlRate)
	}
}

func getTaiwanBankCurrencyRate(c *gin.Context) {
	// queryString := c.Request.URL.Query()
	// fmt.Println(queryString["currency"])
	_, currencyList := getCurrencyList()
	for key, value := range currencyList {
		go crawler.Crawler(key, value)
	}
	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}

func getCurrencyList() (error, map[uint]string) {
	var currency []models.Currency
	if queryErr := database.GetDB().Find(&currency).Error; queryErr != nil {
		Log.Warn("query currency error")
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

func GetCurrencyLatestRate(currency string) (error, models.Rate) {
	var currencyModel models.Currency
	var rateModel models.Rate
	if queryErr := database.GetDB().Where("currency = ?", currency).First(&currencyModel).Error; queryErr != nil {
		Log.Warn("query currency error")
		return queryErr, rateModel
	}
	if queryErr := database.GetDB().Where("currency_id = ?", currencyModel.ID).Last(&rateModel).Error; queryErr != nil {
		Log.Warn("query rate error")
		return queryErr, rateModel
	} else {
		// var currencyList []string
		fmt.Println(rateModel)
		// fmt.Println(currency.Currency)
		return nil, rateModel
	}
}

func doAutoCrawlRate(c *gin.Context) {
	portGo := os.Getenv("PORT_GO")
	defer func(){
		fmt.Println("close crawling ...")
		Log.Info("close crawling ...")
	}()
	go func(){
		for {
			common.CheckGoRoutineNum()
			time.Sleep(10 * time.Second)
		}
	}()
	go func(){
		for {
			fmt.Println("do crawl rate data ...")
			Log.Info("do crawl rate data ...")
			common.CallAPI("http://172.17.0.1:" + portGo + "/currency/taiwan-bank", "GET", nil)
			time.Sleep(5 * time.Minute)
		}
	}()
	c.JSON(http.StatusOK, gin.H{
		"message": "start crawling rate",
		"result": true,
	})
}