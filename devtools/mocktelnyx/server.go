package mocktelnyx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/target/goalert/util/log"
)

// Server implements a mock Telnyx API server.
type Server struct {
	mux       *http.ServeMux
	cfg       Config
	mx        sync.Mutex
	messages  []Message
	calls     map[string]*CallState
	callbacks map[string]string // number -> url
}

// NewServer creates a new Server instance.
func NewServer(cfg Config) *Server {
	s := &Server{
		mux:       http.NewServeMux(),
		cfg:       cfg,
		calls:     make(map[string]*CallState),
		callbacks: make(map[string]string),
	}

	// Telnyx V2 API Routes
	s.mux.HandleFunc("/v2/messages", s.handleMessages)
	// The SDK appends the ApplicationID to the path for TeXML calls
	s.mux.HandleFunc("/v2/texml/calls/", s.handleCalls)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer "+s.cfg.APIKey {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	s.mux.ServeHTTP(w, r)
}

// handleMessages handles POST /v2/messages (Sending SMS)
func (s *Server) handleMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		From string `json:"from"`
		To   string `json:"to"`
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mx.Lock()
	msg := Message{
		ID:   uuid.NewString(),
		From: req.From,
		To:   req.To,
		Body: req.Text,
	}
	s.messages = append(s.messages, msg)
	s.mx.Unlock()

	w.Header().Set("Content-Type", "application/json")
	// Response matching Telnyx V2 Message Object structure
	json.NewEncoder(w).Encode(struct {
		Data struct {
			ID   string `json:"id"`
			From string `json:"from"`
			To   string `json:"to"`
			Text string `json:"text"`
		} `json:"data"`
	}{
		Data: struct {
			ID   string `json:"id"`
			From string `json:"from"`
			To   string `json:"to"`
			Text string `json:"text"`
		}{
			ID:   msg.ID,
			From: msg.From,
			To:   msg.To,
			Text: msg.Body,
		},
	})
}

// handleCalls handles POST /v2/texml/calls/{ApplicationID} (Initiating Calls)
func (s *Server) handleCalls(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract Application ID from URL path: /v2/texml/calls/<AppID>
	parts := strings.Split(r.URL.Path, "/")
	appID := parts[len(parts)-1]

	var req struct {
		To        string `json:"To"`
		From      string `json:"From"`
		Url       string `json:"Url"` // Webhook URL for TeXML
		UrlMethod string `json:"UrlMethod"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := uuid.NewString()

	s.mx.Lock()
	call := &CallState{
		ID:        id,
		Status:    "queued",
		To:        req.To,
		From:      req.From,
		AppID:     appID,
		Direction: "outbound-api",
		TeXMLURL:  req.Url,
	}
	s.calls[id] = call
	s.mx.Unlock()

	// Simulate the call starting asynchronously
	go s.startCall(call)

	w.Header().Set("Content-Type", "application/json")
	// Response must include `sid` (mapped to id in SDK) inside `data`
	json.NewEncoder(w).Encode(struct {
		Data struct {
			Sid    string `json:"sid"`
			Status string `json:"status"`
		} `json:"data"`
	}{
		Data: struct {
			Sid    string `json:"sid"`
			Status string `json:"status"`
		}{
			Sid:    id,
			Status: "queued",
		},
	})
}