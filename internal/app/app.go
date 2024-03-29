package app

import (
	"blog-backend/config"
	v1 "blog-backend/internal/controller/http/v1"
	"blog-backend/internal/repo"
	"blog-backend/internal/usecase"
	"blog-backend/pkg/hasher"
	"blog-backend/pkg/httpserver"
	"blog-backend/pkg/postgres"
	"blog-backend/pkg/validator"
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func Run(configPath string) {
	// Configuration
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Logger
	SetLogrus(cfg.Log.Level)

	// Repositories
	log.Info("Initializing postgres...")
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - pgdb.NewUseCases: %w", err))
	}
	defer pg.Close()

	// Repositories
	log.Info("Initializing repositories...")
	repositories := repo.NewRepositories(pg)

	// UseCases dependencies
	log.Info("Initializing useCases...")
	deps := usecase.UseCasesDependencies{
		Repos:    repositories,
		Hasher:   hasher.NewSHA1Hasher(cfg.Hasher.Salt),
		SignKey:  cfg.JWT.SignKey,
		TokenTTL: cfg.JWT.TokenTTL,
	}
	useCases := usecase.NewUseCases(deps)

	// Echo handler
	log.Info("Initializing handlers and routes...")
	handler := echo.New()
	// setup handler validator as lib validator
	handler.Validator = validator.NewCustomValidator()
	v1.NewRouter(handler, useCases)

	// HTTP server
	log.Info("Starting http server...")
	log.Debugf("Server port: %s", cfg.HTTP.Port)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	log.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Graceful shutdown
	log.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
