package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type messageDTO struct {
	UserID string `json:"userID"`
	ChatID string `json:"chatID"`
	Text   string `json:"text"`
}

type matchResp struct {
	Match bool `json:"match"`
}

type execResp struct {
	Reply string `json:"reply"`
}

func HandleMatch(c *gin.Context) {
	c.JSON(http.StatusOK, matchResp{Match: true})
}

func HandleExecute(c *gin.Context) {
	var msg messageDTO
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, execResp{Reply: msg.Text})
}
