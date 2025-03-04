package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gabrigabs/campaign-message-consumer/config"
	"github.com/gabrigabs/campaign-message-consumer/db"
	consumer "github.com/gabrigabs/campaign-message-consumer/internal/consumers"
	repository "github.com/gabrigabs/campaign-message-consumer/internal/repositories"
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

	mongodb, err := db.NewMongoConnection(cfg.MongoDB, log)
	if err != nil {
		log.Error("Failed to connect to MongoDB", map[string]any{"error": err.Error()})
	}

	campaignRepo := repository.NewCampaignRepository(postgres.DB, log)
	messageRepo := repository.NewMessageRepository(mongodb.Database, log)

	rabbitConsumer, err := consumer.NewRabbitMQConsumer(
		cfg.RabbitMQ.URL,
		cfg.RabbitMQ.Queue,
		messageRepo,
		campaignRepo,
		log,
	)

	if err != nil {
		log.Error("Failed to create RabbitMQ consumer", map[string]any{"error": err.Error()})
	}

	rabbitConsumer.Start()

	log.Info("Campaign message consumer is running.", nil)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	log.Info("Received shutdown signal", map[string]interface{}{"signal": sig.String()})

}
