package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/lmnq/test-thai/config"
	"github.com/lmnq/test-thai/database"
	"github.com/lmnq/test-thai/database/postgres"
	"github.com/lmnq/test-thai/fasthttpserver"
	"github.com/lmnq/test-thai/internal/controller"
	"github.com/lmnq/test-thai/internal/repo"
	"github.com/lmnq/test-thai/internal/service"
	"github.com/lmnq/test-thai/logger"
)

// Run configures and runs the app -.
func Run(cfg *config.Config) {
	// logger - init zap or zerolog logger. Default is zap.
	// If you want to use zerolog, call logger.NewZerolog(cfg.Log.Level)
	l := logger.NewZap(cfg.Log.Level)

	// database - init database connection
	pg, err := postgres.New(cfg.Db.PgURL)
	if err != nil {
		l.Fatal("database error", err)
	}
	defer pg.Close()

	// migrate - init database migration
	database.Migrate(cfg.Db.PgURL)

	// repos
	repos := repo.New(pg)

	// services
	services := service.New(repos)

	// HTTP server
	fiberApp := fiber.New(fiber.Config{AppName: cfg.App.Name})
	controller.New(fiberApp, l, services)
	fastHTTPServer := fasthttpserver.New(fiberApp.Handler(), cfg.HTTP.Port)

	// signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("signal received", s.String())
	case err = <-fastHTTPServer.Notify():
		l.Error(fmt.Errorf("fastHTTPServer error: %w", err))
	}

	// shutdown
	fastHTTPServer.Shutdown()
	l.Info("server shutdown")
}
