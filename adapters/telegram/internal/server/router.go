package server

import (
	"github.com/NeF2le/chatbot-microservices/adapters/telegram/internal/config"
	"github.com/NeF2le/common-lib-golang/logger"
	"github.com/NeF2le/common-lib-golang/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter(logger_ logger.Logger, cfg *config.Config, client *http.Client) *gin.Engine {
	router := gin.New()

	router.Use(middleware.RequestIDMiddleware(), middleware.GinLoggingMiddleware(logger_), gin.Recovery())

	router.POST("/webhook", HandleWebhook(cfg, client))

	router.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	return router
}
