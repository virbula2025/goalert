package telnyx

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type APIError struct {
	Code   string `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

type APIErrorResponse struct {
	Errors []APIError `json:"errors"`
}

func (e *APIErrorResponse) Error() string {
	var msgs []string
	for _, err := range e.Errors {
		msgs = append(msgs, fmt.Sprintf("%s: %s", err.Code, err.Detail))
	}
	return "telnyx: " + strings.Join(msgs, ", ")
}

func DecodeError(r io.Reader) error {
	var e APIErrorResponse
	if err := json.NewDecoder(r).Decode(&e); err != nil {
		return fmt.Errorf("telnyx: unknown error (parse fail): %v", err)
	}
	return &e
}