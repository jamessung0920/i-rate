package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"

	"app/common"
	"app/currency"
)

var Log = common.NewLogger()

type IncomingMessageInfo struct {
	Events []struct {
		Type       string `json:"type"`
		ReplyToken string `json:"replyToken"`
		Source     struct {
			UserID string `json:"userId"`
			Type   string `json:"type"`
		} `json:"source"`
		Timestamp int64 `json:"timestamp"`
		Message   struct {
			Type string `json:"type"`
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"events"`
	Destination string `json:"destination"`
}

type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func AddRoute(r *gin.Engine) {
	hook := r.Group("webhook")
	{
		hook.POST("/", webhook)
	}
}

func webhook(c *gin.Context) {
	fmt.Println("webhook!")

	var incomingData IncomingMessageInfo
	body, _ := ioutil.ReadAll(c.Request.Body)

	json.Unmarshal(body, &incomingData)
	replyToken := incomingData.Events[0].ReplyToken
	incomingUserID := incomingData.Events[0].Source.UserID
	incomingMessage := incomingData.Events[0].Message.Text
	incomingMessage = strings.TrimSpace(incomingMessage)

	_, currencyList := common.GetCurrencyList()
	_, currencyStatement := currency.GetCurrencyStatement(currencyList)
	_, followList := common.GetUnsureFollowList(incomingUserID)

	var responseMessage string
	var isCurrency bool
	var isRateNumber bool
	var oriCurrencyID uint = followList.CurrencyID
	if incomingMessage == "貨幣列表" {
		var currencyBuffer bytes.Buffer
		for key, currency := range currencyList {
			currencyBuffer.WriteString(currencyStatement[strings.ToLower(currency)][0])
			currencyBuffer.WriteString(" - ")
			currencyBuffer.WriteString(currency)
			if key != len(currencyList)-1 {
				currencyBuffer.WriteString("\n")
			}
		}
		responseMessage = currencyBuffer.String()
	} else if incomingMessage == "關注" {
		common.CreateUserFollowData(incomingUserID)
		responseMessage = "請輸入幣別"
	} else if incomingMessage == "使用說明" {
		responseMessage = fmt.Sprintf("點選 \"貨幣列表\" 可查看目前支援貨幣匯率\n----------------------\n點選 \"關注\" 後依照指示設定理想匯率主動通知")
	} else if currencyKey, findOk := common.Mapkey(currencyStatement, strings.ToUpper(incomingMessage)); findOk {
		isCurrency = true
		if followList.UserID != "" && oriCurrencyID == 0 {
			var currencyID uint
			_, currencyListWithKey := common.GetCurrencyListWithKey()
			for k, v := range currencyListWithKey {
				if v == strings.ToUpper(currencyKey) {
					currencyID = k
					break
				}
			}
			followList.CurrencyID = currencyID
			common.UpdateUserFollowData(followList)
			responseMessage = "請輸入期望的匯率"
		} else {
			fmt.Println("rate response!")
			_, bankLatestRateData := currency.GetCurrencyLatestRate(strings.ToUpper(currencyKey))
			var rateBuffer bytes.Buffer
			for _, eachBankRate := range bankLatestRateData {
				rateTimeString := eachBankRate.RateTime.Format("2006-01-02 15:04:05")
				rateBuffer.WriteString(fmt.Sprintf("%s\n匯率時間: %s\n本行買入匯率: %s\n本行賣出匯率: %s\n===============\n", eachBankRate.CrawlFrom, rateTimeString, eachBankRate.BuyRate, eachBankRate.SellRate))
			}
			responseMessage = rateBuffer.String()
		}
	} else if common.IsNumeric(incomingMessage) {
		isRateNumber = true
		if oriCurrencyID != 0 {
			if _, hasExist, hasFollowList := common.HasCurrencyFollowed(incomingUserID, oriCurrencyID); hasExist {
				common.CleanUserFollowData(hasFollowList)
			}
			followList.WishBuyInRate = incomingMessage
			common.UpdateUserFollowData(followList)
			responseMessage = "關注成功"
		}
	} else {
		// responseMessage = incomingMessage
	}

	if !isCurrency && followList.UserID != "" && oriCurrencyID == 0 {
		common.CleanUserFollowData(followList)
	} else if !isRateNumber && oriCurrencyID != 0 {
		common.CleanUserFollowData(followList)
	}
	common.ReplyMessage(replyToken, responseMessage)
}
