package telnyx

import (
	"context"
	"fmt"

	"github.com/target/goalert/config"
	"github.com/target/goalert/notification"
)

// SendSMSAlert formats the SMS specifically for an alert notification.
func (c *Config) SendSMSAlert(ctx context.Context, a notification.Alert) (*notification.SentMessage, error) {
	cfg := config.FromContext(ctx)

	// Fix: Dest.Value is a function in your version, so we call it.
	// We handle the error or cast carefully.
	valRaw, err := a.Dest.Value()
	if err != nil {
		return nil, fmt.Errorf("telnyx: failed to get destination value: %w", err)
	}
	destStr, ok := valRaw.(string)
	if !ok {
		return nil, fmt.Errorf("telnyx: destination value is not a string: %T", valRaw)
	}

	// Fix: a.Status is undefined, so we use Summary and ID.
	msg := fmt.Sprintf("Alert #%d: %s", a.AlertID, a.Summary)

	// Basic truncation
	if len(msg) > 160 {
		msg = msg[:157] + "..."
	}

	// Re-implement the body of SendSMS inline here to avoid circular dependencies or signature mismatch
	input := SendSMSInput{
		To:   destStr,
		From: cfg.Telnyx.FromNumber,
		Text: msg,
	}

	var resp MessageResponse
	err = c.postJSON(ctx, "messages", input, &resp)
	if err != nil {
		return nil, err
	}

	return &notification.SentMessage{
		ExternalID: resp.Data.ID,
		State:      mapStatus(resp.Data.Status),
		SrcValue:   resp.Data.From,
	}, nil
}