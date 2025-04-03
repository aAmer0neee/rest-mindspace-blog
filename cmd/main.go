package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/aAmer0neee/rest-mindspace-blog/internal/auth"
	"github.com/aAmer0neee/rest-mindspace-blog/internal/config"
	"github.com/aAmer0neee/rest-mindspace-blog/internal/handlers"
	"github.com/aAmer0neee/rest-mindspace-blog/internal/logger"
	"github.com/aAmer0neee/rest-mindspace-blog/internal/logo"
	"github.com/aAmer0neee/rest-mindspace-blog/internal/repository"
	"github.com/aAmer0neee/rest-mindspace-blog/internal/service"
)

func main() {

	cfg := config.LoadConfig()

	logger := logger.ConfigureLogger(cfg.Server.Env)

	repo, err := repository.ConnectRepository(*cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	service := service.ConfigureService(repo, logger)

	authService := auth.ConfigureJWT(cfg.Auth.SecretJWT)

	r := gin.Default()
	handlers.Configurehandlers(r, service, authService)

	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}

	logo.Create()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			slog.Info(err.Error())
			os.Exit(1)
		}
	}()

	shutdown(logger, srv)
}

func shutdown(logger *slog.Logger, srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Graceful Shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	srv.Shutdown(ctx)
}
