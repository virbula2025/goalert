package smoke

import (
	"testing"

	"github.com/target/goalert/test/smoke/harness"
)

// TestTelnyxSMSVerification tests that adding an SMS contact method triggers a verification code.
func TestTelnyxSMSVerification(t *testing.T) {
	t.Parallel()

	sql := `
	insert into users (id, name, email) 
	values ({{uuid "user"}}, 'bob', 'joe');
	`
	h := harness.NewHarness(t, sql, "ids-to-uuids")
	defer h.Close()

	// Create contact method
	h.GraphQLToken(h.UUID("user")).CreateContactMethod(harness.CreateContactMethodInput{
		Name:  "SMS",
		Type:  "SMS",
		Value: h.Phone("1"),
	})

	// Trigger verification
	h.GraphQLToken(h.UUID("user")).SendContactMethodVerification(h.UUID("cm1"))

	// Expect SMS with 6-digit code
	h.Telnyx().Device(h.Phone("1")).ExpectSMS("verification code")
}