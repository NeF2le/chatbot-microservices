package server

import (
	"github.com/NeF2le/common-lib-golang/logger"
	"github.com/NeF2le/common-lib-golang/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(logger_ logger.Logger) *gin.Engine {
	router := gin.New()

	router.Use(middleware.RequestIDMiddleware(), middleware.GinLoggingMiddleware(logger_), gin.Recovery())

	router.POST("/match", HandleMatch)
	router.POST("/execute", HandleExecute)

	router.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	return router
}
