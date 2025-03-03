package db

import (
	"context"
	"time"

	"github.com/gabrigabs/campaign-message-consumer/config"
	"github.com/gabrigabs/campaign-message-consumer/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
	logger   logger.Logger
}

func NewMongoConnection(cfg config.MongoDBConfig, log logger.Logger) (*MongoDB, error) {
	log.Info("Connecting to MongoDB", map[string]any{
		"uri": cfg.URI,
		"db":  cfg.DBName,
	})

	clientOptions := options.Client().ApplyURI(cfg.URI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Error("Failed to connect to MongoDB", map[string]any{
			"error": err.Error(),
		})
		return nil, err
	}

	database := client.Database(cfg.DBName)

	log.Info("Connected to MongoDB", map[string]any{
		"database": cfg.DBName,
	})

	return &MongoDB{
		client:   client,
		database: database,
		logger:   log,
	}, nil
}
