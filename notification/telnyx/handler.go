package telnyx

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/target/goalert/alert"
	"github.com/target/goalert/config"
)

// Handler processes incoming webhooks from Telnyx.
type Handler struct {
	aStore *alert.Store
}

// NewHandler creates a new Telnyx webhook handler with access to the AlertStore.
func NewHandler(aStore *alert.Store) *Handler {
	return &Handler{
		aStore: aStore,
	}
}

// ServeHTTP handles the incoming webhook.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cfg := config.FromContext(ctx)

	// 1. Determine what type of message this is based on URL query params
	q := r.URL.Query()
	msgType := q.Get("type")

	var msgText string
	var err error

	switch msgType {
	case "alert":
		// Fetch the actual alert details from the DB
		idStr := q.Get("alertID")
		id, _ := strconv.Atoi(idStr)
		
		var a *alert.Alert
		a, err = h.aStore.FindOne(ctx, id)
		if err != nil {
			// Fallback if the alert is deleted or lookup fails
			fmt.Printf("Telnyx: failed to lookup alert %d: %v\n", id, err)
			msgText = "Critical Alert from GoAlert. Please check your dashboard."
		} else {
			// Construct the voice message
			// "Alert 123: Server Down. Details: CPU is at 100%..."
			msgText = fmt.Sprintf("Alert %d: %s. %s", a.ID, a.Summary, a.Details)
		}

	case "verify":
		// Read code directly from params
		code := q.Get("code")
		msgText = fmt.Sprintf("Your GoAlert verification code is: %s. Repeat: %s.", code, code)

	case "test":
		msgText = "This is a test message from GoAlert. If you are hearing this, your voice configuration is correct."

	default:
		msgText = "GoAlert Notification System."
	}

	// 2. Define the XML structure
	type Response struct {
		XMLName xml.Name `xml:"Response"`
		Say     SayVerb  `xml:"Say"`
	}

	// 3. Use the helper to apply Voice/Language settings from Config
	say := NewSayVerb(cfg, msgText)

	// 4. Send Response
	w.Header().Set("Content-Type", "application/xml")
	if err := xml.NewEncoder(w).Encode(Response{Say: say}); err != nil {
		http.Error(w, "failed to encode xml", http.StatusInternalServerError)
	}
}