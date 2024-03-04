package app

import (
	"github.com/lmnq/test-thai/config"
	"github.com/lmnq/test-thai/database/postgres"
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

	
}
