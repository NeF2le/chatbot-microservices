package server

import (
	"github.com/NeF2le/chatbot-microservices/api-gateway/internal/config"
	"github.com/NeF2le/common-lib-golang/logger"
	"github.com/NeF2le/common-lib-golang/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config, logger_ logger.Logger) *gin.Engine {
	router := gin.New()

	router.Use(middleware.RequestIDMiddleware(), middleware.GinLoggingMiddleware(logger_), gin.Recovery())

	apiV1 := router.Group("/api/v1")
	apiV1.Any("/telegram/*proxyPath", NewProxy(cfg.Services.TelegramAdapter.URL))
	apiV1.Any("/bot-core/*proxyPath", NewProxy(cfg.Services.BotCore.URL))

	apiV1.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })

	return router
}
