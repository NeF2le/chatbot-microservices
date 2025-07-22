package server_test

import (
	"encoding/json"
	"github.com/NeF2le/chatbot-microservices/skills/echo-skill/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleExecute(t *testing.T) {
	router := gin.New()
	router.POST("/execute", server.HandleExecute)

	wr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/execute", strings.NewReader(`{"text":"hello world"}`))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(wr, req)

	type replyBotCore struct {
		Reply string `json:"reply"`
	}
	reply := replyBotCore{}
	if err := json.Unmarshal(wr.Body.Bytes(), &reply); err != nil {
		t.Errorf("failed to unmarshal reply: %v", err)
	}

	require.Equal(t, "hello world", reply.Reply)
}
