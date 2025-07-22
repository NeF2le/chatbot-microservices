package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/NeF2le/chatbot-microservices/api-gateway/internal/config"
	"github.com/NeF2le/common-lib-golang/logger"
	"net/http"
	"os"
	"os/signal"
)

func RunServer(logger_ logger.Logger, cfg *config.Config, router http.Handler) {
	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	go func() {
		logger_.Info(fmt.Sprintf("api-gateway listening on %s", addr), nil)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger_.Fatal(fmt.Sprintf("error when start api-gateway: %s", err), nil)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger_.Fatal(fmt.Sprintf("shutdown api-gateway with error %s", err), nil)
	} else {
		logger_.Info("graceful shutdown of api-gateway", nil)
	}
}
