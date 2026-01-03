package telnyx

import (
	"context"
	"database/sql"

	"github.com/target/goalert/config"
	"github.com/target/goalert/notification"
)

// SMS implements the notification.Sender interface for text messages.
type SMS struct {
	*Config
}

// SendSMSInput is the JSON structure for sending a message.
type SendSMSInput struct {
	To   string `json:"to"`
	From string `json:"from"`
	Text string `json:"text"`
}

// NewSMS initializes the SMS sender.
func NewSMS(ctx context.Context, db *sql.DB, cfg *Config) (*SMS, error) {
	return &SMS{Config: cfg}, nil
}

// SendSMS is the low-level method to hit the API
func (s *SMS) SendSMS(ctx context.Context, to, body string) (*notification.SentMessage, error) {
	cfg := config.FromContext(ctx)

	input := SendSMSInput{
		To:   to,
		From: cfg.Telnyx.FromNumber,
		Text: body,
	}

	var resp MessageResponse
	err := s.postJSON(ctx, "messages", input, &resp)
	if err != nil {
		return nil, err
	}

	return &notification.SentMessage{
		ExternalID: resp.Data.ID,
		State:      mapStatus(resp.Data.Status),
		SrcValue:   resp.Data.From,
	}, nil
}