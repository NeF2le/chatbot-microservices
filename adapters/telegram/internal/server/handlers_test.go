package server

import (
	"bytes"
	"encoding/json"
	"github.com/NeF2le/chatbot-microservices/adapters/telegram/internal/config"
	"github.com/NeF2le/chatbot-microservices/adapters/telegram/internal/model"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestHandleWebhook(t *testing.T) {
	var botCoreCalled bool
	botCoreSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		botCoreCalled = true
		var msg model.Message
		require.NoError(t, json.NewDecoder(r.Body).Decode(&msg))
		require.Equal(t, "24", msg.ChatId)
		require.Equal(t, "42", msg.UserId)
		require.Equal(t, "Hello, World!", msg.Text)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"reply": "Hello, World!"}`))
	}))
	defer botCoreSrv.Close()

	var tgPayload struct {
		Text   string `json:"text"`
		ChatID string `json:"chat_id"`
	}
	var telegramCalled bool
	telegramSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		telegramCalled = true
		json.NewDecoder(r.Body).Decode(&tgPayload)
		w.Write([]byte(`{"ok": true}`))
	}))
	defer telegramSrv.Close()

	servicesConfig := config.ServicesConfig{
		BotCore: config.ServiceConfig{URL: botCoreSrv.URL},
	}
	httpConfig := config.HTTPConfig{}
	cfg := &config.Config{
		Services: servicesConfig,
		HTTP:     httpConfig,
	}

	client := &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			switch {
			case req.URL.Host == botCoreSrv.Listener.Addr().String():
				return http.DefaultTransport.RoundTrip(req)
			case req.URL.Host == "api.telegram.org":
				req.URL.Scheme = "http"
				req.URL.Host = telegramSrv.Listener.Addr().String()
				return http.DefaultTransport.RoundTrip(req)
			default:
				t.Fatalf("unexpected request host:port: %s", req.URL.Host)
				return nil, nil
			}
		}),
	}

	router := gin.New()
	router.POST("/telegram/webhook", HandleWebhook(cfg, client))

	tgUpdate := tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{ID: 42},
			Chat: &tgbotapi.Chat{ID: 24},
			Text: "Hello, World!",
		},
	}
	updJSON, _ := json.Marshal(tgUpdate)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/telegram/webhook", bytes.NewBuffer(updJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.True(t, botCoreCalled, "botCore handler was not called")
	require.True(t, telegramCalled, "telegram handler was not called")
	require.Equal(t, "24", tgPayload.ChatID)
	require.Equal(t, "Hello, World!", tgPayload.Text)
}
