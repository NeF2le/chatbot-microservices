package models

type Message struct {
	UserID string `json:"user_id"`
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type Reply struct {
	Reply string `json:"reply"`
}
