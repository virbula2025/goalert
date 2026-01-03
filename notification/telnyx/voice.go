package telnyx

import (
	"context"
	"database/sql"
	"fmt"

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

// teXMLRequest is the payload specifically for the TeXML API.
type teXMLRequest struct {
	To   string `json:"to"`
	From string `json:"from"`
	Url  string `json:"url"` // TeXML expects "url", Call Control expects "webhook_url"
}

func (v *Voice) MakeCall(ctx context.Context, to string, callbackURL string) (string, error) {
	cfg := config.FromContext(ctx)

	// 1. Use a local struct or ensure CallInitiateRequest has `json:"url"`
	req := teXMLRequest{
		To:   to,
		From: cfg.Telnyx.FromNumber,
		Url:  callbackURL,
	}

	// 2. Construct the correct TeXML endpoint path
	// The Connection ID must be in the URL path, not the body.
	// We assume v.postJSON appends this path to the base Telnyx URL (https://api.telnyx.com/v2/)
	path := fmt.Sprintf("texml/calls/%s", cfg.Telnyx.ConnectionID)

	var resp CallResponse
	// 3. Send to the TeXML endpoint
	err := v.postJSON(ctx, path, req, &resp)
	if err != nil {
		return "", err
	}

	return resp.Data.CallControlID, nil
}