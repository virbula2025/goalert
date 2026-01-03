package telnyx

import "github.com/target/goalert/notification"

type MessageResponse struct {
	Data struct {
		ID       string `json:"id"`
		To       []struct{ PhonePhoneNumber string `json:"phone_number"` } `json:"to"`
		From     string `json:"from"`
		Text     string `json:"text"`
		Status   string `json:"status"` // queued, sending, sent, delivered
	} `json:"data"`
}

func mapStatus(s string) notification.State {
	switch s {
	case "queued", "sending":
		return notification.StateSending
	case "sent":
		return notification.StateSent
	case "delivered":
		return notification.StateDelivered
	case "failed", "undelivered":
		return notification.StateFailedPerm
	default:
		return notification.StateSent // Default fallback
	}
}