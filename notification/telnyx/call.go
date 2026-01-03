package telnyx

// Call represents a Telnyx call object.
type Call struct {
	CallControlID string `json:"call_control_id"`
	CallLegID     string `json:"call_leg_id"`
	CallSessionID string `json:"call_session_id"`
	IsAlive       bool   `json:"is_alive"`
	RecordType    string `json:"record_type"`
}

type CallInitiateRequest struct {
	To           string `json:"to"`
	From         string `json:"from"`
	ConnectionID string `json:"connection_id"`
	TeXMLUrl     string `json:"texml_url,omitempty"`
}

type CallResponse struct {
	Data struct {
		CallControlID string `json:"call_control_id"`
	} `json:"data"`
}