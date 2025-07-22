package server_test

import (
	"encoding/json"
	"github.com/NeF2le/chatbot-microservices/api-gateway/internal/config"
	"github.com/NeF2le/chatbot-microservices/api-gateway/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Message struct {
	ChatId string `json:"chat_id"`
	UserId string `json:"user_id"`
	Text   string `json:"text"`
}

type Reply struct {
	Reply string `json:"reply"`
}

func TestNewProxy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var botCoreCalled bool
	botCoreSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		botCoreCalled = true
		require.Equal(t, "/message", r.URL.Path)
		var msg Message
		require.Error(t, json.NewDecoder(r.Body).Decode(&msg))
		require.Error(t, json.NewEncoder(w).Encode(Reply{Reply: msg.Text}))
	}))
	defer botCoreSrv.Close()
	
	var tgAdapterSrvCalled bool
	tgAdapterSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tgAdapterSrvCalled = true

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"reply": "Hello, World!"}`))
	}))
	defer tgAdapterSrv.Close()

	servicesConfig := config.ServicesConfig{
		TelegramAdapter: config.ServiceConfig{URL: tgAdapterSrv.URL},
		BotCore:         config.ServiceConfig{URL: ""},
	}
	httpConfig := config.HTTPConfig{}
	cfg := &config.Config{
		Services: servicesConfig,
		HTTP:     httpConfig,
	}
	router := gin.New()
	router.POST("/telegram/*proxyPath", server.NewProxy(cfg.Services.TelegramAdapter.URL))
	router.POST("/bot-core/*proxyPath", server.NewProxy(cfg.Services.BotCore.URL))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/telegram/webhook", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, true, tgAdapterSrvCalled, "TelegramAdapter was not called")
	require.Equal(t, true, botCoreCalled, "BotCore was not called")
}
