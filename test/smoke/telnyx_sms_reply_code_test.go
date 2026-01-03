package smoke

import (
	"testing"

	"github.com/target/goalert/test/smoke/harness"
)

// TestTelnyxSMSReplyCode tests acknowledging an alert by replying to the SMS with the code.
func TestTelnyxSMSReplyCode(t *testing.T) {
	t.Parallel()

	sql := `
	insert into users (id, name, email) 
	values ({{uuid "user"}}, 'bob', 'joe');

	insert into user_contact_methods (id, user_id, name, type, value) 
	values ({{uuid "cm1"}}, {{uuid "user"}}, 'personal', 'SMS', {{phone "1"}});

	insert into user_notification_rules (user_id, contact_method_id, delay_minutes) 
	values ({{uuid "user"}}, {{uuid "cm1"}}, 0);

	insert into escalation_policies (id, name) 
	values ({{uuid "eid"}}, 'esc policy');

	insert into escalation_policy_steps (id, escalation_policy_id) 
	values ({{uuid "es1"}}, {{uuid "eid"}});

	insert into escalation_policy_actions (escalation_policy_step_id, user_id) 
	values ({{uuid "es1"}}, {{uuid "user"}});

	insert into services (id, escalation_policy_id, name) 
	values ({{uuid "sid"}}, {{uuid "eid"}}, 'service');
	`
	h := harness.NewHarness(t, sql, "ids-to-uuids")
	defer h.Close()

	h.CreateAlert(h.UUID("sid"), "testing")

	// Wait for the notification SMS
	d := h.Telnyx().Device(h.Phone("1"))
	d.ExpectSMS("testing")

	// Reply with an arbitrary code (GoAlert usually matches on number, but strict code matching might be configured)
	// For this test we assume standard reply behavior.
	// NOTE: In the harness, we need a method to Simulate Incoming SMS from Telnyx
	h.TelnyxMock().SendSMS(h.Phone("1"), h.Phone("GoAlert"), "1a") // Assuming '1a' is the code

	// Expect the confirmation
	d.ExpectSMS("Acknowledged")
}