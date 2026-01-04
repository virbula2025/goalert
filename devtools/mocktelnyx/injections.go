package mocktelnyx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/target/goalert/util/log"
)

// Messages returns all captured messages.
func (s *Server) Messages() []Message {
	s.mx.Lock()
	defer s.mx.Unlock()
	msgs := make([]Message, len(s.messages))
	copy(msgs, s.messages)
	return msgs
}

// GetTeXML fetches the current TeXML instructions from the callback URL associated with the call.
// This allows the test harness to inspect what GoAlert told Telnyx to do (e.g., "Say 'Hello'").
func (s *Server) GetTeXML(callID string) string {
	s.mx.Lock()
	call, ok := s.calls[callID]
	s.mx.Unlock()

	if !ok {
		return ""
	}
	if call.TeXMLURL == "" {
		return ""
	}

	// Request the TeXML from GoAlert
	resp, err := http.PostForm(call.TeXMLURL, url.Values{
		"CallSid":   {callID},
		"CallStatus": {call.Status},
		"From":      {call.From},
		"To":        {call.To},
	})
	if err != nil {
		log.Log(context.Background(), err)
		return ""
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	return string(data)
}

// Press simulates a user pressing digits on their keypad during a call (DTMF).
// It sends a POST request to the call's webhook URL with the 'Digits' parameter.
func (s *Server) Press(callID, digits string) error {
	s.mx.Lock()
	call, ok := s.calls[callID]
	s.mx.Unlock()

	if !ok {
		return fmt.Errorf("call %s not found", callID)
	}

	// Telnyx (via TeXML) sends inputs back to the same URL usually, or an 'action' URL.
	// For simplicity in smoke tests, we hit the configured TeXML URL.
	resp, err := http.PostForm(call.TeXMLURL, url.Values{
		"CallSid":   {callID},
		"CallStatus": {"in-progress"},
		"Digits":    {digits},
		"From":      {call.From},
		"To":        {call.To},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// RegisterSMSCallback registers a callback URL for a specific phone number.
// In a real scenario, this is done via the Telnyx Portal or API configuration.
func (s *Server) RegisterSMSCallback(from, url string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.callbacks[from] = url
}

// RegisterVoiceCallback registers a callback URL for incoming voice calls.
func (s *Server) RegisterVoiceCallback(from, url string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	// In the mock, we might store this to handle incoming call simulation if needed.
	// For outbound smoke tests, this is less critical but good for completeness.
	s.callbacks[from+"_voice"] = url
}

// UpdateCallStatus allows the harness to force a status change (e.g., "completed").
func (s *Server) UpdateCallStatus(callID, status string) {
	s.mx.Lock()
	if c, ok := s.calls[callID]; ok {
		c.Status = status
	}
	s.mx.Unlock()
}

// startCall simulates the initial webhook callback Telnyx performs when a call is answered.
func (s *Server) startCall(c *CallState) {
	// Simple delay to simulate network
	// Then hit the webhook to ask for initial TeXML
	// This is just a no-op in the mock server itself because the Harness usually
	// asserts the content by calling GetTeXML manually.
	// However, marking it in-progress is useful.
	s.mx.Lock()
	c.Status = "in-progress"
	s.mx.Unlock()
}

// RejectNextCall configures the mock to fail the next voice call attempt to this number.
func (s *Server) RejectNextCall(to string) {
    s.mx.Lock()
    defer s.mx.Unlock()
    // Implementation depends on how you want to simulate failure (e.g. 500 error or "failed" status)
    // For now, let's say we register a "fail-next" flag in the server state.
    s.failNextCall[to] = true
}

// RejectNextSMS configures the mock to simulate a delivery failure for the next SMS to `to`.
func (s *Server) RejectNextSMS(to string) {
	s.mx.Lock()
	defer s.mx.Unlock()
    // You need to handle this flag in s.handleMessages or 
    // trigger a separate callback to the status webhook with "failed".
    // 
    // Telnyx reports status asynchronously via webhook.
    // So this function should likely immediately fire the failure webhook 
    // to the URL registered for this number.
}