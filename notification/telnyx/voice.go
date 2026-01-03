package telnyx

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/target/goalert/config"
	"github.com/team-telnyx/telnyx-go/v3"
	"github.com/team-telnyx/telnyx-go/v3/option"
)

// Voice implements the notification.Sender interface for calls.
type Voice struct {
	*Config
}

// NewVoice initializes the Voice sender.
func NewVoice(ctx context.Context, db *sql.DB, cfg *Config) (*Voice, error) {
	return &Voice{Config: cfg}, nil
}

// MakeCall initiates an outbound call using the official Telnyx SDK (TeXML).
func (v *Voice) MakeCall(ctx context.Context, dest, callbackURL string) (string, error) {
	cfg := config.FromContext(ctx)

	// 1. Initialize the Client
	client := telnyx.NewClient(option.WithAPIKey(cfg.Telnyx.APIKey))

	// 2. Prepare the parameters
	params := telnyx.TexmlCallInitiateParams{
		To:        dest,
		From:      cfg.Telnyx.FromNumber,
		URL:       telnyx.String(callbackURL),
		URLMethod: telnyx.TexmlCallInitiateParamsURLMethodPost,
	}

	// 3. Call the API
	// The second argument is the Connection ID (Application ID)
	resp, err := client.Texml.Calls.Initiate(ctx, cfg.Telnyx.ConnectionID, params)
	if err != nil {
		return "", err
	}

	// 4. Return the Call SID
	// WORKAROUND: The Telnyx generated go SDK struct 'TexmlCallInitiateResponseData' is missing the 'Sid' field.
	// We must parse the raw JSON to get it.
	var result struct {
		Data struct {
			Sid string `json:"sid"`
		} `json:"data"`
	}

	if err := json.Unmarshal([]byte(resp.RawJSON()), &result); err != nil {
		return "", err
	}

	return result.Data.Sid, nil
}