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
	DestTypeTelnyxVoice  = "builtin-telnyx-voice"
	FallbackIconURLVoice = "builtin://phone-voice"
)

var _ nfydest.Provider = (*Voice)(nil)

// ID returns the unique identifier for this provider type.
func (v *Voice) ID() string { return DestTypeTelnyxVoice }

// TypeInfo returns metadata about the provider type.
func (v *Voice) TypeInfo(ctx context.Context) (*nfydest.TypeInfo, error) {
	cfg := config.FromContext(ctx)
	return &nfydest.TypeInfo{
		Type:                       DestTypeTelnyxVoice,
		Name:                       "Voice Call (Telnyx)",
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
func (v *Voice) ValidateField(ctx context.Context, fieldID, value string) error {
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
func (v *Voice) DisplayInfo(ctx context.Context, args map[string]string) (*nfydest.DisplayInfo, error) {
	if args == nil {
		args = make(map[string]string)
	}

	n, err := phonenumbers.Parse(args[FieldPhoneNumber], "")
	if err != nil {
		return nil, validation.WrapError(err)
	}

	return &nfydest.DisplayInfo{
		IconURL:     FallbackIconURLVoice,
		IconAltText: "Voice Call",
		Text:        phonenumbers.Format(n, phonenumbers.INTERNATIONAL),
	}, nil
}

// Send implements the notification.Sender interface.
func (v *Voice) Send(ctx context.Context, msg notification.Message) (*notification.SentMessage, error) {
	cfg := config.FromContext(ctx)
	callback := cfg.CallbackURL("/api/v2/telnyx/voice")

	extractDest := func(d interface{ Value() (interface{}, error) }) (string, error) {
		val, err := d.Value()
		if err != nil {
			return "", err
		}
		if s, ok := val.(string); ok {
			return s, nil
		}
		return "", fmt.Errorf("telnyx: invalid destination type %T", val)
	}

	switch m := msg.(type) {
	case notification.Test:
		dest, err := extractDest(m.Dest)
		if err != nil {
			return nil, err
		}
		id, err := v.MakeCall(ctx, dest, callback)
		if err != nil {
			return nil, err
		}
		return &notification.SentMessage{
			ExternalID: id,
			State:      notification.StateSending,
			SrcValue:   dest,
		}, nil

	case notification.Alert:
		dest, err := extractDest(m.Dest)
		if err != nil {
			return nil, err
		}
		id, err := v.MakeCall(ctx, dest, callback)
		if err != nil {
			return nil, err
		}
		return &notification.SentMessage{
			ExternalID: id,
			State:      notification.StateSending,
			SrcValue:   dest,
		}, nil
	}

	return nil, fmt.Errorf("telnyx: unsupported message type %T", msg)
}

func (v *Voice) Status(ctx context.Context, id, providerID string) (*notification.Status, error) {
	return nil, nil
}