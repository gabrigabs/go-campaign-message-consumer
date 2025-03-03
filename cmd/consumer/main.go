package main

import (
	"github.com/gabrigabs/campaign-message-consumer/config"
	"github.com/gabrigabs/campaign-message-consumer/db"
	"github.com/gabrigabs/campaign-message-consumer/logger"
)

func main() {

	cfg, err := config.LoadConfig()

	log := logger.NewLogger(cfg.App.LogLevel)

	if err != nil {
		log.Error("Failed to load configuration: %v", map[string]any{"error": err.Error()})
	}

	postgres, err := db.NewPostgresConnection(cfg.Postgres, log)
	if err != nil {
		log.Error("Failed to connect to PostgreSQL", map[string]any{"error": err.Error()})
	}

	postgres.DB.Query("SELECT * FROM campaigns")

	log.Info("Hello World", nil)

}
