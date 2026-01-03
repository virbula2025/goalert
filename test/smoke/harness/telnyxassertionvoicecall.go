package harness

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TelnyxAssertionVoiceCall struct {
	h      *Harness
	CallID string
	to     string
}

// ExpectVoice asserts that a voice call is made to the device.
func (d *TelnyxAssertionDevice) ExpectVoice(name string) *TelnyxAssertionVoiceCall {
	d.h.t.Helper()
    // Wait for the outgoing call to hit the mock server
	call := d.h.WaitAndAssertTelnyxCall(d.Number) 
	return &TelnyxAssertionVoiceCall{
		h:      d.h,
		CallID: call.ID,
		to:     d.Number,
	}
}

// ThenPress simulates DTMF input during the call.
func (c *TelnyxAssertionVoiceCall) ThenPress(digits string) *TelnyxAssertionVoiceCall {
	c.h.t.Helper()
    // Post to the mock server to simulate DTMF
	c.h.TelnyxMock().Press(c.CallID, digits)
	return c
}

// ThenExpect asserts that the specific text is spoken (via Polly/Text-to-Speech)
func (c *TelnyxAssertionVoiceCall) ThenExpect(text string) *TelnyxAssertionVoiceCall {
	c.h.t.Helper()
    // Fetch the current TeXML execution state from the mock
	texml := c.h.TelnyxMock().GetTeXML(c.CallID)
	assert.Contains(c.h.t, texml, text)
	return c
}

// Status updates the status of the call (e.g., "completed", "busy").
func (c *TelnyxAssertionVoiceCall) Status(status string) {
    c.h.TelnyxMock().UpdateCallStatus(c.CallID, status)
}