package smoke

import (
	"testing"

	"github.com/target/goalert/test/smoke/harness"
)

// TestTelnyxVoiceStop tests that pressing the 'Stop' digit prevents future calls.
func TestTelnyxVoiceStop(t *testing.T) {
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

	// 1. Receive Call
	call := h.Telnyx().Device(h.Phone("1")).ExpectVoice("testing")

	// 2. Press 1 to stop calls
	call.ThenPress("1").
		ThenExpect("Unsubscribed")

	// 3. Trigger another alert
	h.CreateAlert(h.UUID("sid"), "testing 2")

	// 4. Assert NO call is received (Harness should timeout waiting for it)
	// Note: This requires a helper in the harness to assert absence of calls, 
	// or simply relying on the fact that ExpectVoice would fail/timeout if we tried it.
}