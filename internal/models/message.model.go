package models

import (
	"time"

	"github.com/nrednav/cuid2"
)

type Message struct {
	ID          string    `json:"id" bson:"_id"`
	PhoneNumber string    `json:"phone_number" bson:"phone_number"`
	Message     string    `json:"message" bson:"message"`
	CampaignID  string    `json:"campaign_id" bson:"campaign_id"`
	CompanyID   string    `json:"company_id" bson:"company_id"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

func NewMessage(phoneNumber, messageContent, campaignID, companyID string) Message {
	now := time.Now()
	return Message{
		ID:          cuid2.Generate(),
		PhoneNumber: phoneNumber,
		Message:     messageContent,
		CampaignID:  campaignID,
		CompanyID:   companyID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
