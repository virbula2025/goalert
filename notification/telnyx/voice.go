package telnyx

import (
	"context"
	"database/sql"

	"github.com/target/goalert/config"
)

// Voice implements the notification.Sender interface for calls.
type Voice struct {
	*Config
}

// NewVoice initializes the Voice sender.
func NewVoice(ctx context.Context, db *sql.DB, cfg *Config) (*Voice, error) {
	return &Voice{Config: cfg}, nil
}

func (v *Voice) MakeCall(ctx context.Context, to string, callbackURL string) (string, error) {
	cfg := config.FromContext(ctx)

	req := CallInitiateRequest{
		To:           to,
		From:         cfg.Telnyx.FromNumber,
		ConnectionID: cfg.Telnyx.ConnectionID,
		TeXMLUrl:     callbackURL,
	}

	var resp CallResponse
	err := v.postJSON(ctx, "calls", req, &resp)
	if err != nil {
		return "", err
	}

	return resp.Data.CallControlID, nil
}