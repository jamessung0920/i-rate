package currency

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	// "database/sql"

	"github.com/gin-gonic/gin"
	// "github.com/jinzhu/gorm"
	"net/http"

	"github.com/spf13/viper"

	"app/common"
	"app/crawler"
	"app/database"
	"app/database/models"
)

var Log = common.NewLogger()

type CurrencyConfig struct {
	Statement map[string][]string `json:"statement"`
}

func AddRoute(r *gin.Engine) {
	// Log.info("hahaha")
	currency := r.Group("currency")
	{
		currency.GET("taiwan-bank", getTaiwanBankCurrencyRate)
		currency.GET("ctbc-bank", getCtbcBankCurrencyRate)
		currency.GET("esun-bank", getEsunBankCurrencyRate)
		currency.GET("crawl", doAutoCrawlRate)
		currency.GET("rate/notify", checkFollowRate)
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

func GetCurrencyStatement(currencyList []string) (error, map[string][]string) {
	_, statementConfig := getStatementConfig()

	currencyStatement := make(map[string][]string)
	for _, currency := range currencyList {
		key := strings.ToLower(currency)
		currencyStatement[key] = statementConfig.Statement[key]
	}
	return nil, currencyStatement
}

func GetCurrencyLatestRate(currency interface{}) (error, []models.Rate) {
	var currencyModel models.Currency
	var rateModel []models.Rate
	var currencyID uint

	v := reflect.ValueOf(currency).Kind()
	if v == reflect.Uint {
		currencyID = currency.(uint)
	} else if v == reflect.String {
		if queryErr := database.GetDB().Where("currency = ?", currency).First(&currencyModel).Error; queryErr != nil {
			Log.Error("query currency error")
			return queryErr, rateModel
		}
		currencyID = currencyModel.ID
	}
	queryLatestErr := database.GetDB().Raw(`SELECT T1.* FROM rate T1 inner JOIN
												(
													SELECT crawl_from,
														MAX(created_at) max_created_time
													FROM rate
													GROUP BY crawl_from
												) T2
											ON T1.created_at = T2.max_created_time AND
											   T1.crawl_from = T2.crawl_from
											WHERE currency_id = ?`, currencyID).Find(&rateModel).Error
	if queryLatestErr != nil {
		Log.Error("query latest rate error")
		return queryLatestErr, rateModel
	} else {
		return nil, rateModel
	}
}

func doAutoCrawlRate(c *gin.Context) {
	portGo := os.Getenv("PORT_GO")
	go func() {
		for {
			common.CheckGoRoutineNum()
			time.Sleep(10 * time.Second)
		}
	}()
	go func() {
		defer func() {
			Log.Info("close crawling ...")
		}()
		for {
			Log.Info("do crawl rate data ...")
			common.CallAPI("http://172.17.0.1:"+portGo+"/currency/taiwan-bank", "GET", nil, nil)
			common.CallAPI("http://172.17.0.1:"+portGo+"/currency/ctbc-bank", "GET", nil, nil)
			common.CallAPI("http://172.17.0.1:"+portGo+"/currency/esun-bank", "GET", nil, nil)
			time.Sleep(5 * time.Minute)
		}
	}()
	c.JSON(http.StatusOK, gin.H{
		"message": "start crawling rate",
		"result":  true,
	})
}

func checkFollowRate(c *gin.Context) {
	go func() {
		defer func() {
			Log.Info("close checking ...")
		}()
		for {
			Log.Info("do checking follow rate ...")
			_, followList := common.GetAllFollowList()
			for _, follow := range followList {
				_, latestRateList := GetCurrencyLatestRate(follow.CurrencyID)
				for _, rate := range latestRateList {
					wishRate, _ := strconv.ParseFloat(follow.WishBuyInRate, 32)
					SellRate, _ := strconv.ParseFloat(rate.SellRate, 32)
					if wishRate >= SellRate {
						var rateBuffer bytes.Buffer
						rateTimeString := rate.RateTime.Format("2006-01-02 15:04:05")
						rateBuffer.WriteString(fmt.Sprintf("您關注的貨幣達到理想匯率囉！\n\n%s\n\n貨幣: %s\n匯率時間: %s\n理想匯率: %s\n本行賣出匯率: %s\n\n", rate.CrawlFrom, follow.Currency.Currency, rateTimeString, follow.WishBuyInRate, rate.SellRate))
						common.PushMessage(follow.UserID, rateBuffer.String())
					}
				}
			}
			time.Sleep(5 * time.Minute)
		}
	}()
	c.JSON(http.StatusOK, gin.H{
		"message": "start checking following rate",
		"result":  true,
	})
}

func getStatementConfig() (error, CurrencyConfig) {
	var config CurrencyConfig
	viper.SetConfigName("currency")  // 设置配置文件名 (不带后缀)
	viper.AddConfigPath("./configs") // 第一个搜索路径
	err := viper.ReadInConfig()      // 读取配置数据
	if err != nil {
		return err, config
	}
	viper.Unmarshal(&config) // 将配置信息绑定到结构体上

	return nil, config
}
