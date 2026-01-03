package smoke

import (
	"testing"

	"github.com/target/goalert/test/smoke/harness"
)

// TestTelnyxSMSStopStart tests that the system respects STOP commands.
func TestTelnyxSMSStopStart(t *testing.T) {
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

	// 1. Send first alert
	h.CreateAlert(h.UUID("sid"), "testing 1")
	d := h.Telnyx().Device(h.Phone("1"))
	d.ExpectSMS("testing 1")

	// 2. Reply STOP
	h.TelnyxMock().SendSMS(d.Number, h.Phone("GoAlert"), "STOP")
	
	// Wait for processing (GoAlert usually disables the CM)
	h.Trigger()

	// 3. Send second alert - Expect NO SMS
	h.CreateAlert(h.UUID("sid"), "testing 2")
	d.IgnoreUnexpectedSMS("testing 2") // Should fail if it arrives, or use a helper "ExpectNoSMS" if available
	
	// 4. Reply START (or UNSTOP)
	h.TelnyxMock().SendSMS(d.Number, h.Phone("GoAlert"), "START")
	h.Trigger()

	// 5. Send third alert - Expect Success
	h.CreateAlert(h.UUID("sid"), "testing 3")
	d.ExpectSMS("testing 3")
}