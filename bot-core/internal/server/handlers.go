package server

import (
	"github.com/NeF2le/chatbot-microservices/bot-core/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/NeF2le/chatbot-microservices/bot-core/internal/dispatcher"
)

func HandleMessage(disp *dispatcher.Dispatcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		var msg models.Message
		if err := c.BindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		reply, err := disp.Dispatch(msg)
		if reply.Reply == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "skill returned empty reply"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, reply)
	}
}
