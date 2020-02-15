package common

import (
	"encoding/json"
)

var (
	channelToken = "E6lkGp635WG1GKcKCViXIn5xCPPudDEvefYpFRRmhiOHgLzvUcQpoOo+3ZhGfKhymd1Mujyfj0ddLHpNPLjJm2GPw3pynN7KPuJQ9aUrvcSrQue7ibw8el1eO/Xnm+qUHCnFdrQcPF6Z6n000j2sMQdB04t89/1O/w1cDnyilFU="
	replyUrl     = "https://api.line.me/v2/bot/message/reply"
	pushUrl      = "https://api.line.me/v2/bot/message/push"
)

type ResponseData struct {
	ReplyToken string    `json:"replyToken"`
	Messages   []Message `json:"messages"`
}

type PushData struct {
	To       string    `json:"to"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func ReplyMessage(token string, message string) {
	var response ResponseData
	response.ReplyToken = token
	response.Messages = append(response.Messages, Message{Type: "text", Text: message})
	replyToUser(response, replyUrl)
}

func PushMessage(target string, message string) {
	var response PushData
	response.To = target
	response.Messages = append(response.Messages, Message{Type: "text", Text: message})
	replyToUser(response, pushUrl)
}

func replyToUser(response interface{}, apiUrl string) {
	responseBody, _ := json.Marshal(response)
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + channelToken
	CallAPI(apiUrl, "POST", responseBody, header)
}
