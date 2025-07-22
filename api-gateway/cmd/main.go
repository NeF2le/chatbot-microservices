package main

import (
	"github.com/NeF2le/chatbot-microservices/api-gateway/internal/config"
	"github.com/NeF2le/chatbot-microservices/api-gateway/internal/server"
	"github.com/NeF2le/common-lib-golang/logger"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error when reading config in api-gateway: %s", err)
	}

	logger_ := logger.NewZapLogger(false)
	router := server.NewRouter(cfg, logger_)

	server.RunServer(logger_, cfg, router)
}
