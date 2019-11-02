package currency

import (
	"fmt"
	"os"
	"time"
	"strings"
	// "database/sql"

	"github.com/gin-gonic/gin"
	// "github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"net/http"
	
	"app/common"
	"app/database"
	"app/database/models"
	"app/crawler"
)

var Log = common.NewLogger()

type CurrencyConfig struct {
	Statement	map[string][]string	`json:"statement"`
}

func AddRoute(r *gin.Engine) {
	// Log.info("hahaha")
	currency := r.Group("currency")
	{
		currency.GET("taiwan-bank", getTaiwanBankCurrencyRate)
		currency.GET("ctbc-bank", getCtbcBankCurrencyRate)
		currency.GET("esun-bank", getEsunBankCurrencyRate)
		currency.GET("crawl", doAutoCrawlRate)
	}
}

func getTaiwanBankCurrencyRate(c *gin.Context) {
	go crawler.Crawler("taiwan-bank")
	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}

func getCtbcBankCurrencyRate(c *gin.Context) {
	go crawler.Crawler("ctbc-bank")
	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}

func getEsunBankCurrencyRate(c *gin.Context) {
	go crawler.Crawler("esun-bank")
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

func GetCurrencyStatement(currencyList []string) (error, map[string][]string) {
	_, statementConfig := getStatementConfig()

	currencyStatement := make(map[string][]string)
	for _, currency := range currencyList {
		key := strings.ToLower(currency)
		currencyStatement[key] = statementConfig.Statement[key]
	}
	return nil, currencyStatement
}

func GetCurrencyLatestRate(currency string) (error, []models.Rate) {
	var currencyModel models.Currency
	var rateModel []models.Rate

	if queryErr := database.GetDB().Where("currency = ?", currency).First(&currencyModel).Error; queryErr != nil {
		Log.Warn("query currency error")
		return queryErr, rateModel
	}

	queryLatestErr := database.GetDB().Raw(`SELECT T1.* FROM rate T1 inner JOIN
												(
													SELECT crawl_from,
														MAX(created_at) max_created_time
													FROM rate
													GROUP BY crawl_from
												) T2
											ON T1.created_at = T2.max_created_time
											WHERE currency_id = ?`, currencyModel.ID).Find(&rateModel).Error
	if queryLatestErr != nil {
		Log.Warn("query latest rate error")
		return queryLatestErr, rateModel
	} else {
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
			common.CallAPI("http://172.17.0.1:" + portGo + "/currency/ctbc-bank", "GET", nil)
			common.CallAPI("http://172.17.0.1:" + portGo + "/currency/esun-bank", "GET", nil)
			time.Sleep(5 * time.Minute)
		}
	}()
	c.JSON(http.StatusOK, gin.H{
		"message": "start crawling rate",
		"result": true,
	})
}

func getStatementConfig() (error, CurrencyConfig) {
	var config CurrencyConfig
    viper.SetConfigName("currency")   // 设置配置文件名 (不带后缀)
    viper.AddConfigPath("./configs")        // 第一个搜索路径
    err := viper.ReadInConfig()     // 读取配置数据
    if err != nil {
        return err, config
    }
    viper.Unmarshal(&config)        // 将配置信息绑定到结构体上

	return nil, config
}