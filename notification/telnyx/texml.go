package telnyx

import (
	"encoding/xml"
	"io"
)

// TeXML mirrors Twilio's TwiML but is compatible with Telnyx.
type TeXML struct {
	XMLName xml.Name `xml:"Response"`
	Verbs   []interface{}
}

type Say struct {
	XMLName xml.Name `xml:"Say"`
	Voice   string   `xml:"voice,attr,omitempty"`
	Text    string   `xml:",chardata"`
}

type Play struct {
	XMLName xml.Name `xml:"Play"`
	URL     string   `xml:",chardata"`
}

type Dial struct {
	XMLName xml.Name `xml:"Dial"`
	Timeout int      `xml:"timeout,attr,omitempty"`
	Action  string   `xml:"action,attr,omitempty"`
	Number  string   `xml:",chardata"`
}

type Gather struct {
	XMLName     xml.Name `xml:"Gather"`
	Action      string   `xml:"action,attr,omitempty"`
	NumDigits   int      `xml:"numDigits,attr,omitempty"`
	Timeout     int      `xml:"timeout,attr,omitempty"`
	Verbs       []interface{}
}

type Hangup struct {
	XMLName xml.Name `xml:"Hangup"`
}

func (t *TeXML) Add(v interface{}) {
	t.Verbs = append(t.Verbs, v)
}

func (t *TeXML) Encode(w io.Writer) error {
	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	return xml.NewEncoder(w).Encode(t)
}