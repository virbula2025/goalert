package telnyx

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/target/goalert/user/contactmethod"
)

// Config contains the details needed to interact with Telnyx.
type Config struct {
	BaseURL string
	Client  *http.Client
	CMStore *contactmethod.Store
	DB      *sql.DB
}

func (c *Config) url(path string) string {
	base := c.BaseURL
	if base == "" {
		base = "https://api.telnyx.com/v2"
	}
	return strings.TrimSuffix(base, "/") + "/" + strings.TrimPrefix(path, "/")
}

func (c *Config) httpClient() *http.Client {
	if c.Client != nil {
		return c.Client
	}
	return http.DefaultClient
}