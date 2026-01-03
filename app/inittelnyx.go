package app

import (
	"context"

	"github.com/pkg/errors"
	"github.com/target/goalert/notification/telnyx"
)

func (app *App) initTelnyx(ctx context.Context) error {
	// Initialize the main Telnyx configuration object
	app.telnyxConfig = &telnyx.Config{
		BaseURL: app.cfg.TelnyxBaseURL, // You must add this field to your global Config struct
		CMStore: app.ContactMethodStore,
		DB:      app.db,
		Client:  app.httpClient,
	}

	var err error
	// Initialize the SMS sub-system
	app.telnyxSMS, err = telnyx.NewSMS(ctx, app.db, app.telnyxConfig)
	if err != nil {
		return errors.Wrap(err, "init TelnyxSMS")
	}

	// Initialize the Voice sub-system
	app.telnyxVoice, err = telnyx.NewVoice(ctx, app.db, app.telnyxConfig)
	if err != nil {
		return errors.Wrap(err, "init TelnyxVoice")
	}

	return nil
}