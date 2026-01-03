package telnyx

// Call represents a Telnyx call object.
type Call struct {
	CallControlID string `json:"call_control_id"`
	CallLegID     string `json:"call_leg_id"`
	CallSessionID string `json:"call_session_id"`
	IsAlive       bool   `json:"is_alive"`
	RecordType    string `json:"record_type"`
}

// CallInitiateRequest is the payload for creating a new TeXML call.
// Documentation: https://developers.telnyx.com/docs/voice/texml/rest-api/calls
type CallInitiateRequest struct {
	To   string `json:"to"`
	From string `json:"from"`
	
	// Url is the webhook URL where Telnyx will fetch the TeXML instructions.
	// This must be accessible from the public internet.
	Url string `json:"url,omitempty"`
}

type CallResponse struct {
	Data struct {
		CallControlID string `json:"call_control_id"`
	} `json:"data"`
}