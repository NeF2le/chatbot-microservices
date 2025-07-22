package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/NeF2le/chatbot-microservices/skills/echo-skill/internal/config"
	"github.com/NeF2le/chatbot-microservices/skills/echo-skill/internal/server"
	"github.com/NeF2le/common-lib-golang/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error when reading config in echo-skill service, %s", err)
	}

	logger_ := logger.NewZapLogger(false)

	router := server.NewRouter(logger_)

	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.ShutdownTimeout,
	}

	go func() {
		logger_.Info(fmt.Sprintf("start echo-skill service on %s", addr), nil)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger_.Fatal(fmt.Sprintf("error when starting echo-skill service: %s", err), nil)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger_.Fatal(fmt.Sprintf("echo-skill service shutdown: %s", err), nil)
	} else {
		logger_.Info("echo-skill service exiting", nil)
	}
}
