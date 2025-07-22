package dispatcher

import (
	"bytes"
	"encoding/json"
	"github.com/NeF2le/chatbot-microservices/bot-core/internal/models"
	"net/http"
)

type HTTPSkill struct {
	BaseURL    string
	HTTPClient *http.Client
}

func (h *HTTPSkill) Match(msg models.Message) bool {
	data, _ := json.Marshal(msg)
	resp, err := h.HTTPClient.Post(h.BaseURL+"/match", "application/json", bytes.NewBuffer(data))
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}

	var r struct{ Match bool }
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return false
	}

	return r.Match
}

func (h *HTTPSkill) Execute(msg models.Message) (string, error) {
	data, _ := json.Marshal(msg)
	resp, err := h.HTTPClient.Post(h.BaseURL+"/execute", "application/json", bytes.NewBuffer(data))
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", err
	}

	var r struct{ Reply string }
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", err
	}

	return r.Reply, nil
}
