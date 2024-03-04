package main

import (
	"log"

	"github.com/lmnq/test-thai/config"
	"github.com/lmnq/test-thai/internal/app"
)

func main() {
	// config
	config, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	app.Run(config)
}
