package smoke

import (
	"testing"

	"github.com/target/goalert/test/smoke/harness"
)

// TestTelnyxVoiceClose tests that an alert can be closed by pressing '5' during a call.
func TestTelnyxVoiceClose(t *testing.T) {
	t.Parallel()

	sql := `
	insert into users (id, name, email) 
	values ({{uuid "user"}}, 'bob', 'joe');

	insert into user_contact_methods (id, user_id, name, type, value) 
	values ({{uuid "cm1"}}, {{uuid "user"}}, 'personal', 'VOICE', {{phone "1"}});

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

	// Answer call
	call := h.Telnyx().Device(h.Phone("1")).ExpectVoice("testing")

	// Press 5 to close (Standard GoAlert menu: 4=Ack, 5=Close)
	call.ThenPress("5").
		ThenExpect("Closed")
}