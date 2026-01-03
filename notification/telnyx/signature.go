package telnyx

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"net/http"
)

// ValidateSignature checks the X-Telnyx-Signature-Ed25519 header.
func ValidateSignature(r *http.Request, body []byte, publicKeyStr string) error {
	sigHeader := r.Header.Get("X-Telnyx-Signature-Ed25519")
	timestamp := r.Header.Get("X-Telnyx-Timestamp")
	
	if sigHeader == "" || timestamp == "" {
		return errors.New("telnyx: missing signature headers")
	}

	// 1. Decode Public Key
	pubKey, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		return errors.New("telnyx: invalid public key configuration")
	}

	// 2. Decode Signature
	sig, err := base64.StdEncoding.DecodeString(sigHeader)
	if err != nil {
		return errors.New("telnyx: invalid signature encoding")
	}

	// 3. Construct Payload: Timestamp + Body
	payload := append([]byte(timestamp), body...)

	// 4. Verify
	if !ed25519.Verify(pubKey, payload, sig) {
		return errors.New("telnyx: signature validation failed")
	}

	return nil
}