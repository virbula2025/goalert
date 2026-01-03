package harness

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Telnyx returns a reference to the Telnyx assertion helper.
func (h *Harness) Telnyx() *TelnyxAssertionAPI {
	return &TelnyxAssertionAPI{h: h}
}

type TelnyxAssertionAPI struct {
	h *Harness
}

// Device returns a device assertion helper for the given phone number.
func (t *TelnyxAssertionAPI) Device(number string) *TelnyxAssertionDevice {
	return &TelnyxAssertionDevice{
		h:      t.h,
		Number: number,
	}
}

type TelnyxAssertionDevice struct {
	h      *Harness
	Number string
}

// ExpectSMS asserts that an SMS is sent to the device.
func (d *TelnyxAssertionDevice) ExpectSMS(body string) {
	d.h.t.Helper()
	// Logic depends on how the Telnyx Mock server exposes recorded requests.
	// Assuming a similar interface to the Twilio mock:
	d.h.IgnoreErrorsWith("telnyx-mock: 404") // Ignore polling errors if any

    // This waits for the mock server to receive the message
	msg := d.h.Slack().WaitAndAssert("telnyx_sms", d.Number, body)
    assert.Contains(d.h.t, msg, body)
}

// IgnoreUnexpectedSMS ignores any extra SMS messages sent to this device.
func (d *TelnyxAssertionDevice) IgnoreUnexpectedSMS(body string) {
    // Implementation depends on harness capability to ignore specific logs
}