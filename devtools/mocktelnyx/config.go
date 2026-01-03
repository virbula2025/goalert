package mocktelnyx

import "time"

// Config is the configuration for the Telnyx mock server.
type Config struct {
	APIKey string
}

// CallState represents the current state of a voice call.
type CallState struct {
	ID        string
	Status    string
	To        string
	From      string
	AppID     string
	Direction string

	// TeXMLURL is the webhook URL to fetch instructions from.
	TeXMLURL string
}

// Message represents a captured SMS message.
type Message struct {
	ID        string
	To        string
	From      string
	Body      string
	MediaURLs []string
}