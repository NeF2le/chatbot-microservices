package model

type Message struct {
	ChatId string `json:"chat_id"`
	UserId string `json:"user_id"`
	Text   string `json:"text"`
}

type Reply struct {
	Reply string `json:"reply"`
}
