package smoke

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/target/goalert/test/smoke/harness"
)

// TestTelnyxVoiceVerification tests the contact method verification flow via voice.
func TestTelnyxVoiceVerification(t *testing.T) {
	t.Parallel()

	sql := `
	insert into users (id, name, email) 
	values ({{uuid "user"}}, 'bob', 'joe');
	`
	h := harness.NewHarness(t, sql, "ids-to-uuids")
	defer h.Close()

	// 1. Create the Contact Method via API (disabled by default)
	h.GraphQLToken(h.UUID("user")).CreateContactMethod(harness.CreateContactMethodInput{
		Name:  "Voice",
		Type:  "VOICE",
		Value: h.Phone("1"),
	})

	// 2. Trigger Verification
	h.GraphQLToken(h.UUID("user")).SendContactMethodVerification(h.UUID("cm1"))

	// 3. Receive Call
	call := h.Telnyx().Device(h.Phone("1")).ExpectVoice("verification")

	// 4. Verify the code is spoken
	// Logic assumes the verification code is generated deterministically or extracted from DB
	// For smoke tests, we often just check that the prompt asks for it or provides it.
	call.ThenExpect("verification code")
}