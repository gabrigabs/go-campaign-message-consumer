package repository

import (
	"context"
	"time"

	"github.com/gabrigabs/campaign-message-consumer/internal/models"
	"github.com/gabrigabs/campaign-message-consumer/logger"
	"github.com/nrednav/cuid2"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageRepository struct {
	collection *mongo.Collection
	logger     logger.Logger
}

func NewMessageRepository(db *mongo.Database, log logger.Logger) *MessageRepository {
	return &MessageRepository{
		collection: db.Collection("messages"),
		logger:     log,
	}
}

func (repository *MessageRepository) SaveMessage(ctx context.Context, message models.Message) error {
	repository.logger.Debug("Saving message to MongoDB", map[string]interface{}{
		"messageId":   message.ID,
		"campaignId":  message.CampaignID,
		"phoneNumber": message.PhoneNumber,
	})

	if message.ID == "" {
		message.ID = cuid2.Generate()
	}

	now := time.Now()
	if message.CreatedAt.IsZero() {
		message.CreatedAt = now
	}
	message.UpdatedAt = now

	_, err := repository.collection.InsertOne(ctx, message)
	if err != nil {
		repository.logger.Error("Failed to save message to MongoDB", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	repository.logger.Info("Message saved to MongoDB", map[string]interface{}{
		"messageId": message.ID,
	})
	return nil
}
