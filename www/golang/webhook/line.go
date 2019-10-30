package webhook

import (
	"fmt"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"

	"app/currency"
	"app/common"
)


var (
	Log = common.NewLogger()
	channelToken = "E6lkGp635WG1GKcKCViXIn5xCPPudDEvefYpFRRmhiOHgLzvUcQpoOo+3ZhGfKhymd1Mujyfj0ddLHpNPLjJm2GPw3pynN7KPuJQ9aUrvcSrQue7ibw8el1eO/Xnm+qUHCnFdrQcPF6Z6n000j2sMQdB04t89/1O/w1cDnyilFU="
	replyUrl = "https://api.line.me/v2/bot/message/reply"
)

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

type ResponseData struct {
	ReplyToken string `json:"replyToken"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func AddRoute(r *gin.Engine) {
	hook := r.Group("webhook")
	{
		hook.POST("/", webhook)
		// webhook.GET("some-site", get)
	}
}

func webhook(c *gin.Context) {
	fmt.Println("webhook!")

	var incomingData IncomingMessageInfo
	body,_ := ioutil.ReadAll(c.Request.Body)

	json.Unmarshal(body, &incomingData)
	replyToken := incomingData.Events[0].ReplyToken
	incomingMessage := incomingData.Events[0].Message.Text
	incomingMessage = strings.TrimSpace(incomingMessage)

	_, currencyList := currency.GetCurrencyList()
	_, currencyStatement := currency.GetCurrencyStatement(currencyList)

	var responseMessage string
	if incomingMessage == "貨幣列表" {
		currencyStr := strings.Join(currencyList, "\n")
		responseMessage = currencyStr
	} else if currencyKey, findOk := common.Mapkey(currencyStatement, strings.ToUpper(incomingMessage)); findOk{
		_, rate := currency.GetCurrencyLatestRate(strings.ToUpper(currencyKey))
		fmt.Println("rate response!")
		fmt.Println(rate)
		responseMessage = fmt.Sprintf("匯率時間: %s\n 本行買入匯率: %s\n 本行賣出匯率: %s", rate.RateTime, rate.BuyRate, rate.SellRate)
	} else {
		responseMessage = incomingMessage
	}

	var response ResponseData
	response.ReplyToken = replyToken
	response.Messages = append(response.Messages, Message{Type: "text", Text:responseMessage})

	fmt.Println(incomingMessage, response)
	responseBody, _ := json.Marshal(response)

	req, _ := http.NewRequest("POST", replyUrl, bytes.NewBuffer(responseBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + channelToken)

	clt := http.Client{}
	resp, respErr := clt.Do(req)
	if respErr != nil {
		Log.Warn("call response api error!")
	}
	fmt.Println(resp)
	defer resp.Body.Close()
}