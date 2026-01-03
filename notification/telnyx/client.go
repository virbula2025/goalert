package telnyx

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/target/goalert/config"
)

func (c *Config) postJSON(ctx context.Context, endpoint string, payload interface{}, out interface{}) error {
	cfg := config.FromContext(ctx)
	apiKey := cfg.Telnyx.APIKey

	var buf bytes.Buffer
	if payload != nil {
		if err := json.NewEncoder(&buf).Encode(payload); err != nil {
			return err
		}
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.url(endpoint), &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return DecodeError(resp.Body)
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}