package main

import (
	"github.com/NeF2le/chatbot-microservices/bot-core/internal/config"
	"github.com/NeF2le/chatbot-microservices/bot-core/internal/dispatcher"
	"github.com/NeF2le/chatbot-microservices/bot-core/internal/server"
	"github.com/NeF2le/common-lib-golang/logger"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error when reading config in bot-core service, %s", err)
	}

	logger_ := logger.NewZapLogger(false)

	httpClient := &http.Client{
		Timeout: cfg.HTTP.WriteTimeout + cfg.HTTP.ReadTimeout,
	}

	skills := make(map[string]dispatcher.Skill)
	for name, skillCfg := range cfg.Skills {
		skills[name] = &dispatcher.HTTPSkill{
			BaseURL:    skillCfg.URL,
			HTTPClient: httpClient,
		}
	}

	disp := dispatcher.New(skills)
	router := server.NewRouter(logger_, disp)

	server.RunServer(logger_, router, cfg)
}
