package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/NeF2le/chatbot-microservices/bot-core/internal/config"
	"github.com/NeF2le/common-lib-golang/logger"
	"net/http"
	"os"
	"os/signal"
)

func RunServer(logger_ logger.Logger, router http.Handler, cfg *config.Config) {
	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.ShutdownTimeout,
	}

	go func() {
		logger_.Info(fmt.Sprintf("start bot-core service on %s", addr), nil)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger_.Fatal(fmt.Sprintf("error when start bot-core service: %s", err), nil)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger_.Fatal(fmt.Sprintf("bot-core service shutdown: %s", err), nil)
	} else {
		logger_.Info("bot-core service exiting", nil)
	}
}
