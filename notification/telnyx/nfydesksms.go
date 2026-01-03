package telnyx

import (
	"context"
	"fmt"

	"github.com/nyaruka/phonenumbers"
	"github.com/target/goalert/config"
	"github.com/target/goalert/notification"
	"github.com/target/goalert/notification/nfydest"
	"github.com/target/goalert/validation"
)

const (
	DestTypeTelnyxSMS  = "builtin-telnyx-sms"
	FieldPhoneNumber   = "phone_number"
	FallbackIconURLSMS = "builtin://phone-text"
)

var _ nfydest.Provider = (*SMS)(nil)

// ID returns the unique identifier for this provider type.
func (s *SMS) ID() string { return DestTypeTelnyxSMS }

// TypeInfo returns metadata about the provider type, including UI fields.
func (s *SMS) TypeInfo(ctx context.Context) (*nfydest.TypeInfo, error) {
	cfg := config.FromContext(ctx)
	return &nfydest.TypeInfo{
		Type:                       DestTypeTelnyxSMS,
		Name:                       "Text Message (Telnyx)",
		Enabled:                    cfg.Telnyx.Enable,
		UserDisclaimer:             cfg.General.NotificationDisclaimer,
		SupportsAlertNotifications: true,
		SupportsUserVerification:   true,
		SupportsStatusUpdates:      true,
		UserVerificationRequired:   true,
		RequiredFields: []nfydest.FieldConfig{{
			FieldID:            FieldPhoneNumber,
			Label:              "Phone Number",
			Hint:               "Include country code e.g. +1 (USA), +91 (India), +44 (UK)",
			PlaceholderText:    "11235550123",
			Prefix:             "+",
			InputType:          "tel",
			SupportsValidation: true,
		}},
	}, nil
}

// ValidateField validates the input phone number.
func (s *SMS) ValidateField(ctx context.Context, fieldID, value string) error {
	switch fieldID {
	case FieldPhoneNumber:
		n, err := phonenumbers.Parse(value, "")
		if err != nil {
			return validation.WrapError(err)
		}
		if !phonenumbers.IsValidNumber(n) {
			return validation.NewGenericError("invalid phone number")
		}
		return nil
	}

	return validation.NewGenericError("unknown field ID")
}

// DisplayInfo formats the destination for the UI.
func (s *SMS) DisplayInfo(ctx context.Context, args map[string]string) (*nfydest.DisplayInfo, error) {
	if args == nil {
		args = make(map[string]string)
	}

	n, err := phonenumbers.Parse(args[FieldPhoneNumber], "")
	if err != nil {
		return nil, validation.WrapError(err)
	}

	return &nfydest.DisplayInfo{
		IconURL:     FallbackIconURLSMS,
		IconAltText: "Text Message",
		Text:        phonenumbers.Format(n, phonenumbers.INTERNATIONAL),
	}, nil
}

// Send implements the notification.Sender interface.
func (s *SMS) Send(ctx context.Context, msg notification.Message) (*notification.SentMessage, error) {
	// Helper to handle destination value extraction safely
	getDest := func(d interface{ Value() (interface{}, error) }) (string, error) {
		val, err := d.Value()
		if err != nil {
			return "", err
		}
		if str, ok := val.(string); ok {
			return str, nil
		}
		return "", fmt.Errorf("telnyx: invalid destination type %T", val)
	}

	switch m := msg.(type) {
	case notification.Test:
		dest, err := getDest(m.Dest)
		if err != nil {
			return nil, err
		}
		return s.SendSMS(ctx, dest, "GoAlert Test Message")

	case notification.Alert:
		return s.SendSMSAlert(ctx, m)

	case notification.Verification:
		dest, err := getDest(m.Dest)
		if err != nil {
			return nil, err
		}
		body := fmt.Sprintf("Your GoAlert verification code is: %d", m.Code)
		return s.SendSMS(ctx, dest, body)
	}

	return nil, fmt.Errorf("telnyx: unsupported message type %T", msg)
}

func (s *SMS) Status(ctx context.Context, id, providerID string) (*notification.Status, error) {
	return nil, nil
}