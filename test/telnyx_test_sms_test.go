package smoke

import (
	"testing"

	"github.com/target/goalert/test/smoke/harness"
)

// TestTelnyxSMS tests that SMS messages are sent via Telnyx.
func TestTelnyxSMS(t *testing.T) {
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

	// Assert Telnyx received the SMS request
	h.Telnyx().Device(h.Phone("1")).ExpectSMS("testing")
}