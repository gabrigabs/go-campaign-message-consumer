package repository

import (
	"context"
	"database/sql"

	"github.com/gabrigabs/campaign-message-consumer/logger"
)

type PostgresRepository struct {
	db     *sql.DB
	logger logger.Logger
}

func NewCampaignRepository(db *sql.DB, log logger.Logger) *PostgresRepository {
	return &PostgresRepository{
		db:     db,
		logger: log,
	}
}

const (
	CampaignStatusPending = "PENDING"
	CampaignStatusSent    = "SENT"
)

func (repository *PostgresRepository) UpdateCampaignStatus(ctx context.Context, campaignID string) error {
	repository.logger.Debug("Updating campaign status to SENT in PostgreSQL", map[string]interface{}{
		"campaignId": campaignID,
	})

	query := `
		UPDATE "Campaign"
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := repository.db.ExecContext(ctx, query, CampaignStatusSent, campaignID)
	if err != nil {
		repository.logger.Error("Failed to update campaign status in PostgreSQL", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		repository.logger.Warn("No campaign found to update in PostgreSQL", map[string]interface{}{
			"campaignId": campaignID,
		})
	} else {
		repository.logger.Info("Campaign status updated to SENT in PostgreSQL", map[string]interface{}{
			"campaignId": campaignID,
		})
	}

	return nil
}
