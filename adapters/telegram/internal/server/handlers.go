package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/NeF2le/chatbot-microservices/adapters/telegram/internal/config"
	"github.com/NeF2le/chatbot-microservices/adapters/telegram/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TgUpdate struct {
	UpdateId int `json:"update_id"`
	Message  struct {
		MessageID int    `json:"message_id"`
		Text      string `json:"text,omitempty"`
		Chat      struct {
			ID int64 `json:"id"`
		}
		From struct {
			ID int64 `json:"id"`
		}
	}
}

func HandleWebhook(cfg *config.Config, client *http.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var update TgUpdate
		if err := c.ShouldBindJSON(&update); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "invalid JSON data",
				"details": err.Error(),
			})
			return
		}

		if update.Message.MessageID == 0 || update.Message.Chat.ID == 0 || update.Message.From.ID == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "missing required fields in update",
			})
			return
		}

		msg := model.Message{
			ChatId: strconv.FormatInt(update.Message.Chat.ID, 10),
			UserId: strconv.FormatInt(update.Message.From.ID, 10),
			Text:   update.Message.Text,
		}

		body, err := json.Marshal(msg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "failed to marshal update: " + err.Error(),
			})
			return
		}

		resp, err := client.Post(cfg.Services.BotCore.URL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "failed to call bot-core: " + err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		var reply model.Reply
		if err := json.NewDecoder(resp.Body).Decode(&reply); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   "invalid JSON data",
				"details": err.Error(),
			})
			return
		}

		botToken := c.Query("token")
		if botToken != "" {
			sendMessageURL := fmt.Sprintf(
				"https://api.telegram.org/bot%s/sendMessage",
				botToken,
			)

			payload := map[string]interface{}{
				"chat_id": msg.ChatId,
				"text":    reply.Reply,
			}

			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "failed to marshal payload: " + err.Error(),
				})
			}

			tgResp, err := client.Post(sendMessageURL, "application/json", bytes.NewBuffer(payloadBytes))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send message to telegram: " + err.Error()})
				return
			}
			defer tgResp.Body.Close()

			if tgResp.StatusCode != http.StatusOK {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send message to telegram: " + tgResp.Status})
				return
			}

			c.JSON(http.StatusOK, "message sent to telegram")
		} else {
			c.JSON(http.StatusOK, gin.H{"reply": reply.Reply})
		}

	}
}
