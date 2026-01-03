package telnyx

import (
	"encoding/xml"
	"github.com/target/goalert/config"
)

// SayVerb represents the <Say> element in TeXML
type SayVerb struct {
	XMLName  xml.Name `xml:"Say"`
	Voice    string   `xml:"voice,attr,omitempty"`
	Language string   `xml:"language,attr,omitempty"`
	Text     string   `xml:",chardata"`
}

// Helper to create the Say verb using config
func NewSayVerb(cfg config.Config, text string) SayVerb {
    // Default to 'alice' if not set
    voice := cfg.Telnyx.VoiceName
    if voice == "" {
        voice = "alice"
    }
    
    // Default to 'en-US' if not set
    lang := cfg.Telnyx.VoiceLanguage
    if lang == "" {
        lang = "en-US"
    }

	return SayVerb{
		Voice:    voice,
		Language: lang,
		Text:     text,
	}
}