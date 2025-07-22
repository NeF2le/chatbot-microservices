package server

import (
	"github.com/NeF2le/chatbot-microservices/bot-core/internal/dispatcher"
	"github.com/NeF2le/common-lib-golang/logger"
	"github.com/NeF2le/common-lib-golang/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(logger_ logger.Logger, disp *dispatcher.Dispatcher) *gin.Engine {
	router := gin.New()

	router.Use(middleware.RequestIDMiddleware(), middleware.GinLoggingMiddleware(logger_), gin.Recovery())

	router.POST("/message", HandleMessage(disp))

	router.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	return router
}
