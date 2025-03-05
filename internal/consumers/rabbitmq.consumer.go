package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gabrigabs/campaign-message-consumer/internal/models"
	repository "github.com/gabrigabs/campaign-message-consumer/internal/repositories"
	"github.com/gabrigabs/campaign-message-consumer/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	connection         *amqp.Connection
	channel            *amqp.Channel
	queue              string
	messageRepository  *repository.MessageRepository
	campaignRepository *repository.PostgresRepository
	logger             logger.Logger
	shutdownChannel    chan struct{}
}

func NewRabbitMQConsumer(
	url string,
	queue string,
	messageRepository *repository.MessageRepository,
	campaignRepository *repository.PostgresRepository,
	log logger.Logger,
) (*RabbitMQConsumer, error) {
	log.Info("Connecting to RabbitMQ", map[string]any{
		"url":   url,
		"queue": queue,
	})

	connection, err := amqp.Dial(url)
	if err != nil {
		log.Error("Failed to connect to RabbitMQ", map[string]any{
			"error": err.Error(),
		})
		return nil, err
	}

	channel, err := connection.Channel()
	if err != nil {
		log.Error("Failed to open a channel", map[string]any{
			"error": err.Error(),
		})
		connection.Close()
		return nil, err
	}

	_, err = channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error("Failed to declare a queue", map[string]any{
			"error": err.Error(),
		})
		channel.Close()
		connection.Close()
		return nil, err
	}

	log.Info("Connected to RabbitMQ", map[string]any{
		"queue": queue,
	})

	return &RabbitMQConsumer{
		connection:         connection,
		channel:            channel,
		queue:              queue,
		messageRepository:  messageRepository,
		campaignRepository: campaignRepository,
		logger:             log,
		shutdownChannel:    make(chan struct{}),
	}, nil
}

func (consumer *RabbitMQConsumer) Start() error {
	consumer.logger.Info("Starting RabbitMQ consumer", map[string]any{
		"queue": consumer.queue,
	})

	err := consumer.channel.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		consumer.logger.Error("Failed to set QoS", map[string]any{
			"error": err.Error(),
		})
		return err
	}

	msgs, err := consumer.channel.Consume(
		consumer.queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		consumer.logger.Error("Failed to register a consumer", map[string]any{
			"error": err.Error(),
		})
		return err
	}

	go func() {
		for {
			select {
			case msg, ok := <-msgs:
				if !ok {
					consumer.logger.Warn("RabbitMQ channel closed", nil)
					return
				}
				consumer.processMessage(msg)
			case <-consumer.shutdownChannel:
				consumer.logger.Info("Shutting down consumer", nil)
				return
			}

		}

	}()

	return nil
}

func (consumer *RabbitMQConsumer) processMessage(msg amqp.Delivery) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	consumer.logger.Debug("Processing message", map[string]any{
		"deliveryTag": msg.DeliveryTag,
	})

	var rabbitMessage struct {
		PhoneNumber   string `json:"phone_number"`
		Message       string `json:"message"`
		CampaignID    string `json:"campaign_id"`
		CompanyID     string `json:"company_id"`
		IsLastMessage bool   `json:"is_last_message"`
	}

	err := json.Unmarshal(msg.Body, &rabbitMessage)
	if err != nil {
		consumer.logger.Error("Failed to parse message", map[string]any{
			"error": err.Error(),
			"body":  string(msg.Body),
		})
		msg.Reject(false)
		return
	}

	if rabbitMessage.CampaignID == "" || rabbitMessage.PhoneNumber == "" ||
		rabbitMessage.Message == "" || rabbitMessage.CompanyID == "" {
		consumer.logger.Error("Invalid message format", map[string]any{
			"campaignId":  rabbitMessage.CampaignID,
			"phoneNumber": rabbitMessage.PhoneNumber,
			"companyId":   rabbitMessage.CompanyID,
		})
		msg.Reject(false)
		return
	}

	message := models.Message{
		PhoneNumber: rabbitMessage.PhoneNumber,
		Message:     rabbitMessage.Message,
		CampaignID:  rabbitMessage.CampaignID,
		CompanyID:   rabbitMessage.CompanyID,
	}

	err = consumer.messageRepository.SaveMessage(ctx, message)
	if err != nil {
		consumer.logger.Error("Failed to save message to MongoDB", map[string]any{
			"error": err.Error(),
		})
		msg.Nack(false, true)
		return
	}

	if rabbitMessage.IsLastMessage {

		err = consumer.campaignRepository.UpdateCampaignStatus(ctx, message.CampaignID)
		if err != nil {
			consumer.logger.Error("Failed to update campaign status to SENT", map[string]any{
				"error":      err.Error(),
				"campaignId": message.CampaignID,
			})
			msg.Ack(false)
			return
		}
	}

	consumer.logger.Info("Message processed successfully", map[string]any{
		"messageId":   message.ID,
		"campaignId":  message.CampaignID,
		"phoneNumber": message.PhoneNumber,
	})

	msg.Ack(false)
}
